package env

import (
	"fmt"
	"os"
)

func Get(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		var valInt int
		fmt.Scanf(val, "%d", valInt)
		return valInt
	}

	return defaultValue
}
