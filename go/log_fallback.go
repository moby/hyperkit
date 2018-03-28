// +build !darwin

package hyperkit

import (
	"github.com/sirupsen/logrus"
)

func InsinuateSystemLogger(log *logrus.Logger) *logrus.Logger {
	return log
}
