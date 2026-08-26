package config

import "os"

type Config struct {
	DatabaseURL string
}

func Load() Config {
	return Config{
		DatabaseURL: getEnv(os.Getenv("DB_URL"), "host=localhost port=5432 user=postgres dbname=graphBlog password=postgres"),
	}
}

func getEnv(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
