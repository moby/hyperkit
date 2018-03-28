package hyperkit

import (
	"github.com/sirupsen/logrus"
)

// log receives stdout/stderr of the hyperkit process itself, if set.
// It defaults to the go standard logger.
var log = logrus.StandardLogger()

// SetLogger sets the logger to use.
func SetLogger(l *logrus.Logger) {
	log = l
}
