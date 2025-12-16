package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	API      APIConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	DSN            string
	MaxConnections int
	Timeout        time.Duration
}

type APIConfig struct {
	BaseURL string
	Timeout time.Duration
}

func LoadConfig() (*Config, error) {
	config := &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		API:      loadAPIConfig(),
	}

	if config.API.BaseURL == "" {
		config.API.BaseURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	}

	return config, nil
}

func loadServerConfig() ServerConfig {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return ServerConfig{
		Port: port,
	}
}

func loadDatabaseConfig() DatabaseConfig {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = ".cotacoes.db"
	}

	maxConnectionsStr := os.Getenv("DB_MAX_CONNECTIONS")
	maxConnections := 10
	if maxConnectionsStr != "" {
		if parsed, err := strconv.Atoi(maxConnectionsStr); err == nil {
			maxConnections = parsed
		}
	}

	timeoutStr := os.Getenv("DB_TIMEOUT")
	timeout := 10 * time.Millisecond
	if timeoutStr != "" {
		if parsed, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = parsed
		}
	}

	return DatabaseConfig{
		DSN:            dsn,
		MaxConnections: maxConnections,
		Timeout:        timeout,
	}
}

func loadAPIConfig() APIConfig {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	}

	timeoutStr := os.Getenv("API_TIMEOUT")
	timeout := 200 * time.Millisecond
	if timeoutStr != "" {
		if parsed, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = parsed
		}
	}

	return APIConfig{
		BaseURL: baseURL,
		Timeout: timeout,
	}
}
