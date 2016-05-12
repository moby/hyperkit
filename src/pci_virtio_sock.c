/*-
 * Copyright (c) 2016 Docker, Inc.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer
 *    in this position and unchanged.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

/*
 * virtio vsock emulation based on v4 specification
 *    http://markmail.org/message/porhou5zv3wqjz6h
 * Tested against the Linux implementation at
 *    git@github.com:stefanha/linux.git#vsock @ 563d2a770dfa
 * Backported to v4.1.19:
 *    git cherry-pick -x 11aa9c2 f6a835b 4ef7ea9 8566b86 \
 *                       ea3803c a9f9df1 1bb5b77 0c734eb \
 *                       139bbcd 563d2a7
 */

#include <sys/types.h>
#include <sys/socket.h>
#include <sys/uio.h>
#include <sys/un.h>
#include <sys/time.h>

#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include <fcntl.h>
#include <inttypes.h>
#include <strings.h>
#include <unistd.h>
#include <errno.h>

#include <xhyve/pci_emul.h>
#include <xhyve/virtio.h>
#include <xhyve/xhyve.h>

#define VTSOCK_RINGSZ 256

#define VTSOCK_QUEUE_RX		0
#define VTSOCK_QUEUE_TX		1
#define VTSOCK_QUEUE_EVT	2
#define VTSOCK_QUEUES		3

#define VTSOCK_MAXSEGS		32

#define VTSOCK_MAXSOCKS	1024
#define VTSOCK_MAXFWDS	4

/*
 * Host capabilities
 */
#define VTSOCK_S_HOSTCAPS 0
#if 0
	(VIRTIO_RING_F_INDIRECT_DESC) /* indirect descriptors */
#endif

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wpacked"
/*
 * Config space "registers"
 */
struct vtsock_config {
	uint32_t guest_cid;
} __packed;

/*
 * Fixed-size block header
 */

struct virtio_sock_hdr {
	uint32_t src_cid;
	uint32_t src_port;
	uint32_t dst_cid;
	uint32_t dst_port;
	uint32_t len;
#define VIRTIO_VSOCK_TYPE_STREAM 1
	uint16_t type;
#define VIRTIO_VSOCK_OP_INVALID 0
	/* Connect operations */
#define VIRTIO_VSOCK_OP_REQUEST 1
#define VIRTIO_VSOCK_OP_RESPONSE 2
#define VIRTIO_VSOCK_OP_RST 3
#define VIRTIO_VSOCK_OP_SHUTDOWN 4
	/* To send payload */
#define VIRTIO_VSOCK_OP_RW 5
	/* Tell the peer our credit info */
#define VIRTIO_VSOCK_OP_CREDIT_UPDATE 6
	/* Request the peer to send the credit info to us */
#define VIRTIO_VSOCK_OP_CREDIT_REQUEST 7
	uint16_t op;
	uint32_t flags;
#define VIRTIO_VSOCK_FLAG_SHUTDOWN_RX (1U<<0) /* Peer will not receive any more data */
#define VIRTIO_VSOCK_FLAG_SHUTDOWN_TX (1U<<1) /* Peer will not transmit any more data */
#define VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL (VIRTIO_VSOCK_FLAG_SHUTDOWN_RX|VIRTIO_VSOCK_FLAG_SHUTDOWN_TX)
	uint32_t buf_alloc;
	uint32_t fwd_cnt;
} __packed;

#pragma clang diagnostic pop

/*
 * Debug printf
 */
static int pci_vtsock_debug = 0;
#define DPRINTF(params) if (pci_vtsock_debug) printf params
/* Protocol logging */
#define PPRINTF(params) if (0) printf params

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wpadded"

/* XXX need to use rx and tx more consistently */

struct vsock_addr {
	uint32_t cid, port;
};
#define PRIcid "%08"PRIx32
#define PRIport "%08"PRIx32

#ifdef PRI_ADDR_PREFIX
#define PRIaddr PRI_ADDR_PREFIX PRIcid "." PRIport
#else
#define PRIaddr PRIcid "." PRIport
#endif

#ifndef CONNECT_SOCKET_NAME
#define CONNECT_SOCKET_NAME "connect"
#endif

#define FMTADDR(a) a.cid, a.port

#define WRITE_BUF_LENGTH (128*1024)

struct pci_vtsock_sock {
	pthread_mutex_t mtx;

	/*
	 * To allocate a sock:
	 *
	 *   Grab alloc_mtx
	 *
	 *   Loop over socks[] looking for a FREE sock
	 *
	 *   If a FREE sock is found take its lock before setting it to state == CONNECTING.
	 *
	 *   Finally drop alloc_mtx
	 *
	 * To free a sock:
	 *
	 *   Set state to CLOSING and kick the rx sock.
	 *
	 *   The RX loop will close the fd and set state == FREE
	 *
	 *   This does not require the alloc_mtx.
	 */
	enum {
		SOCK_FREE, /* Initial state */
		SOCK_CONNECTING,
		SOCK_CONNECTED,
		SOCK_CLOSING,
	} state;
	/* fd is:
	 *   >= 0	When state == CONNECTED,
	 *     -1	Otherwise
	 */
	int fd;
	/* valid when SOCK_CONNECTED only */
	uint32_t local_shutdown, peer_shutdown;

	struct vsock_addr local_addr;
	struct vsock_addr peer_addr;

	uint32_t buf_alloc;
	uint32_t fwd_cnt;

	bool credit_update_required;
	uint32_t rx_cnt; /* Amount we have sent to the peer */
	uint32_t peer_buf_alloc; /* From the peer */
	uint32_t peer_fwd_cnt; /* From the peer */

	/* Write buffer. We do not update fwd_cnt until we drain the _whole_ buffer */
	uint8_t write_buf[WRITE_BUF_LENGTH];
	unsigned int write_buf_head, write_buf_tail;
};

struct pci_vtsock_forward {
	int listen_fd;
	uint32_t port;
};

/*
 * Per-device softc
 */
/*
 * Lock order (outer most first): XXX more thought needed.
 *
 *   vssc_mtx is taken by the core and is often held during callbacks
 *   (e.g. it is held during a vq_notify or pci cfg access). It
 *   protects virtio resources.
 *
 *   tx_mtx protects the tx data structures, including the queue.
 *
 *   rx_mtx protects the rx data structures, including the queue.
 *
 *   alloc_mtx protects transitions from socks[...].state == FREE => *
 *
 *   reply_mtx protects reply_{ring,prod,cons}
 *
 *   sock->mtx protects the contents of the sock struct, including the
 *   state.
 *
 */
struct pci_vtsock_softc {
	struct virtio_softc vssc_vs;
	pthread_mutex_t vssc_mtx;
	char *path;
	struct vqueue_info vssc_vqs[VTSOCK_QUEUES];
	struct vtsock_config vssc_cfg;

	pthread_mutex_t alloc_mtx;
	struct pci_vtsock_sock socks[VTSOCK_MAXSOCKS];

	struct pci_vtsock_forward fwds[VTSOCK_MAXFWDS];
	int nr_fwds;

	/* Protects the following plus VTSOCK_QUEUE_TX */
	pthread_mutex_t tx_mtx;
	pthread_t tx_thread;
	int tx_kick_fd, tx_wake_fd; /* Write to kick, select on wake */
	int connect_fd; /* */

	/* Protects the following plus VTSOCK_QUEUE_RX */
	pthread_mutex_t rx_mtx;
	pthread_t rx_thread;
	int rx_kick_fd, rx_wake_fd; /* Write to kick, select on wake */

	pthread_mutex_t reply_mtx;
#define VTSOCK_REPLYRINGSZ (2*VTSOCK_RINGSZ)
	struct virtio_sock_hdr reply_ring[VTSOCK_REPLYRINGSZ];
	int reply_prod, reply_cons;
};

#pragma clang diagnostic pop

/* Protocol stuff */

/* Reserved CIDs */
#define VMADDR_CID_ANY -1U
//#define VMADDR_CID_HYPERVISOR 0
//#define VMADDR_CID_RESERVED 1
#define VMADDR_CID_HOST 2

static void pci_vtsock_reset(void *);
static void pci_vtsock_notify_tx(void *, struct vqueue_info *);
static void pci_vtsock_notify_rx(void *, struct vqueue_info *);
static int pci_vtsock_cfgread(void *, int, int, uint32_t *);
static int pci_vtsock_cfgwrite(void *, int, int, uint32_t);
static void *pci_vtsock_rx_thread(void *vssc);

static struct virtio_consts vtsock_vi_consts = {
	"vtsock", /* our name */
	VTSOCK_QUEUES,
	sizeof(struct vtsock_config), /* config reg size */
	pci_vtsock_reset, /* reset */
	NULL, /* no device-wide qnotify */
	pci_vtsock_cfgread, /* read PCI config */
	pci_vtsock_cfgwrite, /* write PCI config */
	NULL, /* apply negotiated features */
	VTSOCK_S_HOSTCAPS, /* our capabilities */
};

static void pci_vtsock_reset(void *vsc)
{
	struct pci_vtsock_softc *sc = vsc;

	DPRINTF(("vtsock: device reset requested !\n"));
	vi_reset_dev(&sc->vssc_vs);
	/* XXX TODO: close/reset all socks */
}

static const char * const opnames[] = {
	[VIRTIO_VSOCK_OP_INVALID] = "INVALID",
	[VIRTIO_VSOCK_OP_REQUEST] = "REQUEST",
	[VIRTIO_VSOCK_OP_RESPONSE] = "RESPONSE",
	[VIRTIO_VSOCK_OP_RST] = "RST",
	[VIRTIO_VSOCK_OP_SHUTDOWN] = "SHUTDOWN",
	[VIRTIO_VSOCK_OP_RW] = "RW",
	[VIRTIO_VSOCK_OP_CREDIT_UPDATE] = "CREDIT_UPDATE",
	[VIRTIO_VSOCK_OP_CREDIT_REQUEST] = "CREDIT_REQUEST"
};

