package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config is a struct that holds the configuration for the application.
type Config struct {
	// Server related
	ListenPort     string
	CookieSecret   string
	RequestsPerMin int
	TrustedOrigins []string

	// Upload related
	FileUploadPath string
}

// New returns a new Config struct.
func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		CookieSecret:   getEnvString("COOKIE_SECRET", "mysecret"),
		RequestsPerMin: getEnvInt("REQUESTS_PER_MIN", 30),
		ListenPort:     getEnvString("LISTEN_PORT", ":8080"),
		TrustedOrigins: strings.Split(getEnvString("CORS_TRUSTED_ORIGINS", "http://localhost:3000"), ","),
		FileUploadPath: getEnvString("FILE_UPLOAD_PATH", "uploads"),
	}
}

// getEnvString gets the environment variable or returns the default value.
func getEnvString(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}

// getEnvInt gets the environment variable or returns the default value.
func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	num, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return num
}
