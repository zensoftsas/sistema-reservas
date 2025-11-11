package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration settings
type Config struct {
	ServerPort        string
	DatabaseURL       string // Changed from DatabasePath
	JWTSecret         string
	JWTExpirationHrs  int
	SendGridAPIKey    string
	SendGridFromEmail string
	SendGridFromName  string
	AllowedOrigins    string
}

// LoadConfig loads configuration from environment variables and .env file
// Returns a Config struct with all necessary application settings
func LoadConfig() *Config {
	// Try to load .env file (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Read environment variables with default values
	serverPort := getEnv("SERVER_PORT", "8080")
	databaseURL := getEnv("DATABASE_URL", "./clinica.db")
	jwtSecret := getEnv("JWT_SECRET", "")
	jwtExpHours := getEnvAsInt("JWT_EXPIRATION_HOURS", 24)

	// SendGrid configuration (optional)
	sendGridAPIKey := getEnv("SENDGRID_API_KEY", "")
	sendGridFromEmail := getEnv("SENDGRID_FROM_EMAIL", "noreply@clinica.com")
	sendGridFromName := getEnv("SENDGRID_FROM_NAME", "Clinica Internacional")

	// CORS configuration
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:8080,http://localhost:8081")

	// Validate required configuration
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required in environment variables")
	}

	// Return configuration
	return &Config{
		ServerPort:        serverPort,
		DatabaseURL:       databaseURL,
		JWTSecret:         jwtSecret,
		JWTExpirationHrs:  jwtExpHours,
		SendGridAPIKey:    sendGridAPIKey,
		SendGridFromEmail: sendGridFromEmail,
		SendGridFromName:  sendGridFromName,
		AllowedOrigins:    allowedOrigins,
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt retrieves an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
