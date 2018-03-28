#include <fcntl.h>
#include <stdio.h>
#include <time.h>
#include <SystemConfiguration/SystemConfiguration.h>

#include "log_darwin.h"

// logs

static aslclient asl = NULL;
static aslmsg log_msg = NULL;

// asl is deprecated in favor of os_log starting with macOS 10.12.
#pragma GCC diagnostic ignored "-Wdeprecated-declarations"

void apple_asl_logger_init(const char* sender, const char* facility) {
    free(asl);
    // I believe that these guys, sender and facility, are useless.  I
    // never managed to have them show in the actual logs.
    asl = asl_open(sender, facility, 0);
    log_msg = asl_new(ASL_TYPE_MSG);
}

static int is_initialized() {
  if (asl) {
    return true;
  } else {
      aslclient tmp_asl = asl_open("Docker", "com.docker.docker", ASL_OPT_STDERR);
      log_msg = asl_new(ASL_TYPE_MSG);
      asl_log(tmp_asl, log_msg, ASL_LEVEL_ERR, "asl_logger_init must be called before asl_logger_log");
      free(tmp_asl);
      return false;
  }
}

void apple_asl_logger_log(int level, const char *message) {
  if (!is_initialized())
    return;

  // The max length for log entries is 1024 bytes.  Beyond, they are
  // truncated.  In that case, split into several log entries.
  const size_t len = strlen(message);
  if (len < 1024)
    asl_log(asl, log_msg, level, "%s", message);
  else {
    enum { step = 1000 };
    for (int pos = 0; pos < len; pos += step) {
      asl_log(asl, log_msg, level,
              "%s%.*s%s",
              pos ? "[...] " : "",
              step, message + pos,
              pos + step < len ? " [...]" : "");
    }
  }
}
