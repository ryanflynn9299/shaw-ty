package utils

import (
	"github.com/subosito/gotenv"
	"time"
)

// InitializeEnv loads environment variables and sets up JWT configuration
func InitializeEnv() {
	// Load .env file if it exists
	err := gotenv.Load()
	if err != nil {
		panic(err)
	}
}

func ValidateMachineId(id int) bool {
	return id >= 0 && id < 4096
}

func ConvertUnixToLocalDateString(timestamp int64) string {
	timeObj := time.Unix(timestamp, 0)
	local := timeObj.Local()
	return local.Format(time.RFC822)
}

func GetPasswordPepper() string {
	return readFromEnvFile("APP_PEPPER")
}

func GetBaseURL() string {
	return readFromEnvFile("BASE_URL")
}

func readFromEnvFile(key string) string {
	env, err := gotenv.Read(".env")
	if err != nil {
		panic(err)
	}

	value := env[key]
	if value != "" {
		return value
	} else {
		panic(key + " env var is missing")
	}
}