static int max_fd(int a, int b)
{
	if (a > b)
		return a;
	else
		return b;
}

static size_t iovec_clip(struct iovec **iov, int *iov_len, size_t bytes)
{
	size_t ret = 0;
	int i;
	for (i = 0; i < *iov_len && ret < bytes; i++) {
		if ((bytes-ret) < (*iov)[i].iov_len)
			(*iov)[i].iov_len = bytes - ret;
		ret += (*iov)[i].iov_len;
	}
	*iov_len = i;
	return ret;
}

/* Pulls @bytes from @iov into @buf. @buf can be NULL, in which case this just discards @bytes */
static size_t iovec_pull(struct iovec **iov, int *iov_len, void *buf, size_t bytes)
{
	size_t res = 0;

	//DPRINTF(("iovec_pull %zd bytes into %p. iov=%p, iov_len=%d\n",
	//	 bytes, (void *)buf, (void *)*iov, *iov_len));

	while (res < bytes && *iov_len) {
		size_t c = (bytes - res) < (*iov)[0].iov_len ? (bytes - res) : (*iov)[0].iov_len;

		//DPRINTF(("Copy %zd/%zd bytes from base=%p to buf=%p\n",
		//	 c, (*iov)[0].iov_len, (void*)(*iov)[0].iov_base, (void*)buf));

		if (buf) memcpy(buf, (*iov)[0].iov_base, c);

		(*iov)[0].iov_len -= c;
		(*iov)[0].iov_base = (char *)(*iov)[0].iov_base + c;

		//DPRINTF(("iov %p is now %zd bytes at %p\n", (void *)*iov,
		//	 (*iov)[0].iov_len, (void *)(*iov)[0].iov_base));

		if ((*iov)[0].iov_len == 0) {
			(*iov)++;
			(*iov_len)--;
			//DPRINTF(("iov elem consumed, now iov=%p, iov_len=%d\n", (void *)*iov, *iov_len));
		}

		if (buf) buf = (char *)buf + c;
		//DPRINTF(("buf now %p\n", (void *)buf));

		res += c;
	}
	//DPRINTF(("iovec_pull pulled %zd/%zd bytes\n", res, bytes));

	return res;
}

static size_t iovec_push(struct iovec **iov, int *iov_len, void *buf, size_t bytes)
{
	size_t res = 0;

	//DPRINTF(("iovec_push %zd bytes from %p. iov=%p, iov_len=%d\n",
	//	 bytes, (void *)buf, (void *)*iov, *iov_len));

	while (res < bytes && *iov_len) {
		size_t c = (bytes - res) < (*iov)[0].iov_len ? (bytes - res) : (*iov)[0].iov_len;

		//DPRINTF(("Copy %zd/%zd bytes from buf=%p to base=%p\n",
		//	 c, (*iov)[0].iov_len, (void *)buf, (void *)(*iov)[0].iov_base));

		memcpy((*iov)[0].iov_base, buf, c);

		(*iov)[0].iov_len -= c;
		(*iov)[0].iov_base = (char *)(*iov)[0].iov_base + c;

		//DPRINTF(("iov %p is now %zd bytes at %p\n", (void *)*iov,
		//	 (*iov)[0].iov_len, (void *)(*iov)[0].iov_base));

		if ((*iov)[0].iov_len == 0) {
			(*iov)++;
			(*iov_len)--;
			//DPRINTF(("iov elem consumed, now iov=%p, iov_len=%d\n", (void *)*iov, *iov_len));
		}

		buf = (char *)buf + c;
		//DPRINTF(("buf now %p\n", (void *)buf));

		res += c;
	}

	return res;
}

static void dprint_iovec(struct iovec *iov, int iovec_len, const char *ctx)
{
	int i;
	DPRINTF(("%s: IOV:%p ELEMS:%d\n", ctx, (void *)iov, iovec_len));
	for (i = 0; i < iovec_len; i++)
		DPRINTF(("%s:  %d = %zu @ %p\n",
			 ctx, i, iov[i].iov_len, (void *)iov[i].iov_base));
}

static void dprint_chain(struct iovec *iov, int iovec_len, uint16_t *flags,
			 const char *ctx)
{
	int i;
	DPRINTF(("%s: CHAIN:%p ELEMS:%d\n", ctx, (void *)iov, iovec_len));
	for (i = 0; i < iovec_len; i++)
		DPRINTF(("%s:  %d = %zu @ %p (%"PRIx16")\n",
			 ctx, i, iov[i].iov_len, (void *)iov[i].iov_base, flags[i]));
}


static void dprint_header(struct virtio_sock_hdr *hdr, bool tx, const char *ctx)
{
	assert(hdr->op < nitems(opnames));

	DPRINTF(("%s: %sSRC:"PRIaddr" DST:"PRIaddr"\n",
		 ctx, tx ? "<=" : "=>",
		 hdr->src_cid, hdr->src_port, hdr->dst_cid, hdr->dst_port));
	DPRINTF(("%s:   LEN:%08"PRIx32" TYPE:%04"PRIx16" OP:%"PRId16"=%s\n",
		 ctx, hdr->len, hdr->type, hdr->op,
		 opnames[hdr->op] ? opnames[hdr->op] : "<unknown>"));
	DPRINTF(("%s:  FLAGS:%08"PRIx32" BUF_ALLOC:%08"PRIx32" FWD_CNT:%08"PRIx32"\n",
		 ctx, hdr->flags, hdr->buf_alloc, hdr->fwd_cnt));
}

static void put_sock(struct pci_vtsock_sock *s)
{
	int err = pthread_mutex_unlock(&s->mtx);
	assert(err == 0);
}

static struct pci_vtsock_sock *get_sock(struct pci_vtsock_sock *s)
{
	int err = pthread_mutex_lock(&s->mtx);
	assert(err == 0);
	return s;
}

/* Returns a _locked_, sock by idx */
static struct pci_vtsock_sock *lookup_sock_by_idx(struct pci_vtsock_softc *sc, int i)
{
	struct pci_vtsock_sock *s = &sc->socks[i];

	/*
	 * Avoid locking overhead if the socket is free. Since any
	 * alloc will trigger a kick of the rx or tx threads there is
	 * no danger of missing something which is being allocated
	 * right now.
	 *
	 * Since alloc_sock takes the sock->mtx before setting state
	 * we won't see a half constructed socket here either, since
	 * the caller of alloc_sock will complete the init before
	 * dropping the sock->mtx.
	 */
	if (s->state == SOCK_FREE) return NULL;
	return get_sock(s);
}

static struct pci_vtsock_sock *lookup_sock(struct pci_vtsock_softc *sc,
					   uint16_t type,
					   struct vsock_addr local_addr,
					   struct vsock_addr peer_addr)
{
	int i;

	assert(type == VIRTIO_VSOCK_TYPE_STREAM);

	for (i = 0 ; i < VTSOCK_MAXSOCKS; i++) {
		struct pci_vtsock_sock *s = lookup_sock_by_idx(sc, i);
		if (!s) continue;

		if ((s->state == SOCK_CONNECTED || s->state == SOCK_CONNECTING) &&
		    s->peer_addr.cid == peer_addr.cid &&
		    s->peer_addr.port == peer_addr.port &&
		    s->local_addr.cid == local_addr.cid &&
		    s->local_addr.port == local_addr.port) {
			return s;
		}

		put_sock(s);
	}

	return NULL;
}


/* Returns NULL on failure or a locked socket on success */
static struct pci_vtsock_sock *alloc_sock(struct pci_vtsock_softc *sc)
{
	struct pci_vtsock_sock *s = NULL; /* XXX init otherwise cc thinks return s uses undefined s! */
	int i;

	pthread_mutex_lock(&sc->alloc_mtx);
	for (i = 0 ; i < VTSOCK_MAXSOCKS; i++) {
		s = &sc->socks[i];
		if (s->state == SOCK_FREE) {
			get_sock(s);
			s->state = SOCK_CONNECTING;
			break;
		}
	}
	pthread_mutex_unlock(&sc->alloc_mtx);

	assert(s != NULL);
	assert(i == VTSOCK_MAXSOCKS || s->state == SOCK_CONNECTING);

	if (i == VTSOCK_MAXSOCKS)
		return NULL;

	s->buf_alloc = WRITE_BUF_LENGTH;
	s->fwd_cnt = 0;

	s->peer_buf_alloc = 0;
	s->peer_fwd_cnt = 0;
	s->rx_cnt = 0;
	s->credit_update_required = false;

	s->local_shutdown = 0;
	s->peer_shutdown = 0;

	s->write_buf_head = s->write_buf_tail = 0;

	return s;
}

static int set_socket_options(struct pci_vtsock_sock *s)
{
	int rc, buf_alloc = (int)s->buf_alloc;
	socklen_t opt_len;

	rc = setsockopt(s->fd, SOL_SOCKET, SO_SNDBUF,
			&buf_alloc, sizeof(buf_alloc));
	if ( rc < 0 ) {
		DPRINTF(("Failed to set SO_SNDBUF on fd %d: %s\n",
			 s->fd, strerror(errno)));
		return rc;
	}

	rc = setsockopt(s->fd, SOL_SOCKET, SO_RCVBUF,
			&buf_alloc, sizeof(buf_alloc));
	if ( rc < 0 ) {
		DPRINTF(("Failed to set SO_RCVBUF on fd %d: %s\n",
			 s->fd, strerror(errno)));
		return rc;
	}

	opt_len = sizeof(buf_alloc);
	rc = getsockopt(s->fd, SOL_SOCKET, SO_SNDBUF,
			&buf_alloc, &opt_len);
	if ( rc < 0 ) {
		DPRINTF(("Failed to get SO_SNDBUF on fd %d: %s\n",
			 s->fd, strerror(errno)));
		return rc;
	}
	/* If we didn't get what we asked for then expose this to the other end */
	if (buf_alloc < (int)s->buf_alloc) {
		PPRINTF(("fd %d SO_SNDBUF is 0x%x not 0x%x as requested, clamping\n",
			 s->fd, buf_alloc, s->buf_alloc));
		s->buf_alloc = (uint32_t)buf_alloc;
	}

	return 0;
}

