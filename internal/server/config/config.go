package config

import (
	"os"
	"strconv"
	"time"
)

// Config contém todas as configurações da aplicação
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	API      APIConfig
}

// ServerConfig contém configurações do servidor HTTP
type ServerConfig struct {
	Port string
}

// DatabaseConfig contém configurações do banco de dados
type DatabaseConfig struct {
	DSN            string
	MaxConnections int
	Timeout        time.Duration
}

// APIConfig contém configurações da API externa
type APIConfig struct {
	BaseURL string
	Timeout time.Duration
}

// LoadConfig carrega configurações das variáveis de ambiente com valores padrão
func LoadConfig() (*Config, error) {
	config := &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		API:      loadAPIConfig(),
	}

	// Validação básica
	if config.API.BaseURL == "" {
		config.API.BaseURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	}

	return config, nil
}

// loadServerConfig carrega configurações do servidor
func loadServerConfig() ServerConfig {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return ServerConfig{
		Port: port,
	}
}

// loadDatabaseConfig carrega configurações do banco de dados
func loadDatabaseConfig() DatabaseConfig {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = ".cotacoes.db"
	}

	maxConnectionsStr := os.Getenv("DB_MAX_CONNECTIONS")
	maxConnections := 10 // Valor padrão
	if maxConnectionsStr != "" {
		if parsed, err := strconv.Atoi(maxConnectionsStr); err == nil {
			maxConnections = parsed
		}
	}

	timeoutStr := os.Getenv("DB_TIMEOUT")
	timeout := 10 * time.Millisecond // Valor padrão
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

// loadAPIConfig carrega configurações da API externa
func loadAPIConfig() APIConfig {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	}

	timeoutStr := os.Getenv("API_TIMEOUT")
	timeout := 200 * time.Millisecond // Valor padrão
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
