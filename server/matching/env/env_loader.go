package env

import (
	"log"
	"os"
)

func LoadEnvVar(envVar string) string {
	envVarStr := os.Getenv(envVar)

	if envVarStr == "" {
		log.Fatalf("Error: %s not provided, add env var", envVar)
	}

	return envVarStr
}