static struct pci_vtsock_sock *connect_sock(struct pci_vtsock_softc *sc,
					    struct vsock_addr local_addr,
					    struct vsock_addr peer_addr,
					    uint32_t peer_buf_alloc,
					    uint32_t peer_fwd_cnt)
{
	struct pci_vtsock_sock *s;
	struct sockaddr_un un;
	int rc, fd = -1;

	s = alloc_sock(sc);
	if (s == NULL) {
		DPRINTF(("TX: No available socks\n"));
		goto err;
	}

	DPRINTF(("TX: Assigned sock %ld at %p\n",
		 s - &sc->socks[0], (void *)s));

	bzero(&un, sizeof(un));

	un.sun_len = 0; /* Unused? */
	un.sun_family = AF_UNIX;
	rc = snprintf(un.sun_path, sizeof(un.sun_path),
		     "%s/"PRIaddr, sc->path, FMTADDR(local_addr));
	if (rc < 0) {
		DPRINTF(("TX: Failed to format socket path\n"));
		goto err;
	}

	fd = socket(AF_UNIX, SOCK_STREAM, 0);
	if (fd < 0) {
		DPRINTF(("TX: socket failed for %s: %s\n",
			 un.sun_path, strerror(errno)));
		goto err;
	}

	rc = connect(fd, (struct sockaddr *)&un, sizeof(un));
	if (rc < 0) {
		DPRINTF(("TX: connect failed for %s: %s\n",
			 un.sun_path, strerror(errno)));
		goto err;
	}

	rc = fcntl(fd, F_SETFL, O_NONBLOCK);
	if (rc < 0) {
		DPRINTF(("TX: O_NONBLOCK failed for %s: %s\n",
			 un.sun_path, strerror(errno)));
		goto err;
	}

	DPRINTF(("TX: Socket path %s opened on fd %d\n", un.sun_path, fd));

	s->fd = fd;
	s->peer_addr = peer_addr;
	s->local_addr = local_addr;

	s->peer_buf_alloc = peer_buf_alloc;
	s->peer_fwd_cnt = peer_fwd_cnt;

	rc = set_socket_options(s);
	if (rc < 0) goto err;

	PPRINTF(("TX: SOCK connected (%d) "PRIaddr" <=> "PRIaddr"\n",
		 s->fd, FMTADDR(s->local_addr), FMTADDR(s->peer_addr)));
	s->state = SOCK_CONNECTED;

	put_sock(s);

	return s;

err:
	/* s is static, no need to free(), but do set state back to free */
	if (fd >= 0) close(fd);
	if (s) {
		s->state = SOCK_FREE;
		put_sock(s);
	}
	return NULL;
}

static void kick_rx(struct pci_vtsock_softc *sc, const char *why)
{
	char dummy;
	ssize_t nr;
	nr = write(sc->rx_kick_fd, &dummy, 1);
	assert(nr == 1);
	DPRINTF(("RX: kicked rx thread: %s\n", why));
}

static void kick_tx(struct pci_vtsock_softc *sc, const char *why)
{
	char dummy;
	ssize_t nr;
	nr = write(sc->tx_kick_fd, &dummy, 1);
	assert(nr == 1);
	DPRINTF(("TX: kicked tx thread: %s\n", why));
}

/* Reflect peer_shutdown into local fd */
static uint32_t shutdown_peer_local_fd(struct pci_vtsock_sock *s, uint32_t mode,
				       const char *ctx)
{
	int rc;
	int how;
	const char *how_str;
	uint32_t new = mode | s->peer_shutdown;
	uint32_t set = s->peer_shutdown ^ new;

	assert((mode & ~VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL) == 0);
	assert(mode != 0);

	DPRINTF(("%s: PEER CUR %"PRIx32", MODE %"PRIx32", NEW %"PRIx32", SET %"PRIx32"\n",
		 ctx, s->peer_shutdown, mode, new, set));

	switch (set) {
	case 0:
		return 0;
	case VIRTIO_VSOCK_FLAG_SHUTDOWN_TX:
		how = SHUT_WR;
		how_str = "SHUT_WR";
		break;
	case VIRTIO_VSOCK_FLAG_SHUTDOWN_RX:
		how = SHUT_RD;
		how_str = "SHUT_RD";
		break;
	case VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL:
		how = SHUT_RDWR;
		how_str = "SHUT_RDWR";
		break;
	default:
		abort();
	}

	rc = shutdown(s->fd, how);
	DPRINTF(("%s: shutdown_peer: shutdown(%d, %s)\n", ctx, s->fd, how_str));
	if (rc < 0 && errno != ENOTCONN) {
		DPRINTF(("%s: shutdown(%d, %s) for peer shutdown failed: %s\n",
			 ctx, s->fd, how_str, strerror(errno)));
		abort();
	}

	s->peer_shutdown = new;
	return set;
}

/* The caller must have sent something (probably OP_RST, but perhaps
 * OP_SHUTDOWN) to the peer already.
 */
static void close_sock(struct pci_vtsock_softc *sc,  struct pci_vtsock_sock *s,
		       const char *ctx)
{
	if (!s) return;
	DPRINTF(("%s: Closing sock %p\n", ctx, (void *)s));

	shutdown_peer_local_fd(s, VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL, ctx);

	/* The call to peer_local_fd will have done any required
	 * shutdown() call on s->fd */
	s->local_shutdown = VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL;

	s->state = SOCK_CLOSING;
	kick_rx(sc, "sock closed");
}

static void shutdown_peer_sock(struct pci_vtsock_softc *sc, struct pci_vtsock_sock *s,
			       uint32_t mode, const char *ctx)
{
	bool kick = false;
	uint32_t set;

	assert((mode & ~VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL) == 0);

	if (s->state != SOCK_CONNECTED) goto done;

	assert(s->local_shutdown != VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL);

	DPRINTF(("%s: fd %d: PEER SHUTDOWN 0x%"PRIx32" (0x%"PRIx32")\n",
		 ctx, s->fd, mode, s->peer_shutdown));

	set = shutdown_peer_local_fd(s, mode, ctx);

	DPRINTF(("%s: setting 0x%"PRIx32" mode is now 0x%"PRIx32" (local 0x%"PRIx32")\n",
		 ctx, set, s->peer_shutdown, s->local_shutdown));

	/* Is everything now closed? */
	if (/*set &&*/ mode == VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL &&
	    /*s->local_shutdown == VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL &&*/
	    s->peer_shutdown == VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL) {
		s->state = SOCK_CLOSING;
		kick = true;
	}

done:
	if (kick) kick_rx(sc, "shutdown (peer)");
}

/* Caller should send OP_SHUTDOWN with flags == s->local_shutdown after calling this */
static void shutdown_local_sock(struct pci_vtsock_softc *sc, struct pci_vtsock_sock *s,
				uint32_t mode, const char *ctx)
{
	bool kick = false;
	uint32_t new, set;

	assert((mode & ~VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL) == 0);

	if (s->state != SOCK_CONNECTED) goto done;

	assert(s->local_shutdown != VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL);

	DPRINTF(("%s: fd %d: LOCAL SHUTDOWN 0x%"PRIx32" (0x%"PRIx32")\n",
		 ctx, s->fd, mode, s->peer_shutdown));

	new = mode | s->local_shutdown;
	set = s->local_shutdown ^ new;
	s->local_shutdown = new;

	DPRINTF(("%s: setting 0x%"PRIx32" mode is now 0x%"PRIx32" (peer 0x%"PRIx32")\n",
		 ctx, set, s->local_shutdown, s->peer_shutdown));

	/* Did we do something and is everything now closed? */
	if (set &&
	    s->local_shutdown == VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL /*&&
	    s->peer_shutdown == VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL*/) {
		s->state = SOCK_CLOSING;
		kick = true;
	}

done:
	if (kick) kick_rx(sc, "shutdown (local)");
}

static void set_credit_update_required(struct pci_vtsock_softc *sc,
				       struct pci_vtsock_sock *sock)
{
	if (sock->credit_update_required) return;
	sock->credit_update_required = true;
	kick_rx(sc, "credit update required");
}

static void send_response_common(struct pci_vtsock_softc *sc,
				 struct vsock_addr local_addr,
				 struct vsock_addr peer_addr,
				 uint16_t op, uint16_t type, uint32_t flags,
				 uint32_t buf_alloc, uint32_t fwd_cnt)
{
	struct virtio_sock_hdr *hdr;
	int slot;

	assert(op != VIRTIO_VSOCK_OP_RW);
	assert(flags == 0 || op == VIRTIO_VSOCK_OP_SHUTDOWN);

	pthread_mutex_lock(&sc->reply_mtx);

	slot = sc->reply_prod++;
	if (sc->reply_prod == VTSOCK_REPLYRINGSZ)
		sc->reply_prod = 0;
	/* check for reply ring overflow */
	/* XXX correct check? */
	DPRINTF(("TX: QUEUING REPLY IN SLOT %x (prod %x, cons %x)\n",
		 slot, sc->reply_prod, sc->reply_cons));
	assert(sc->reply_cons != sc->reply_prod);

	hdr = &sc->reply_ring[slot];

	hdr->src_cid = local_addr.cid;
	hdr->src_port = local_addr.port;

	hdr->dst_cid = peer_addr.cid;
	hdr->dst_port = peer_addr.port;

	hdr->len = 0;
	hdr->type = type;
	hdr->op = op;
	hdr->flags = flags;

	hdr->buf_alloc = buf_alloc;
	hdr->fwd_cnt = fwd_cnt;

	dprint_header(hdr, 0, "TX");

	kick_rx(sc, "tx thread queued response");
	pthread_mutex_unlock(&sc->reply_mtx);
}

