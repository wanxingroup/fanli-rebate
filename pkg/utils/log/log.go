package log

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func init() {
	logger = logrus.NewEntry(logrus.New())
}

func GetLogger() *logrus.Entry {
	return logger
}
