package config

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Host string
	Port int
	Mode string // "debug" or "release"
}

type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	Path            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type JWTConfig struct {
	Secret       string
	ExpirationHs int // hours
}

var dotenvBaseDir string

func Load() (*Config, error) {
	dotenvBaseDir = ""
	if err := loadDotEnv(); err != nil {
		return nil, err
	}
	defaultDBPath := "../data/cdk_recharge.db"
	if dotenvBaseDir != "" {
		defaultDBPath = "data/cdk_recharge.db"
	}

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "sqlite"),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "cdk_recharge"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			Path:            resolvePath(getEnv("DB_PATH", defaultDBPath)),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME", 5)) * time.Minute,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			PoolSize: getEnvInt("REDIS_POOL_SIZE", 10),
		},
		JWT: JWTConfig{
			Secret:       getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpirationHs: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		},
	}

	return cfg, nil
}

func (c *DatabaseConfig) DSN() string {
	connectionURL := &url.URL{
		Scheme: "postgres",
		Host:   net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:   "/" + strings.TrimPrefix(c.DBName, "/"),
	}
	connectionURL.User = url.UserPassword(c.User, c.Password)
	query := connectionURL.Query()
	query.Set("sslmode", c.SSLMode)
	connectionURL.RawQuery = query.Encode()
	return connectionURL.String()
}

func loadDotEnv() error {
	paths, err := envFileCandidates()
	if err != nil {
		return err
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		if err := loadDotEnvFile(path); err != nil {
			return fmt.Errorf("load %s: %w", path, err)
		}
		dotenvBaseDir = filepath.Dir(path)
		log.Printf("已加载环境文件: %s", path)
		return nil
	}
	return nil
}

func resolvePath(path string) string {
	if filepath.IsAbs(path) || dotenvBaseDir == "" {
		return path
	}
	return filepath.Join(dotenvBaseDir, path)
}

func envFileCandidates() ([]string, error) {
	if explicit := strings.TrimSpace(os.Getenv("ENV_FILE")); explicit != "" {
		return []string{explicit}, nil
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var paths []string
	for dir := workingDir; ; dir = filepath.Dir(dir) {
		paths = append(paths, filepath.Join(dir, ".env"))
		if filepath.Base(dir) != "backend" {
			paths = append(paths, filepath.Join(dir, "backend", ".env"))
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	return paths, nil
}

func loadDotEnvFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for lineNumber, rawLine := range strings.Split(string(content), "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) == "" {
			return fmt.Errorf("invalid assignment on line %d", lineNumber+1)
		}
		key = strings.TrimSpace(key)
		if !isEnvKey(key) {
			return fmt.Errorf("invalid variable name %q on line %d", key, lineNumber+1)
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		value = strings.TrimSpace(value)
		if len(value) >= 2 && ((value[0] == '\'' && value[len(value)-1] == '\'') || (value[0] == '"' && value[len(value)-1] == '"')) {
			value = value[1 : len(value)-1]
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set %s: %w", key, err)
		}
	}
	return nil
}

func isEnvKey(key string) bool {
	if key == "" {
		return false
	}
	for index, r := range key {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' || (index > 0 && r >= '0' && r <= '9') {
			continue
		}
		return false
	}
	return true
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if result, err := strconv.Atoi(value); err == nil {
			return result
		}
	}
	return defaultValue
}
