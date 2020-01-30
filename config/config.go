package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

// GetVar returns required environment variable.
// If required variable not found will return empty value
func GetVar(v string) (string, error) {
	result, exists := os.LookupEnv(v)
	if !exists {
		return "", fmt.Errorf("Env variable not found %s", v)
	}

	return result, nil
}
