package env

import (
	"github.com/pkg/errors"
	"os"
)

var (
	MODE       string
	DB_URL     string
	JWT_SECRET string
	REDIS_URL  string
)

func Load() error {
	var ok bool

	MODE, ok = os.LookupEnv("MODE")
	if !ok {
		return errors.New("MODE is not set")
	}

	DB_URL, ok = os.LookupEnv("DB_URL")
	if !ok {
		return errors.New("DB_URL is not set")
	}

	if MODE == "production" {
		JWT_SECRET, ok = os.LookupEnv("JWT_SECRET")
		if !ok {
			return errors.New("JWT_SECRET is not set")
		}

		REDIS_URL, ok = os.LookupEnv("REDIS_URL")
		if !ok {
			return errors.New("REDIS_URL is not set")
		}
	}

	return nil
}
