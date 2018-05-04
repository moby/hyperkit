#pragma once

/* Initialize logger. */
void log_init(void);

/* Send one character to the console logger. */
void log_console_put(char _c);

/* Specify where logs should be sent: `stderr` or `log`.  Dies on error.  */
void log_set_destination(const char* dst);

/* If logging is enabled, intercept outputs to stdout/stderr to the
 * logger, otherwise forward to fprintf. */
__attribute__ ((format (printf, 2, 3)))
int log_fprintf(FILE* f, const char *fmt, ...);

/* Intercept all the calls to fprintf to honor log_set_destination.  */
#define fprintf log_fprintf
#define printf(...) log_fprintf(stdout, __VA_ARGS__)
