package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewViper(log *logrus.Logger) *viper.Viper {
	config := viper.New()

	config.AutomaticEnv()
	config.SetConfigName(".env")
	config.SetConfigType("dotenv")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")

	err := config.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.WithError(err).Panic("fatal error config file")
		}
	}

	return config
}