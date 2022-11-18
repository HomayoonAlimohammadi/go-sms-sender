package main

import (
	"log"
	"time"

	"github.com/homayoonalimohammadi/go-sms-sender/smssender/internal/database"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SmsSenderConfig struct {
	ApiKey          string
	Port            string
	Host            string
	GracefulTimeout time.Duration
}

type Config struct {
	SmsSender  SmsSenderConfig
	PostgresDB database.PostgresConfig
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

	log.Printf("Config: %+v \n", config)
	return &config, nil
}
