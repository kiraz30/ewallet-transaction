package helpers

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

func SetupLogger() {
	//use farmework logrus
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	log.Info("Logger initialized Using logrus")
	Logger = log
}
