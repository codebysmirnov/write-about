package config

import "os"

// Config structure
type Config struct {
	DB      *DBConfig
	LogPath string
}

// DBConfig contain main fields for connect to DB
type DBConfig struct {
	Dialect  string
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Charset  string
}

// GetConfig initial configuration from env
func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  getEnv("DIALECT", "postgres"),
			Host:     getEnv("HOST", "127.0.0.1"),
			Port:     getEnv("PORT", "5432"),
			Username: getEnv("DB_USER", "smirnov"),
			Password: getEnv("DB_PASSWORD", "baikal"),
			Name:     getEnv("DB_NAME", "api"),
		},
		LogPath: getEnv("LOG_PATH", "app.log"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
