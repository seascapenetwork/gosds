package env

import (
	"fmt"
)

func Require(names []string) error {
	for _, n := range names {
		v := GetString(n)
		if len(v) == 0 {
			return fmt.Errorf("The environment variable is missing: %s", n)
		}
	}
	return nil
}
