package config

import "github.com/sirupsen/logrus"

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}