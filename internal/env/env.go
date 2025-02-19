package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if valueInt, err := strconv.Atoi(value); err == nil {
			return valueInt
		}
		return fallback
	}
	return fallback
}
