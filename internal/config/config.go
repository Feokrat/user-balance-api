package config

import (
	"log"
	"strings"

	_ "github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type (
	Config struct {
		HTTP       HTTPConfig
		Postgresql PGConfig
	}

	HTTPConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	PGConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"ssl"`
	}
)

func Init(path string, logger *log.Logger) (*Config, error) {
	if err := parseConfigFile(path); err != nil {
		logger.Printf("failed to parse path to config file: %s", err)
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg, logger); err != nil {
		logger.Printf("failed to unmarshal config: %s", err)
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config, logger *log.Logger) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		logger.Printf("failed to unmarshal http key in config: %s", err)
		return err
	}

	if err := viper.UnmarshalKey("postgres", &cfg.Postgresql); err != nil {
		logger.Printf("failed to unmarshal postgres key in config: %s", err)
		return err
	}

	return nil
}

func parseConfigFile(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0])
	viper.SetConfigName(path[1])

	return viper.ReadInConfig()
}
