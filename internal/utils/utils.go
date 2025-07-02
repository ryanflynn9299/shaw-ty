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
	env, err := gotenv.Read(".env")
	if err != nil {
		panic(err)
	}

	pepper := env["APP_PEPPER"]
	if pepper != "" {
		return pepper
	} else {
		panic("APP_PEPPER env var is missing")
	}
}
