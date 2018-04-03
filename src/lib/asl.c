#include <asl.h>
#include <pwd.h>
#include <fcntl.h>
#include <stdio.h>
#include <time.h>

#include <SystemConfiguration/SystemConfiguration.h>

#include <xhyve/asl.h>

static aslclient asl = NULL;
static aslmsg log_msg = NULL;

static unsigned char buf[4096];
static size_t buf_size = 0;

/* asl is deprecated in favor of os_log starting with macOS 10.12.  */
#pragma GCC diagnostic ignored "-Wdeprecated-declarations"

/* Initialize ASL logger and local buffer. */
void asl_init(void)
{
	asl = asl_open(NULL, NULL, 0);
	log_msg = asl_new(ASL_TYPE_MSG);
}


/* Send one character to the logger: wait for full lines before actually sending. */
void asl_put(uint8_t c)
{
	if ((c == '\n') || (c == 0)) {
		buf[buf_size] = 0;
		asl_log(asl, log_msg, ASL_LEVEL_NOTICE, "%s", buf);
		buf_size = 0;
	} else {
		if (buf_size + 2 >= sizeof buf) {
			/* Running out of space, flush.  */
			asl_put('\n');
		}
		buf[buf_size] = c;
		++buf_size;
	}
}
