package env

import (
	"github.com/blocklords/gosds/argument"
	"github.com/joho/godotenv"
)

// Load all .env files
func LoadAnyEnv() error {
	opts, optErr := argument.EnvPaths()
	if optErr != nil {
		return optErr
	}

	godotenv.Load()

	if opts != nil {
		return godotenv.Load(opts...)
	}
	return nil
}
