/* The environment package's file category handles loading
.env or any other environment variable that is provided by the user
*/
package env

import (
	"os"

	"github.com/joho/godotenv"
)

func optionalPaths() ([]string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, nil
	}

	return args, nil
}

func LoadAnyEnv() error {
	opts, optErr := optionalPaths()
	if optErr != nil {
		return optErr
	}

	godotenv.Load()

	if opts != nil {
		return godotenv.Load(opts...)
	}
	return nil
}
