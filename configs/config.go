package configs

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Log      LogConfig
	TLS      TLSConfig
}

// AppConfig holds all application-specific configuration
type AppConfig struct {
	Env  string
	Port string
	Host string
	Name string
	URL  string
}

// DatabaseConfig holds all database-specific configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	Charset         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// JWTConfig holds all JWT-specific configuration
type JWTConfig struct {
	Secret               string
	AccessTokenExpiry    time.Duration
	RefreshTokenExpiry   time.Duration
}

// CORSConfig holds all CORS-specific configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// LogConfig holds all logging-specific configuration
type LogConfig struct {
	Level  string
	Format string
}

// TLSConfig holds all TLS-specific configuration
type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

// LoadConfig loads configuration from environment variables and/or config files
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, using environment variables")
	}

	// Set defaults
	config := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
			Host: getEnv("APP_HOST", "0.0.0.0"),
			Name: getEnv("APP_NAME", "zione-backend"),
			URL:  getEnv("APP_URL", "http://localhost:8080"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "3306"),
			Name:            getEnv("DB_NAME", "zione_db"),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", ""),
			Charset:         getEnv("DB_CHARSET", "utf8mb4"),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", time.Hour),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "default-jwt-secret-change-in-production"),
			AccessTokenExpiry:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRY", 7*24*time.Hour), // 7 days
		},
		CORS: CORSConfig{
			AllowedOrigins: getStringSliceEnv("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
			AllowedMethods: getStringSliceEnv("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getStringSliceEnv("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		TLS: TLSConfig{
			Enabled:  getBoolEnv("TLS_ENABLED", false),
			CertFile: getEnv("TLS_CERT_FILE", "./certs/server.crt"),
			KeyFile:  getEnv("TLS_KEY_FILE", "./certs/server.key"),
		},
	}

	return config, nil
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	v := viper.New()
	v.Set(key, value)
	return v.GetInt(key)
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	v := viper.New()
	v.Set(key, value)
	return v.GetBool(key)
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	v := viper.New()
	v.Set(key, value)
	return v.GetDuration(key)
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return strings.Split(value, ",")
} 