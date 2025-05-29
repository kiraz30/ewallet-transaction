package helpers

import (
	"log"

	"github.com/joho/godotenv"
)

var Env = map[string]string{}

func SetupConfig() {
	var err error
	Env, err = godotenv.Read(".env")
	if err != nil {
		log.Fatalf("Failed to read env: %v", err)
	}
}

func GetEnv(key, val string) string {
	//val is used as a default value
	result := Env[key]
	if result == "" {
		result = val
	}
	return result
}
