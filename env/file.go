package env

import (
	"os"

	"github.com/joho/godotenv"
)

// any command line data that comes after the files are .env file paths
func optional_paths() ([]string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, nil
	}

	return args, nil
}

// Load all .env files
func LoadAnyEnv() error {
	opts, optErr := optional_paths()
	if optErr != nil {
		return optErr
	}

	godotenv.Load()

	if opts != nil {
		return godotenv.Load(opts...)
	}
	return nil
}