static void send_response_sock(struct pci_vtsock_softc *sc,
				 uint16_t op, uint32_t flags,
				 const struct pci_vtsock_sock *sock)
{
	send_response_common(sc, sock->local_addr, sock->peer_addr,
			     op, VIRTIO_VSOCK_TYPE_STREAM, flags,
			     sock->buf_alloc, sock->fwd_cnt);
}

static void send_response_nosock(struct pci_vtsock_softc *sc, uint16_t op,
				 uint16_t type,
				 struct vsock_addr local_addr,
				 struct vsock_addr peer_addr)
{
	send_response_common(sc, local_addr, peer_addr,
			     op, type, 0, 0, 0);
}

static bool sock_is_buffering(struct pci_vtsock_sock *sock)
{
	return sock->write_buf_tail > 0;
}

static int buffer_write(struct pci_vtsock_sock *sock,
			uint32_t len, struct iovec *iov, int iov_len)
{
	size_t nr;
	if (sock->write_buf_tail + len > WRITE_BUF_LENGTH) {
		DPRINTF(("TX: fd %d unable to buffer write of 0x%"PRIx32" bytes,"
			 " buffer use 0x%x/0x%x, 0x%x remaining\n",
			 sock->fd, len, sock->write_buf_tail,
			 WRITE_BUF_LENGTH,
			 WRITE_BUF_LENGTH < sock->write_buf_tail));
		return -1;
	}

	nr = iovec_pull(&iov, &iov_len,
			&sock->write_buf[sock->write_buf_tail], len);
	assert(nr == len);
	assert(iov_len == 0);

	sock->write_buf_tail += nr;
	DPRINTF(("TX: fd %d buffered 0x%"PRIx32" bytes (0x%x/0x%x)\n",
		 sock->fd, len, sock->write_buf_tail, WRITE_BUF_LENGTH));

	return 0;
}

static void buffer_drain(struct pci_vtsock_softc *sc,
			 struct pci_vtsock_sock *sock)
{
	struct timeval before, after;
	ssize_t nr;

	DPRINTF(("TX: buffer drain on fd %d 0x%x-0x%x/0x%x\n",
		 sock->fd, sock->write_buf_head, sock->write_buf_tail,
		 WRITE_BUF_LENGTH));

	assert(sock_is_buffering(sock));
	assert(sock->write_buf_head < sock->write_buf_tail);

	gettimeofday(&before, NULL);
	nr = write(sock->fd, &sock->write_buf[sock->write_buf_head],
		   sock->write_buf_tail - sock->write_buf_head);
	gettimeofday(&after, NULL);
	if ((after.tv_sec - before.tv_sec) > 5)
		fprintf(stderr, "TX: WARNING: write on fd %d took %ld seconds\n",
			sock->fd, (unsigned long)(after.tv_sec - before.tv_sec));
	if (nr == -1) {
		if (errno == EPIPE) {
			/* Assume EOF and shutdown */
			PPRINTF(("TX: writev fd=%d failed with EPIPE => SHUTDOWN_RX\n", sock->fd));
			shutdown_local_sock(sc, sock, VIRTIO_VSOCK_FLAG_SHUTDOWN_RX, "RX");
			send_response_sock(sc, VIRTIO_VSOCK_OP_SHUTDOWN,
					   sock->local_shutdown, sock);
			return;
		} else if (errno == EAGAIN) {
			return;
		} else {
			PPRINTF(("TX: write fd=%d failed with %d %s\n", sock->fd,
				 errno, strerror(errno)));
			send_response_sock(sc, VIRTIO_VSOCK_OP_RST, 0, sock);
			close_sock(sc, sock, "TX");
			return;
		}
	}

	DPRINTF(("TX: drained %zd/%"PRId32" bytes in %ld seconds\n", nr,
		 sock->write_buf_tail - sock->write_buf_head,
		 (unsigned long)(after.tv_sec - before.tv_sec)));
	sock->write_buf_head += nr;
	if (sock->write_buf_head < sock->write_buf_tail)
		return;

	/* Buffer completely drained, reset and update peer.  NB: We
	 * only update fwd_cnt once the buffer is empty rather than as
	 * we go, in the hopes that we then won't need to buffer so
	 * much as we go on.
	 */
	DPRINTF(("TX: fd %d buffer drained of 0x%x bytes\n",
		 sock->fd, sock->write_buf_head));
	sock->fwd_cnt += sock->write_buf_head;
	sock->write_buf_head = sock->write_buf_tail = 0;
	set_credit_update_required(sc, sock);
}

/* -> 1 == success, update peer credit
 * -> 0 == success, don't update peer credit
 */
static int handle_write(struct pci_vtsock_softc *sc,
			struct pci_vtsock_sock *sock,
			uint32_t len, struct iovec *iov, int iov_len)
{
	struct timeval before, after;
	ssize_t num;

	if (sock_is_buffering(sock)) {
		return buffer_write(sock, len, iov, iov_len);
	}

	gettimeofday(&before, NULL);
	num = writev(sock->fd, iov, iov_len);
	gettimeofday(&after, NULL);
	if ((after.tv_sec - before.tv_sec) > 5)
		fprintf(stderr, "TX: WARNING: writev on fd %d took %ld seconds\n",
			sock->fd, (unsigned long)(after.tv_sec - before.tv_sec));
	if (num == -1) {
		if (errno == EPIPE) {
			/* Assume EOF and shutdown */
			PPRINTF(("TX: writev fd=%d failed with EPIPE => SHUTDOWN_RX\n", sock->fd));
			shutdown_local_sock(sc, sock, VIRTIO_VSOCK_FLAG_SHUTDOWN_RX, "RX");
			send_response_sock(sc, VIRTIO_VSOCK_OP_SHUTDOWN,
					   sock->local_shutdown, sock);
			return 0;
		} else if (errno == EAGAIN) {
			num = 0;
		} else {
			PPRINTF(("TX: writev fd=%d failed with %d %s\n", sock->fd,
				 errno, strerror(errno)));
			return -1;
		}
	}

	DPRINTF(("TX: wrote %zd/%"PRId32" bytes in %ld seconds\n", num, len,
		 (unsigned long)(after.tv_sec - before.tv_sec)));
	if (num == len) {
		sock->fwd_cnt += num;
		return 1;
	} else { /* Buffer the rest */
		size_t pulled = iovec_pull(&iov, &iov_len, NULL, (size_t)num);
		assert(pulled == (size_t)num);
		return buffer_write(sock, len - (uint32_t)num, iov, iov_len);
	}
}

static void pci_vtsock_proc_tx(struct pci_vtsock_softc *sc,
			       struct vqueue_info *vq)
{
	struct pci_vtsock_sock *sock;
	struct iovec iov_array[VTSOCK_MAXSEGS], *iov = iov_array;
	uint16_t idx, flags[VTSOCK_MAXSEGS];
	struct virtio_sock_hdr hdr;
	int iovec_len;
	size_t pulled;

	iovec_len = vq_getchain(vq, &idx, iov, VTSOCK_MAXSEGS, flags);
	assert(iovec_len <= VTSOCK_MAXSEGS);

	DPRINTF(("TX: chain with %d buffers at idx %"PRIx16"\n",
		 iovec_len, idx));
	dprint_chain(iov, iovec_len, flags, "TX");
	//assert(iov[0].iov_len >= sizeof(*hdr));
	//hdr = iov[0].iov_base;

	pulled = iovec_pull(&iov, &iovec_len, &hdr, sizeof(hdr));
	assert(pulled == sizeof(hdr));

	dprint_header(&hdr, 1, "TX");

	dprint_iovec(iov, iovec_len, "TX");

	if (hdr.src_cid != sc->vssc_cfg.guest_cid ||
	    hdr.dst_cid != VMADDR_CID_HOST ||
	    hdr.type != VIRTIO_VSOCK_TYPE_STREAM) {
		DPRINTF(("TX: Bad src/dst address/type\n"));
		send_response_nosock(sc, VIRTIO_VSOCK_OP_RST,
				     hdr.type,
				     (struct vsock_addr) {
					     .cid = hdr.dst_cid,
					     .port =hdr.dst_port
				     },
				     (struct vsock_addr) {
					     .cid = hdr.src_cid,
					     .port =hdr.src_port
				     });
		vq_relchain(vq, idx, 0);
		return;
	}

	sock = lookup_sock(sc, VIRTIO_VSOCK_TYPE_STREAM,
			   (struct vsock_addr) {
				   .cid = hdr.dst_cid,
					   .port =hdr.dst_port
			   },
			   (struct vsock_addr) {
				   .cid = hdr.src_cid,
					   .port =hdr.src_port
			   });

	if (sock) {
		sock->peer_buf_alloc = hdr.buf_alloc;
		sock->peer_fwd_cnt = hdr.fwd_cnt;
	}

