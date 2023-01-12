// The security package has the utility functions to make the
// Requests secure.
package security

import (
	"github.com/blocklords/gosds/env"
)

const PLAIN = "plain"

// Whether the socket connections are in plain or not.
// Accepts the list of the arguments passed to the application.
// If one of the arguments is "--plain", then this application runs in a plain mode.
func IsPlain() (bool, error) {
	arguments, err := env.LoadArguments()
	if err != nil {
		return false, err
	}

	for _, arg := range arguments {
		if arg == PLAIN {
			return true, nil
		}
	}

	return false, nil
}
