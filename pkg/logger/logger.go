package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// todo change to ZAP
var Log *logrus.Logger

func init() {
	Log = logrus.New()

	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	Log.SetOutput(os.Stdout)

	Log.SetLevel(logrus.InfoLevel)
}

func SetLogLevel(level logrus.Level) {
	Log.SetLevel(level)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}
