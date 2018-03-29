#include <asl.h>
#include <pwd.h>
#include <fcntl.h>
#include <stdio.h>
#include <time.h>

#include <SystemConfiguration/SystemConfiguration.h>

#include <xhyve/asl.h>

static aslclient asl = NULL;
static aslmsg log_msg = NULL;

static unsigned char *buf = NULL;
static size_t buf_size = 0;
static size_t buf_capacity = 0;

/* asl is deprecated in favor of os_log starting with macOS 10.12.  */
#pragma GCC diagnostic ignored "-Wdeprecated-declarations"

/* Grow buf/buf_capacity. */
static void buf_grow(void)
{
	buf_capacity = buf_capacity ? 2 * buf_capacity : 1024;
	buf = realloc(buf, buf_capacity);
	if (!buf) {
		perror("buf_grow");
		exit(1);
	}
}


/* Initialize ASL logger and local buffer. */
void asl_init(void)
{
	asl = asl_open(NULL, NULL, 0);
	log_msg = asl_new(ASL_TYPE_MSG);
	buf_grow();
}


/* Send one character to the logger: wait for full lines before actually sending. */
void asl_put(uint8_t c)
{
	if (buf_size + 1 >= buf_capacity) {
		buf_grow();
	}
	if (c == '\n') {
		buf[buf_size] = 0;
		asl_log(asl, log_msg, ASL_LEVEL_NOTICE, "%s", buf);
		buf_size = 0;
	} else {
		buf[buf_size] = c;
		++buf_size;
	}
}
