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

	paths := make([]string, 0)

	for _, arg := range args {
		if arg[:2] != "--" {
			paths = append(paths, arg)
		}
	}

	return paths, nil
}

// Load arguments, not the environment variable paths.
// Arguments are with --prefix
func LoadArguments() ([]string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, nil
	}

	parameters := make([]string, 0)

	for _, arg := range args {
		if arg[:2] == "--" {
			parameters = append(parameters, arg[2:])
		}
	}

	return parameters, nil
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
