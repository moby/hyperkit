#pragma once

#include <asl.h>
#include <pwd.h>

// aslInit must called before calling aslLog.
extern void apple_asl_logger_init(const char* sender, const char* facility);
extern void apple_asl_logger_log(int level, const char *message);

