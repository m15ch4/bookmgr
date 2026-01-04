package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabaseHost     string
	DatabasePort     int
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	SkipBootstrap    bool
	ServerPort       int
}

func Load() (*Config, error) {
	config := &Config{
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabaseName:     getEnv("DATABASE_NAME", "bookdb"),
		DatabaseUser:     getEnv("DATABASE_USER", "root"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", ""),
		ServerPort:       getEnvAsInt("SERVER_PORT", 8080),
	}

	port := getEnvAsInt("DATABASE_PORT_NUMBER", 3306)
	config.DatabasePort = port

	skipBootstrap := getEnv("SKIP_BOOTSTRAP", "false")
	config.SkipBootstrap = skipBootstrap == "true" || skipBootstrap == "1"

	return config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
