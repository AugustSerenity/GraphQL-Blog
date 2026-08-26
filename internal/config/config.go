package config

import "os"

type Config struct {
	DatabaseURL    string
	RepositoryType string
}

func Load() Config {
	return Config{
		DatabaseURL:    getEnv(os.Getenv("DB_URL"), ""),
		RepositoryType: getEnv(os.Getenv("REPOSITORY_TYPE"), ""),
	}
}

func getEnv(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
