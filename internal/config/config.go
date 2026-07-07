package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
	// App      AppConfig
}
type ServerConfig struct {
	Port string `envconfig:"PORT" required:"true"`
	Host string `envconfig:"HOST" required:"true"`
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type LogConfig struct {
	Level  string `envconfig:"LEVEL" required:"true"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func NewServConfig() (ServerConfig, error) {
	var config ServerConfig
	if err := envconfig.Process("SERVER", &config); err != nil {
		return ServerConfig{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}
func NewServConfigMust() ServerConfig {
	config, err := NewServConfig()
	if err != nil {
		err = fmt.Errorf("get Server config : %w", err)
		panic(err)
	}

	return config

}

func NewLogConfig() (LogConfig, error) {
	var config LogConfig
	if err := envconfig.Process("LOGGER", &config); err != nil {
		return LogConfig{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}
func NewLogConfigMust() LogConfig {
	config, err := NewLogConfig()
	if err != nil {
		err = fmt.Errorf("get Logger config : %w", err)
		panic(err)
	}

	return config

}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		panic("pupupu")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "user-123"),
			Password: getEnv("POSTGRES_PASSWORD", "password-123"),
			DBName:   getEnv("POSTGRES_DB", "subscription-db"),
		},
	}, nil
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.DBName)
}
