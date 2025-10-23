package config

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Config holds all configuration settings
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// LoggerConfig holds logger-related configuration
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	OutputPath string `mapstructure:"output_path"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// NewConfig creates a new Config instance from configuration sources
func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("conf.d")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// ProviderSet is a Wire provider set for configuration
var ProviderSet = wire.NewSet(NewConfig)
