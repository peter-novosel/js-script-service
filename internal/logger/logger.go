package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Init() *logrus.Logger {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	env := os.Getenv("ENV")
	if env == "production" {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logrus.DebugLevel)
	}

	return log
}
