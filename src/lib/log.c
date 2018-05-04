#include <asl.h>
#include <fcntl.h>
#include <pwd.h>
#include <stdio.h>
#include <sysexits.h>
#include <time.h>

#include <SystemConfiguration/SystemConfiguration.h>

#include <xhyve/log.h>

/* Some functions below invoke functions with non-literal format
 * strings, and it's ok: their (top-level) fmt string was checked. */
#pragma GCC diagnostic ignored "-Wformat-nonliteral"

static aslclient log_client = NULL;
static aslmsg log_msg = NULL;

static char buf[4096];
/* Index of the _next_ character to insert in the buffer. */
static size_t buf_idx = 0;

/* asl is deprecated in favor of os_log starting with macOS 10.12.  */
#pragma GCC diagnostic ignored "-Wdeprecated-declarations"

void log_init(void)
{
	log_client = asl_open(NULL, NULL, 0);
	log_msg = asl_new(ASL_TYPE_MSG);
}


/* Send the content of the buffer to the logger. */
static void log_flush(void)
{
	buf[buf_idx] = 0;
	asl_log(log_client, log_msg, ASL_LEVEL_NOTICE, "%s", buf);
	buf_idx = 0;
}


/* Send one character to the logger: wait for full lines before actually sending. */
static void log_put(char c)
{
	if ((c == '\n') || (c == 0)) {
		log_flush();
	} else {
		if (buf_idx + 2 >= sizeof(buf)) {
			log_flush();
		}
		buf[buf_idx] = c;
		++buf_idx;
	}
}


/* Send a string to the logger: wait for full lines before actually sending. */
static void log_puts(const char *s)
{
	for (; *s; ++s) {
		log_put(*s);
	}
}


/* Send a string to the logger: wait for full lines before actually sending. */
static int log_vprintf(const char *fmt, va_list args)
{
	static char buf[4096];

	int res = vsnprintf(buf, sizeof(buf), fmt, args);

	log_puts(buf);
	return (res);
}


/* Destination for logs.  */
static enum {
	LOG_DST_LOG,
	LOG_DST_STDERR
}
log_dst = LOG_DST_STDERR;

void log_set_destination(const char *dst)
{
	if (!dst) {
		log_fprintf(stderr,
		    "log_set_destination: invalid NULL argument");
		exit(EX_USAGE);
	} else if (!strcmp(dst, "log")) {
		log_dst = LOG_DST_LOG;
	} else if (!strcmp(dst, "stderr")) {
		log_dst = LOG_DST_STDERR;
	} else {
		log_fprintf(stderr, "log_set_destination: invalid argument: %s",
		    dst);
		exit(EX_USAGE);
	}
}


int log_fprintf(FILE *f, const char *fmt, ...)
{
	va_list args;
	int res = 0;

	va_start(args, fmt);
	if ((log_dst == LOG_DST_LOG) &&
	    ((f == stdout) || (f == stderr))) {
		res = log_vprintf(fmt, args);
	} else {
		res = vfprintf(f, fmt, args);
	}
	va_end(args);
	return (res);
}


/* A buffer for the console.  */
static char console_buf[4096];
/* Index of the _next_ character to insert in console_buf. */
static size_t console_idx = 0;


/* Send the content of the buffer to the logger. */
static void log_console_flush(void)
{
	console_buf[console_idx] = 0;
	log_fprintf(stderr, "%s\n", console_buf);
	console_idx = 0;
}


/* Send one character to the logger: wait for full lines before actually sending. */
void log_console_put(char c)
{
	if ((c == '\n') || (c == 0)) {
		log_console_flush();
	} else {
		if (console_idx + 2 >= sizeof(console_buf)) {
			log_console_flush();
		}
		console_buf[console_idx] = c;
		++console_idx;
	}
}
