package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	GRPC     GRPCConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	Auth     AuthConfig
}

func LoadConfig(envPath string) Config {
	if envPath != "" {
		_ = godotenv.Load(envPath)
	}

	return Config{
		App: AppConfig{
			ServiceName: getEnv("SERVICE_NAME", "customer-service"),
			Env:         getEnv("APP_ENV", "local"),
		},
		HTTP: HTTPConfig{
			Host:     getEnv("HTTP_HOST", "0.0.0.0"),
			Port:     getEnv("HTTP_PORT", "8443"),
			CertFile: getEnv("HTTP_CERT_FILE", "certs/server.crt"),
			KeyFile:  getEnv("HTTP_KEY_FILE", "certs/server.key"),
		},
		GRPC: GRPCConfig{
			Host:                getEnv("GRPC_HOST", "0.0.0.0"),
			Port:                getEnv("GRPC_PORT", "9091"),
			SupplierServiceAddr: getEnv("SUPPLIER_SERVICE_GRPC_ADDR", "localhost:9091"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "softplace"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConns: int32(getEnvInt("DB_MAX_CONNS", 10)),
			MinConns: int32(getEnvInt("DB_MIN_CONNS", 2)),
		},
		Logger: LoggerConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "text"),
			Outputs:    splitEnv("LOG_OUTPUTS", "stdout"),
			FilePath:   getEnv("LOG_FILE_PATH", ""),
			FileDir:    getEnv("LOG_FILE_DIR", "logs"),
			FilePrefix: getEnv("LOG_FILE_PREFIX", getEnv("SERVICE_NAME", "app")),
			AddSource:  getEnvBool("LOG_ADD_SOURCE", true),
		},
		Auth: AuthConfig{
			JWTSecret:        getEnv("JWT_SECRET", "change-me-secret"),
			JWTLifetimeHours: getEnvInt("JWT_LIFETIME_HOURS", 24),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func splitEnv(key string, defaultValue string) []string {
	value := getEnv(key, defaultValue)

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}
