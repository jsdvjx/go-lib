package utils

import "os"

func Env(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func RunEnv() string {
	return Env("ENV", "dev")
}
