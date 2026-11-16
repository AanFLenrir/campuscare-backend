package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	JWTExpireHour int
	AppEnv        string
}

// loadDotEnv reads .env file if present and sets environment variables if not already set.
func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			// Strip optional quotes
			val = strings.Trim(val, `"'`)
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}
}

func Load() *Config {
	// Automatically try loading .env from root or parent directory
	loadDotEnv(".env")
	loadDotEnv("../.env")

	expire, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOUR", "24"))
	if err != nil {
		expire = 24
	}

	// Prefer 127.0.0.1 as default host to avoid Windows IPv4/IPv6 resolution ambiguity
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	if dbHost == "localhost" {
		dbHost = "127.0.0.1"
	}

	return &Config{
		AppPort:       getEnv("APP_PORT", "3000"),
		DBHost:        dbHost,
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "password"),
		DBName:        getEnv("DB_NAME", "petriva"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		JWTSecret:     getEnv("JWT_SECRET", "change_me_secret_petriva_jwt_key_2026"),
		JWTExpireHour: expire,
		AppEnv:        getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (c *Config) DSN() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=" + c.DBSSLMode
}

func (c *Config) JWTExpireDuration() time.Duration {
	return time.Hour * time.Duration(c.JWTExpireHour)
}
