package hyperkit

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -framework CoreFoundation -framework Foundation -framework ServiceManagement -framework SystemConfiguration -framework CoreServices -framework IOKit

#include "log_darwin.h"
*/
import "C"

// InsinuateSystemLogger modifies a logrus Logger to sent its logs to
// the system's logger (e.g, Apple's ASL).
func InsinuateSystemLogger(log *logrus.Logger) *logrus.Logger {
	log.SetLevel(logrus.DebugLevel)
	log.AddHook(NewLogrusASLHook())
	// In addition to the default destiniation, our logs are sent
	// to ASL via the previous hook.  But since our stderr is also
	// redirected to ASL, each entry would appear twice: don't
	// send logs to stderr.
	log.Out = ioutil.Discard
	return log
}

// LogrusASLHook defines a hook for Logrus that redirects logs
// to ASL API (to be displayed in Console application)
type LogrusASLHook struct {
}

// NewLogrusASLHook returns a new LogrusASLHook
func NewLogrusASLHook() *LogrusASLHook {
	hook := new(LogrusASLHook)
	C.apple_asl_logger_init(C.CString("Docker"), C.CString(filepath.Base(os.Args[0])))
	return hook
}

// Levels returns the available ASL log levels
func (t *LogrusASLHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func aslLevel(l logrus.Level) C.int {
	switch l {
	case logrus.PanicLevel:
		return C.ASL_LEVEL_ALERT
	case logrus.FatalLevel:
		return C.ASL_LEVEL_CRIT
	case logrus.ErrorLevel:
		return C.ASL_LEVEL_ERR
	case logrus.WarnLevel:
		return C.ASL_LEVEL_WARNING
	case logrus.InfoLevel:
		return C.ASL_LEVEL_NOTICE
	case logrus.DebugLevel:
		return C.ASL_LEVEL_DEBUG
	}
	return C.ASL_LEVEL_DEBUG
}

// Fire sends a log entry to ASL
func (t *LogrusASLHook) Fire(entry *logrus.Entry) error {
	C.apple_asl_logger_log(aslLevel(entry.Level), C.CString(entry.Message))
	return nil
}