	switch (hdr.op) {
	case VIRTIO_VSOCK_OP_INVALID:
		PPRINTF(("TX: => INVALID\n"));
		goto do_rst;

	case VIRTIO_VSOCK_OP_REQUEST:
		/* Attempt to (re)connect existing sock? Naughty! */
		/* Or is it -- what are the semantics? */
		if (sock) {
			PPRINTF(("TX: Attempt to reconnect sock\n"));
			goto do_rst;
		}

		if (hdr.dst_cid == sc->vssc_cfg.guest_cid) {
			PPRINTF(("TX: Attempt to connect back to guest\n!"));
			goto do_rst;
		}

		sock = connect_sock(sc,
				    (struct vsock_addr){
					    .cid = hdr.dst_cid, .port = hdr.dst_port
				    },
				    (struct vsock_addr){
					    .cid = hdr.src_cid, .port = hdr.src_port
				    }, hdr.buf_alloc, hdr.fwd_cnt);
		if (!sock) {
			PPRINTF(("TX: Failed to open sock\n"));
			goto do_rst;
		}

		send_response_sock(sc, VIRTIO_VSOCK_OP_RESPONSE, 0, sock);
		vq_relchain(vq, idx, 0);
		/* No rx kick required, send_response_sock did one */
		break;

	case VIRTIO_VSOCK_OP_RESPONSE:
		if (!sock) {
			PPRINTF(("TX: RESPONSE to non-existent sock\n"));
			goto do_rst;
		}
		if (sock->state != SOCK_CONNECTING) {
			PPRINTF(("TX: RESPONSE to non-connecting sock (state %d)\n",
				 sock->state));
			goto do_rst;
		}
		PPRINTF(("TX: SOCK connected (%d) "PRIaddr" <=> "PRIaddr"\n",
			 sock->fd, FMTADDR(sock->local_addr), FMTADDR(sock->peer_addr)));
		sock->state = SOCK_CONNECTED;
		vq_relchain(vq, idx, 0);
		kick_rx(sc, "new outgoing sock");
		break;

	case VIRTIO_VSOCK_OP_RST:
		/* No response */
		if (!sock)
			PPRINTF(("TX: RST to non-existent sock\n"));
		close_sock(sc, sock, "TX");
		vq_relchain(vq, idx, 0);
		break;

	case VIRTIO_VSOCK_OP_SHUTDOWN:
		if (!sock) {
			PPRINTF(("TX: SHUTDOWN to non-existent sock\n"));
			goto do_rst;
		}
		if (sock->state != SOCK_CONNECTED) {
			PPRINTF(("TX: SHUTDOWN to non-connected sock (state %d)\n",
				 sock->state));
			goto do_rst;
		}
		if (hdr.flags & ~VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL) {
			PPRINTF(("TX: SHUTDOWN with reserved flags %"PRIx32"\n",
				 hdr.flags));
			goto do_rst; /* ??? */
		}
		if (!(hdr.flags & VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL)) {
			PPRINTF(("TX: SHUTDOWN with no flags %"PRIx32"\n",
				 hdr.flags));
			goto do_rst; /* ??? */
		}

		shutdown_peer_sock(sc, sock, hdr.flags, "TX");

		vq_relchain(vq, idx, 0);
		break;

	case VIRTIO_VSOCK_OP_RW:
	{
		int rc;

		if (!sock) {
			PPRINTF(("TX: RW with no sock\n"));
			goto do_rst;
		}
		if (sock->state != SOCK_CONNECTED) {
			PPRINTF(("TX: RW to non-connected sock (state %d)\n",
				 sock->state));
			goto do_rst;
		}
		if (sock->peer_shutdown & VIRTIO_VSOCK_FLAG_SHUTDOWN_TX) {
			PPRINTF(("TX: RW to socket with peer_shutdown.TX\n"));
			goto do_rst;
		}
		if (sock->local_shutdown & VIRTIO_VSOCK_FLAG_SHUTDOWN_RX) {
			PPRINTF(("TX: RW to socket with local_shutdown.RX\n"));
			goto do_rst;
		}
		rc = handle_write(sc, sock, hdr.len, iov, iovec_len);
		if (rc < 0) goto do_rst;
		vq_relchain(vq, idx, 0);
		if (rc == 1)
			set_credit_update_required(sc, sock);
		break;
	}

	case VIRTIO_VSOCK_OP_CREDIT_UPDATE:
		if (!sock) {
			PPRINTF(("TX: CREDIT_UPDATE to non-existent sock\n"));
			goto do_rst;
		}
		if (sock->state != SOCK_CONNECTED) {
			PPRINTF(("TX: CREDIT_UPDATE to non-connected sock (state %d)\n",
				 sock->state));
			goto do_rst;
		}
		/* No response needed, we updated above */
		vq_relchain(vq, idx, 0);
		/* But kick rx thread to attempt to send more */
		kick_rx(sc, "credit update");
		break;

	case VIRTIO_VSOCK_OP_CREDIT_REQUEST:
		if (!sock) {
			PPRINTF(("TX: CREDIT_REQUEST to non-existent sock\n"));
			goto do_rst;
		}
		if (sock->state != SOCK_CONNECTED) {
			PPRINTF(("TX: CREDIT_REQUEST to non-connected sock (state %d)\n",
				 sock->state));
			goto do_rst;
		}
		vq_relchain(vq, idx, 0);
		set_credit_update_required(sc, sock);
		break;
	}

	if (sock)
		put_sock(sock);

	return;

do_rst:
	if (sock)
		send_response_sock(sc, VIRTIO_VSOCK_OP_RST, 0, sock);
	else
		send_response_nosock(sc, VIRTIO_VSOCK_OP_RST, hdr.type,
				     (struct vsock_addr) {
					     .cid = hdr.dst_cid,
					     .port =hdr.dst_port
				     },
				     (struct vsock_addr) {
					     .cid = hdr.src_cid,
					     .port =hdr.src_port
				     });
	vq_relchain(vq, idx, 0);
	close_sock(sc, sock, "TX");
	if (sock) put_sock(sock);
	return;
}

static void handle_connect_fd(struct pci_vtsock_softc *sc, int accept_fd, uint32_t cid, uint32_t port)
{
	int fd, rc;
	char buf[8 + 1 + 8 + 1 + 1]; /* %08x.%08x\n\0 */
	ssize_t bytes;
	struct pci_vtsock_sock *sock = NULL;

	fd = accept(accept_fd, NULL, NULL);
	if (fd < 0) {
		fprintf(stderr,
			"TX: Unable to accept incoming connection: %d (%s)\n",
			errno, strerror(errno));
		return;
	}

	DPRINTF(("TX: Connect attempt on connect fd => %d\n", fd));

	if (cid == VMADDR_CID_ANY) {
		do {
			bytes = read(fd, buf, sizeof(buf)-1);
		} while (bytes == -1 && errno == EAGAIN);

		if (bytes != sizeof(buf) - 1) {
			DPRINTF(("TX: Short read on connect %zd/%zd\n", bytes, sizeof(buf)-1));
			if (bytes == -1) DPRINTF(("TX: errno: %s\n", strerror(errno)));
			goto err;
		}
		buf[sizeof(buf)-1] = '\0';

		if (buf[sizeof(buf)-2] != '\n') {
			DPRINTF(("TX: No newline on connect %s\n", buf));
			goto err;
		}

		DPRINTF(("TX: Connect to %s", buf));

		rc = sscanf(buf, "%08x.%08x\n", &cid, &port);
		if (rc != 2) {
			DPRINTF(("TX: Failed to parse connect attempt\n"));
			goto err;
		}
		DPRINTF(("TX: Connection requested to %08x.%08x\n", cid, port));
	} else {
		DPRINTF(("TX: Forwarding connection to %08x.%08x\n", cid, port));
	}

	if (cid != sc->vssc_cfg.guest_cid) {
		DPRINTF(("TX: Attempt to connect to non-guest CID\n"));
		goto err;
	}

	sock = alloc_sock(sc);

	if (sock == NULL) {
		DPRINTF(("TX: No available sockets for connect\n"));
		goto err;
	}

	DPRINTF(("TX: Assigned sock %ld at %p for connect\n",
		 sock - &sc->socks[0], (void *)sock));

	sock->fd = fd;
	sock->peer_addr.cid = cid;
	sock->peer_addr.port = port;
	sock->local_addr.cid = VMADDR_CID_HOST;
	/* Start at 2^16 to be larger than a TCP port
	 * XXX Allocate properly.
         */
	sock->local_addr.port = 65536U + (uint32_t)(sock - &sc->socks[0]);

	rc = set_socket_options(sock);
	if (rc < 0) goto err;

	put_sock(sock);

	PPRINTF(("TX: SOCK connecting (%d) "PRIaddr" <=> "PRIaddr"\n",
		 sock->fd, FMTADDR(sock->local_addr), FMTADDR(sock->peer_addr)));
	send_response_sock(sc, VIRTIO_VSOCK_OP_REQUEST, 0, sock);

	return;
err:
	if (sock) {
		sock->state = SOCK_FREE;
		put_sock(sock);
	}
	close(fd);
}

