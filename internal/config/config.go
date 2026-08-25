package config

import "os"

type Config struct {
	DBPath string
}

func Load() Config {
	return Config{
		DBPath: getEnv(os.Getenv("DB_PATH"), "host=localhost port=5432 user=postgres dbname=graphBlog password=postgres"),
	}
}

func getEnv(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
