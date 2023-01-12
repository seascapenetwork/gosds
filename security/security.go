// The security package has the utility functions to make the
// Requests secure.
package security

const PLAIN = "plain"

// Whether the socket connections are in plain or not.
// Accepts the list of the arguments passed to the application.
// If one of the arguments is "--plain", then this application runs in a plain mode.
func IsPlain(arguments []string) bool {
	for _, arg := range arguments {
		if arg == PLAIN {
			return true
		}
	}

	return false
}