static void *pci_vtsock_tx_thread(void *vsc)
{
	struct pci_vtsock_softc *sc = vsc;
	struct vqueue_info *vq = &sc->vssc_vqs[VTSOCK_QUEUE_TX];
	fd_set rfd, wfd;

	assert(sc);
	assert(sc->tx_wake_fd != -1);
	assert(sc->connect_fd != -1);

	while(1) {
		int i, nrfd, maxfd, nr;
		int buffering = 0;

		FD_ZERO(&rfd);
		FD_ZERO(&wfd);

		FD_SET(sc->tx_wake_fd, &rfd);
		maxfd = sc->tx_wake_fd;

		FD_SET(sc->connect_fd, &rfd);
		maxfd = max_fd(sc->connect_fd, maxfd);
		nrfd = 2;

		for (i = 0; i < sc->nr_fwds; i++) {
			struct pci_vtsock_forward *fwd = &sc->fwds[i];
			assert(fwd->listen_fd != -1);
			FD_SET(fwd->listen_fd, &rfd);
			maxfd = max_fd(fwd->listen_fd, maxfd);
			nrfd++;
		}

		for(i = 0; i < VTSOCK_MAXSOCKS; i++) {
			struct pci_vtsock_sock *s = lookup_sock_by_idx(sc, i);
			if (!s) continue;
			if (s->state != SOCK_CONNECTED) {
				put_sock(s);
				continue;
			}
			assert(s->fd >= 0);
			assert(s->fd < FD_SETSIZE);
			if (sock_is_buffering(s)) {
				FD_SET(s->fd, &wfd);
				maxfd = max_fd(s->fd, maxfd);
				buffering++;
				nrfd++;
			}
			put_sock(s);
		}

		DPRINTF(("TX: *** selecting on %d fds (buffering: %d)\n",
			 nrfd, buffering));
		nr = select(maxfd + 1, &rfd, &wfd, NULL, NULL);
		if (nr < 0) DPRINTF(("TX select returned %zd errno %d\n", nr, errno));
		assert(nr >= 0);
		DPRINTF(("TX:\nTX: *** %d/%d fds are readable/writeable\n", nr, nrfd));

		if (FD_ISSET(sc->tx_wake_fd, &rfd)) {
			/* Eat the notification(s) */
			char dummy[128];
			ssize_t rd_dummy = read(sc->tx_wake_fd, &dummy, sizeof(dummy));
			assert(rd_dummy >= 1);
			/* Restart select now that we have some descriptors */
			DPRINTF(("TX: thread got %zd kicks (have descs: %s)\n",
				 rd_dummy,
				 vq_has_descs(vq) ? "yes" : "no"));
		}

		if (FD_ISSET(sc->connect_fd, &rfd)) {
			DPRINTF(("TX: Handling connect fd\n"));
			handle_connect_fd(sc, sc->connect_fd, VMADDR_CID_ANY, 0);
		}

		for (i = 0; i < sc->nr_fwds; i++) {
			struct pci_vtsock_forward *fwd = &sc->fwds[i];
			if (FD_ISSET(fwd->listen_fd, &rfd)) {
				DPRINTF(("Attempt to connect to forwarded guest port %"PRId32"\n", fwd->port));
				handle_connect_fd(sc, fwd->listen_fd, sc->vssc_cfg.guest_cid, fwd->port);
			}
		}

		if (buffering) {
			for(i = 0; i < VTSOCK_MAXSOCKS; i++) {
				struct pci_vtsock_sock *s = lookup_sock_by_idx(sc, i);
				if (!s) continue;
				if (s->state != SOCK_CONNECTED) {
					put_sock(s);
					continue;
				}
				if (FD_ISSET(s->fd, &wfd)) {
					buffer_drain(sc, s);
				}
				put_sock(s);
			}
		}

		pthread_mutex_lock(&sc->tx_mtx);

		while (vq_has_descs(vq))
			pci_vtsock_proc_tx(sc, vq);

		if (vq_ring_ready(vq))
			vq_endchains(vq, 1);

		pthread_mutex_unlock(&sc->tx_mtx);

		DPRINTF(("TX: All work complete\n"));
	}
}

static void pci_vtsock_notify_tx(void *vsc, struct vqueue_info *vq)
{
	struct pci_vtsock_softc *sc = vsc;

	assert(vq == &sc->vssc_vqs[VTSOCK_QUEUE_TX]);
	kick_tx(sc, "notify");
}

/*
 * Returns:
 *  -1 == no descriptors available
 *   0 == nothing done (sock has shutdown, peer has no buffers, nothing on Unix socket)
 *  >0 == number of bytes read
 */
static ssize_t pci_vtsock_proc_rx(struct pci_vtsock_softc *sc,
				  struct vqueue_info *vq,
				  struct pci_vtsock_sock *s)
{
	struct virtio_sock_hdr *hdr;
	struct iovec iov_array[VTSOCK_MAXSEGS], *iov = iov_array;
	uint16_t flags[VTSOCK_MAXSEGS];
	uint16_t idx;
	uint32_t peer_free;
	int iovec_len;
	size_t pushed;
	ssize_t len;
	struct timeval before, after;

	assert(s->fd >= 0);

	if (!vq_has_descs(vq)) {
		DPRINTF(("RX: no queues!\n"));
		return -1;
	}

	peer_free = s->peer_buf_alloc - (s->rx_cnt - s->peer_fwd_cnt);
	DPRINTF(("RX:\tpeer free = %"PRIx32"\n", peer_free));
	if (!peer_free) return 0; /* No space */

	iovec_len = vq_getchain(vq, &idx, iov, VTSOCK_MAXSEGS, flags);
	DPRINTF(("RX: virtio-vsock: got %d elem rx chain\n", iovec_len));
	dprint_chain(iov, iovec_len, flags, "RX");

	assert(iovec_len >= 1);
	/* XXX needed so we can update len after the read */
	assert(iov[0].iov_len >= sizeof(*hdr));

	hdr = iov[0].iov_base;
	hdr->src_cid = s->local_addr.cid;
	hdr->src_port = s->local_addr.port;
	hdr->dst_cid = s->peer_addr.cid;
	hdr->dst_port = s->peer_addr.port;
	hdr->len = 0; /* XXX */
	hdr->type = VIRTIO_VSOCK_TYPE_STREAM;
	hdr->op = VIRTIO_VSOCK_OP_RW;
	hdr->flags = 0;
	hdr->buf_alloc = s->buf_alloc;
	hdr->fwd_cnt = s->fwd_cnt;

	pushed = iovec_push(&iov, &iovec_len, hdr, sizeof(*hdr));
	assert(pushed == sizeof(*hdr));

	iovec_clip(&iov, &iovec_len, peer_free);

	gettimeofday(&before, NULL);
	len = readv(s->fd, iov, iovec_len);
	gettimeofday(&after, NULL);
	if ((after.tv_sec - before.tv_sec) > 5)
		fprintf(stderr, "RX: WARNING: readv on fd %d took %ld seconds\n",
			s->fd, (unsigned long)(after.tv_sec - before.tv_sec));
	if (len == -1) {
		if (errno == EAGAIN) { /* Nothing to read/would block */
			DPRINTF(("RX: readv fd=%d EAGAIN in %ld seconds\n", s->fd,
				 (unsigned long)(after.tv_sec - before.tv_sec)));
			vq_retchain(vq);
			return 0;
		}
		PPRINTF(("RX: readv fd=%d failed with %d %s in %ld seconds\n",
			 s->fd, errno, strerror(errno),
			 (unsigned long)(after.tv_sec - before.tv_sec)));
		hdr->op = VIRTIO_VSOCK_OP_RST;
		hdr->flags = 0;
		hdr->len = 0;
		dprint_header(hdr, 0, "RX");
		vq_relchain(vq, idx, sizeof(*hdr));
		close_sock(sc, s, "RX");
		return 0;
	}
	DPRINTF(("RX: readv put %zd bytes into iov in %ld seconds\n",
		 len,
		 (unsigned long)(after.tv_sec - before.tv_sec)));
	if (len == 0) { /* Not actually anything to read -- EOF */
		PPRINTF(("RX: readv fd=%d EOF => SHUTDOWN_TX\n", s->fd));
		shutdown_local_sock(sc, s, VIRTIO_VSOCK_FLAG_SHUTDOWN_TX, "RX");
		hdr->op = VIRTIO_VSOCK_OP_SHUTDOWN;
		hdr->flags = s->local_shutdown;
		hdr->len = 0;
		dprint_header(hdr, 0, "RX");
		vq_relchain(vq, idx, sizeof(*hdr));
		return 0;
	}
	hdr->len = (uint32_t)len;

	s->rx_cnt += len;

	dprint_header(hdr, 0, "RX");

	vq_relchain(vq, idx, sizeof(*hdr) + (uint32_t)len);

	return len;
}

/* True if there is more to do */
static bool rx_do_one_reply(struct pci_vtsock_softc *sc,
			    struct vqueue_info *vq)
{
	struct virtio_sock_hdr *hdr;
	struct iovec iov_array[VTSOCK_MAXSEGS], *iov = iov_array;
	int iovec_len;
	uint16_t idx;
	size_t pushed;
	int slot;
	bool more_to_do = false;

	if (sc->reply_cons == sc->reply_prod)
		goto done;

	slot = sc->reply_cons++;
	if (sc->reply_cons == VTSOCK_REPLYRINGSZ)
		sc->reply_cons = 0;

	hdr = &sc->reply_ring[slot];

	iovec_len = vq_getchain(vq, &idx, iov, VTSOCK_MAXSEGS, NULL);
	DPRINTF(("RX: reply: got %d elem rx chain for slot %x (prod %x, cons %x)\n",
		 iovec_len, slot, sc->reply_prod, sc->reply_cons));

	assert(iovec_len >= 1);

	pushed = iovec_push(&iov, &iovec_len, hdr, sizeof(*hdr));
	assert(pushed == sizeof(*hdr));

	vq_relchain(vq, idx, sizeof(*hdr));

	more_to_do = sc->reply_cons != sc->reply_prod;

done:
	return more_to_do;
}

/* true on success, false if no descriptors */
static bool send_credit_update(struct vqueue_info *vq,
			       struct pci_vtsock_sock *s)
{
	struct virtio_sock_hdr *hdr;
	struct iovec iov_array[VTSOCK_MAXSEGS], *iov = iov_array;
	uint16_t idx;
	int iovec_len;

	assert(s->fd >= 0);

	if (!vq_has_descs(vq)) {
		DPRINTF(("RX: no queues for credit update!\n"));
		return false;
	}

	iovec_len = vq_getchain(vq, &idx, iov, VTSOCK_MAXSEGS, NULL);
	DPRINTF(("RX: virtio-vsock: got %d elem rx chain for credit update\n", iovec_len));
	dprint_chain(iov, iovec_len, NULL, "RX");

	assert(iovec_len >= 1);
	assert(iov[0].iov_len >= sizeof(*hdr));

	hdr = iov[0].iov_base;
	hdr->src_cid = s->local_addr.cid;
	hdr->src_port = s->local_addr.port;
	hdr->dst_cid = s->peer_addr.cid;
	hdr->dst_port = s->peer_addr.port;
	hdr->len = 0;
	hdr->type = VIRTIO_VSOCK_TYPE_STREAM;
	hdr->op = VIRTIO_VSOCK_OP_CREDIT_UPDATE;
	hdr->flags = 0;
	hdr->buf_alloc = s->buf_alloc;
	hdr->fwd_cnt = s->fwd_cnt;

