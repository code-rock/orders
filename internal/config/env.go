package envconfig

import (
	"os"
	"strconv"
	"strings"
)

type SDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type SGitHubConfig struct {
	Username string
	APIKey   string
}

type Config struct {
	GitHub SGitHubConfig
	BD     SDBConfig
	// DebugMode bool
	// UserRoles []string
	// MaxUsers  int
}

func New() *Config {
	return &Config{
		GitHub: SGitHubConfig{
			Username: getEnv("GITHUB_USERNAME", ""),
			APIKey:   getEnv("GITHUB_API_KEY", ""),
		},
		BD: SDBConfig{
			Host:     getEnv("POSTGRE_SQL_HOST", ""),
			Port:     getEnvAsInt("POSTGRE_SQL_PORT", 5432),
			User:     getEnv("POSTGRE_SQL_USER", ""),
			Password: getEnv("POSTGRE_SQL_PASSWORD", ""),
			DBName:   getEnv("POSTGRE_SQL_NAME", ""),
		},
		// DebugMode: getEnvAsBool("DEBUG_MODE", true),
		// UserRoles: getEnvAsSlice("USER_ROLES", []string{"admin"}, ","),
		// MaxUsers:  getEnvAsInt("MAX_USERS", 1),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
