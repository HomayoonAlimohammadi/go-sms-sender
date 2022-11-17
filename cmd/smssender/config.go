package main

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SmsSenderConfig struct {
	ApiKey          string
	Host            string
	Port            string
	GracefulTimeout time.Duration
}

type PostgresConfig struct {
	Host     string
	Port     string
	Password string
	DbName   string
	SslMode  string
}

type Config struct {
	SmsSender SmsSenderConfig
	DB        PostgresConfig
}

func loadConfig(cmd *cobra.Command, args []string) (*Config, error) {

	// load config from file
	configFile, err := cmd.Flags().GetString("config-file")
	if err == nil && configFile != "" {
		viper.SetConfigFile(configFile)

		if err := viper.ReadInConfig(); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Println(config)
	return &config, nil
}