	dprint_header(hdr, 0, "RX");

	vq_relchain(vq, idx, sizeof(*hdr));

	return true;
}

static void *pci_vtsock_rx_thread(void *vsc)
{
	struct pci_vtsock_softc *sc = vsc;
	struct vqueue_info *vq = &sc->vssc_vqs[VTSOCK_QUEUE_RX];
	fd_set rfd;
	bool poll_socks = true;

	assert(sc);
	assert(sc->rx_wake_fd != -1);

	while (1) {
		int nrfd, maxfd, i, nr;
		bool did_some_work = true;

		FD_ZERO(&rfd);


		FD_SET(sc->rx_wake_fd, &rfd);
		maxfd = sc->rx_wake_fd;
		nrfd = 1;

		if (poll_socks) {
			for(i = 0; i < VTSOCK_MAXSOCKS; i++) {
				struct pci_vtsock_sock *s = lookup_sock_by_idx(sc, i);
				uint32_t peer_free;
				if (!s) continue;
				if (s->state != SOCK_CONNECTED) {
					put_sock(s);
					continue;
				}
				if (s->local_shutdown & VIRTIO_VSOCK_FLAG_SHUTDOWN_TX) {
					put_sock(s);
					continue;
				}
				if (s->peer_shutdown & VIRTIO_VSOCK_FLAG_SHUTDOWN_RX) {
					put_sock(s);
					continue;
				}
				assert(s->fd >= 0);
				assert(s->fd < FD_SETSIZE);
				peer_free = s->peer_buf_alloc - (s->rx_cnt - s->peer_fwd_cnt);
				DPRINTF(("RX: sock %d (%d): peer free = %"PRId32"\n",
					 i, s->fd, peer_free));
				if (peer_free == 0) {
					put_sock(s);
					continue;
				}
				FD_SET(s->fd, &rfd);
				maxfd = max_fd(s->fd, maxfd);
				nrfd++;
				put_sock(s);
			}
		}

		/* Unlocked during select */

		DPRINTF(("RX: *** thread selecting on %d fds (socks: %s)\n",
			 nrfd, poll_socks ? "yes" : "no"));
		nr = select(maxfd + 1, &rfd, NULL, NULL, NULL);
		if (nr < 0) DPRINTF(("RX: select returned %zd errno %d\n", nr, errno));
		assert(nr >= 0);
		DPRINTF(("RX:\nRX: *** %d/%d fds are readable (descs: %s)\n",
			 nr, nrfd, vq_has_descs(vq) ? "yes" : "no"));

		pthread_mutex_lock(&sc->rx_mtx);

		if (FD_ISSET(sc->rx_wake_fd, &rfd)) {
			/* Eat the notification(s) */
			char dummy[128];
			ssize_t rd_dummy = read(sc->rx_wake_fd, &dummy, 128);
			assert(rd_dummy >= 1);
			/* Restart select now that we have some
			 * descriptors. It's possible that synchronous
			 * responses sent from the tx thread have
			 * eaten them all though, so check.
			 */
			DPRINTF(("RX: thread got %zd kicks (have descs: %s)\n",
				 rd_dummy, vq_has_descs(vq) ? "yes" : "no"));

// XXX need to check sockets in order to process closing, so cannot make this
// tempting looking optimisation.
//
//			if (nr == 1) {
//				 /* Must have been the kicker fd, in
//				  * which case there is no point
//				  * checking the socks.
//				  */
//				DPRINTF(("RX: Kicked w/ no other fds -- restarting select()\n"));
//				goto rx_done;
//			}

			/* We might have some descriptors, so it might be worth polling the socks again */
			poll_socks = true;
		}

		if (!vq_has_descs(vq)) {
			DPRINTF(("RX: No descs -- restarting select()\n"));
			poll_socks = false; /* Don't poll socks next time */
			goto rx_done;
		}

		while (did_some_work) {
			bool more_replies_pending = true; /* Assume there is */
			did_some_work = false;

			DPRINTF(("RX: Handling pending replies first\n"));
			pthread_mutex_lock(&sc->reply_mtx);
			while (vq_has_descs(vq)) {
				more_replies_pending = rx_do_one_reply(sc, vq);
				if (!more_replies_pending) break;
				did_some_work = true;
			}
			pthread_mutex_unlock(&sc->reply_mtx);

			if (more_replies_pending) {
				DPRINTF(("RX: No more descriptors for pending replies\n"));
				poll_socks = false; /* Still replies to send, so don't handle socks yet */
				vq_endchains(vq, 1);
				goto rx_done;
			}

			DPRINTF(("RX: Checking all socks\n"));

			for(i = 0; i < VTSOCK_MAXSOCKS; i++) {
				struct pci_vtsock_sock *s = lookup_sock_by_idx(sc, i);

				if (!s) continue;

				if (s->state == SOCK_CLOSING) { /* Closing comes through here */
					assert(s->local_shutdown == VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL ||
					       s->peer_shutdown == VIRTIO_VSOCK_FLAG_SHUTDOWN_ALL);
					DPRINTF(("RX: Closing sock %p fd %d local %"PRIx32" peer %"PRIx32"\n",
						 (void *)s, s->fd,
						 s->local_shutdown,
						 s->peer_shutdown));
					PPRINTF(("RX: SOCK closed (%d) "PRIaddr" <=> "PRIaddr"\n",
						 s->fd,
						 FMTADDR(s->local_addr), FMTADDR(s->peer_addr)));
					close(s->fd);
					s->fd = -1;
					s->state = SOCK_FREE;
					put_sock(s);
					continue;
				}

				if (s->state != SOCK_CONNECTED) {
					put_sock(s);
					continue;
				}

				assert(s->fd >= 0);

				if (FD_ISSET(s->fd, &rfd)) {
					ssize_t bytes;
					DPRINTF(("RX: event on sock %p fd %d\n",
						 (void *)s, s->fd));
					bytes = pci_vtsock_proc_rx(sc, vq, s);
					if (bytes == -1) {
						/* Consumed all descriptors, stop */
						DPRINTF(("RX: No more descriptors\n"));
						vq_endchains(vq, 1);
						put_sock(s);
						goto rx_done;
					} else if (bytes == 0) {
						FD_CLR(s->fd, &rfd);
					} else {
						did_some_work = true;
					}
					/* We sent either an OP_RW or an OP_SHUTDOWN in proc_rx */
					s->credit_update_required = false;
				} else if (s->credit_update_required) {
					if (send_credit_update(vq, s)) {
						s->credit_update_required = false;
					} else {
						/* Consumed all descriptors, stop */
						DPRINTF(("RX: No more descriptors\n"));
						vq_endchains(vq, 1);
						put_sock(s);
						goto rx_done;
					}
				}

				put_sock(s);
			}
		}

		DPRINTF(("RX: All work complete\n"));
		vq_endchains(vq, 0);
 rx_done:
		pthread_mutex_unlock(&sc->rx_mtx);

	}
}


static void pci_vtsock_notify_rx(void *vsc, struct vqueue_info *vq)
{
	struct pci_vtsock_softc *sc = vsc;

	assert(vq == &sc->vssc_vqs[VTSOCK_QUEUE_RX]);
	assert(sc->rx_wake_fd >= 0);

	kick_rx(sc, "notify");
}

static void pci_vtsock_notify_evt(void *vsc, struct vqueue_info *vq)
{
	struct pci_vtsock_softc *sc = vsc;

	DPRINTF(("vtsock: evt notify sc=%p vq=%zd(%p)\n",
		 (void *)sc, vq - &sc->vssc_vqs[VTSOCK_QUEUE_RX], (void *)vq));
}

static int listen_un(struct sockaddr_un *un)
{
	int fd, rc;

	rc = unlink(un->sun_path);
	if (rc < 0 && errno != ENOENT) {
		perror("Failed to unlink unix socket path");
		return -1;
	}

	fd = socket(AF_UNIX, SOCK_STREAM, 0);
	if (fd < 0) {
		perror("Failed to open unix socket");
		return -1;
	}

	rc = bind(fd, (struct sockaddr *)un, sizeof(*un));
	if (rc < 0) {
		perror("Failed to bind() unix socket");
		return -1;
	}

	rc = listen(fd, SOMAXCONN);
	if (rc < 0) {
		perror("Failed to listen() unix socket");
		return -1;
	}

	/* XXX Any chown/chmod needed? */

	rc = fcntl(fd, F_SETFL, O_NONBLOCK);
	if (rc < 0) {
		perror("O_NONBLOCK failed for unix socket\n");
		return -1;
	}

	return fd;
}

static int open_connect_socket(struct pci_vtsock_softc *sc)
{
	struct sockaddr_un un;
	int fd, rc;

	assert(sc->connect_fd = -1);

	bzero(&un, sizeof(un));

	un.sun_len = 0; /* Unused? */
	un.sun_family = AF_UNIX;
	rc = snprintf(un.sun_path, sizeof(un.sun_path),
		     "%s/"CONNECT_SOCKET_NAME, sc->path);
	if (rc < 0) {
		perror("Failed to format connect socket path");
		return 1;
	}
	DPRINTF(("Connect socket is %s\n", un.sun_path));

	fd = listen_un(&un);
	if (fd < 0) {
		fprintf(stderr, "failed to open connect socket\n");
		return 1;
	}

	sc->connect_fd = fd;
	DPRINTF(("Connect socket %s is fd %d\n", un.sun_path, fd));

	return 0;
}

