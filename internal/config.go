package main

import (
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	mainEnvFile = "internal/.env"
)

type Config struct {
	LogLevel          string        `mapstructure:"LOG_LEVEL"`
	HTTPPort          string        `mapstructure:"HTTP_PORT"`
	EthereumRPCAPIURL string        `mapstructure:"ETHEREUM_RPC_API_URL"`
	PoolingTime       time.Duration `mapstructure:"POOLING_TIME"`
}

func (c *Config) IsValid() error {
	return nil
}

func (c *Config) LogFields() map[string]interface{} {
	return map[string]interface{}{
		"LogLevel":          c.LogLevel,
		"HTTPPort":          c.HTTPPort,
		"EthereumRPCAPIURL": c.EthereumRPCAPIURL,
	}
}

func loadConfig(filenames ...string) (*Config, error) {
	var cfg = &Config{}

	filenames = append(filenames, mainEnvFile)

	viper.SetConfigType("env")
	viper.AutomaticEnv()

	for _, filename := range filenames {
		if _, err := os.Stat(filename); err != nil {
			log.Printf("Skipping load env file %s: %v", filename, err)
			continue
		}

		viper.SetConfigFile(filename)

		if err := viper.ReadInConfig(); err != nil {
			return nil, errors.Wrapf(err, "error to read config, path: %s", mainEnvFile)
		}

		if err := viper.MergeInConfig(); err != nil {
			return nil, errors.Wrapf(err, "error to merge config, filename: %s", filename)
		}

		if err := viper.Unmarshal(&cfg); err != nil {
			return nil, errors.Wrapf(err, "error to unmarshal config, filename: %s", filename)
		}
	}

	return cfg, nil
}