static int open_one_forward_socket(struct pci_vtsock_softc *sc, uint32_t port)
{
	struct sockaddr_un un, sl;
	struct pci_vtsock_forward *fwd;
	int fd, rc;

	if (sc->nr_fwds == VTSOCK_MAXFWDS)  {
		fprintf(stderr, "Too many forwards\n");
		return 1;
	}

	fwd = &sc->fwds[sc->nr_fwds++];
	assert(fwd->listen_fd == -1);

	bzero(&un, sizeof(un));

	un.sun_len = 0; /* Unused? */
	un.sun_family = AF_UNIX;
	rc = snprintf(un.sun_path, sizeof(un.sun_path),
		     "%s/"PRIaddr, sc->path, sc->vssc_cfg.guest_cid, port);
	if (rc < 0) {
		perror("Failed to format forward socket path");
		return 1;
	}
	rc = snprintf(sl.sun_path, sizeof(sl.sun_path),
		     "%s/guest."PRIport, sc->path, port);
	if (rc < 0) {
		perror("Failed to format forward socket symlink path");
		return 1;
	}

	rc = unlink(sl.sun_path);
	if (rc < 0 && errno != ENOENT) {
		perror("Failed to unlink forward socket symlink path");
		return 1;
	}

	fd = listen_un(&un);
	if (fd < 0) {
		fprintf(stderr, "Failed to open forward socket\n");
		return 1;
	}

	rc = symlink(&un.sun_path[strlen(sc->path) + 1], sl.sun_path);
	if (rc < 0) {
		perror("Failed to create forward socket symlink\n");
		close(fd);
		return 1;
	}

	fwd->listen_fd = fd;
	fwd->port = port;

	DPRINTF(("forwarding port %"PRId32" to the guest\n", port));

	return 0;
}

static int open_forward_sockets(struct pci_vtsock_softc *sc,
				char *guest_forwards)
{
	char *s = guest_forwards, *e;
	int rc;

	if (!guest_forwards) return 0;

	while (*s != '\0') {
		unsigned long ul;

		rc = 1;
		errno = 0;
		ul = strtoul(s, &e, 0);

		if (errno) {
			fprintf(stderr, "failed to parse forward \"%s\": %s\n",
				s, strerror(errno));
			goto err;
		}
		if (ul >= UINT32_MAX) {
			fprintf(stderr, "invalid guest port forward %ld\n", ul);
			goto err;
		}

		rc = open_one_forward_socket(sc, (uint32_t)ul);
		if (rc) goto err;

		s = e;
		if (*s == ';') s++;
	}

	rc = 0;
err:
	free(guest_forwards);
	return rc;

}

static int pci_vtsock_cfgread(void *, int, int, uint32_t *);
static int pci_vtsock_cfgwrite(void *, int, int, uint32_t);

static char *
copy_up_to_comma(const char *from)
{
	char *comma = strchr(from, ',');
	char *tmp = NULL;
	if (comma == NULL) {
		tmp = strdup(from); /* rest of string */
	} else {
		size_t length = (size_t)(comma - from);
		tmp = strndup(from, length);
	}
	return tmp;
}

static int
pci_vtsock_init(struct pci_devinst *pi, char *opts)
{
	uint32_t guest_cid = VMADDR_CID_ANY;
	const char *path = NULL;
	char *guest_forwards = NULL;
	struct pci_vtsock_softc *sc;
	struct sockaddr_un un;
	int i, pipefds[2];

	if (opts == NULL) {
		printf("virtio-sock: configuration required\n");
		return (1);
	}

	while (1) {
		char *next;
		if (! opts)
			break;
		next = strchr(opts, ',');
		if (next)
			next[0] = '\0';
		if (strncmp(opts, "guest_cid=", 10) == 0) {
			int tmp = atoi(&opts[10]);
			if (tmp <= 0) {
				fprintf(stderr, "bad guest cid: %s\r\n", opts);
				return 1;
			}
			guest_cid = (uint32_t)tmp;
		} else if (strncmp(opts, "path=", 5) == 0) {
			path = copy_up_to_comma(opts + 5);
		} else if (strncmp(opts, "guest_forwards=", 15) == 0) {
			guest_forwards = copy_up_to_comma(opts + 15);
		} else {
			fprintf(stderr, "invalid option: %s\r\n", opts);
			return 1;
		}

		if (! next)
			break;
		opts = &next[1];
	}
	if (guest_cid == VMADDR_CID_ANY || path == NULL) {
		fprintf(stderr, "guest_cid and path options are both required.\n");
		return 1;
	}

	if (guest_cid <= VMADDR_CID_HOST) {
		fprintf(stderr, "invalid guest_cid %"PRIx32"\n", guest_cid);
		return 1;
	}

	/*
	 * We need to be able to construct socket paths of the form
	 * "%08x.%08x" cid,port.
	 */
	if (strlen(path) + sizeof("/00000000.00000000") > sizeof(un.sun_path)) {
		printf("virtio-sock: path too long\n");
		return (1);
	}

	/* XXX confirm path exists and is a directory */

	fprintf(stderr, "vsock init %d:%d = %s, guest_cid = %"PRIx32"\n\r",
		pi->pi_slot, pi->pi_func, path, guest_cid);

	sc = calloc(1, sizeof(struct pci_vtsock_softc));
	for (i = 0; i < VTSOCK_MAXSOCKS; i++) {
		struct pci_vtsock_sock *str = &sc->socks[i];
		int err = pthread_mutex_init(&str->mtx, NULL);
		assert(err == 0);
		str->state = SOCK_FREE;
		str->fd = -1;
	}

	sc->nr_fwds = 0;
	for (i = 0; i < VTSOCK_MAXFWDS; i++) {
		struct pci_vtsock_forward *fwd = &sc->fwds[i];
		fwd->listen_fd = -1;
	}

	pthread_mutex_init(&sc->vssc_mtx, NULL);
	pthread_mutex_init(&sc->tx_mtx, NULL);
	pthread_mutex_init(&sc->rx_mtx, NULL);
	pthread_mutex_init(&sc->reply_mtx, NULL);
	pthread_mutex_init(&sc->alloc_mtx, NULL);

	sc->path = strdup(path);

	/* init virtio softc and virtqueues */
	vi_softc_linkup(&sc->vssc_vs, &vtsock_vi_consts, sc, pi, sc->vssc_vqs);
	sc->vssc_vs.vs_mtx = &sc->vssc_mtx;

	sc->vssc_vqs[VTSOCK_QUEUE_RX].vq_qsize = VTSOCK_RINGSZ;
	sc->vssc_vqs[VTSOCK_QUEUE_RX].vq_notify = pci_vtsock_notify_rx;

	sc->vssc_vqs[VTSOCK_QUEUE_TX].vq_qsize = VTSOCK_RINGSZ;
	sc->vssc_vqs[VTSOCK_QUEUE_TX].vq_notify = pci_vtsock_notify_tx;

	/* Unused, make it small */
	sc->vssc_vqs[VTSOCK_QUEUE_EVT].vq_qsize = 4;
	sc->vssc_vqs[VTSOCK_QUEUE_EVT].vq_notify = pci_vtsock_notify_evt;

	/* setup virtio sock config space */
	sc->vssc_cfg.guest_cid = guest_cid;

	/*
	 * Should we move some of this into virtio.c?  Could
	 * have the device, class, and subdev_0 as fields in
	 * the virtio constants structure.
	 */
	pci_set_cfgdata16(pi, PCIR_DEVICE, VIRTIO_DEV_SOCK);
	pci_set_cfgdata16(pi, PCIR_VENDOR, VIRTIO_VENDOR);
	pci_set_cfgdata8(pi, PCIR_REVID, 0 /*LEGACY 1*/);
	pci_set_cfgdata8(pi, PCIR_CLASS, PCIC_NETWORK);
	pci_set_cfgdata16(pi, PCIR_SUBDEV_0, VIRTIO_TYPE_SOCK);
	pci_set_cfgdata16(pi, PCIR_SUBVEND_0, VIRTIO_VENDOR);

	if (vi_intr_init(&sc->vssc_vs, 1, fbsdrun_virtio_msix()))
		return (1);
	vi_set_io_bar(&sc->vssc_vs, 0);

	sc->connect_fd = -1;
	if (open_connect_socket(sc))
		return (1);

	if (open_forward_sockets(sc, guest_forwards))
		return (1);

	if (pipe(pipefds))
		return (1);
	sc->tx_wake_fd = pipefds[0];
	sc->tx_kick_fd = pipefds[1];

	if (pthread_create(&sc->tx_thread, NULL,
			   pci_vtsock_tx_thread, sc))
		return (1);

	if (pipe(pipefds))
		return (1);
	sc->rx_wake_fd = pipefds[0];
	sc->rx_kick_fd = pipefds[1];

	sc->reply_prod = 0;
	sc->reply_cons = 0;

	if (pthread_create(&sc->rx_thread, NULL,
			   pci_vtsock_rx_thread, sc))
		return (1);

	return (0);
}

static int
pci_vtsock_cfgwrite(UNUSED void *vsc, int offset, UNUSED int size,
	UNUSED uint32_t value)
{
	DPRINTF(("vtsock: write to readonly reg %d\n\r", offset));
	return (1);
}

static int
pci_vtsock_cfgread(void *vsc, int offset, int size, uint32_t *retval)
{
	struct pci_vtsock_softc *sc = vsc;
	void *ptr;

	DPRINTF(("vtsock: %d byte read pci reg %d\n\r", size, offset));

	/* our caller has already verified offset and size */
	ptr = (uint8_t *)&sc->vssc_cfg + offset;
	memcpy(retval, ptr, size);
	return (0);
}

static struct pci_devemu pci_de_vsock = {
	.pe_emu =	"virtio-sock",
	.pe_init =	pci_vtsock_init,
	.pe_barwrite =	vi_pci_write,
	.pe_barread =	vi_pci_read
};
PCI_EMUL_SET(pci_de_vsock);
