package env

import "os"

func GetString(name string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		println("no " + name + "' environment variable set")
	}
	return value
}

// Returns the path of Spaghetti Publisher as host:port
func SpaghettiPublisher() string {
	return GetString("SPAGHETTI_PUBLISHER_HOST") + ":" + SpaghettiPublisherPort()
}

func SpaghettiPublisherPort() string {
	return GetString("SPAGHETTI_PUBLISHER_PORT")
}

// Returns the path of Spaghetti Controller as host:port
func SpaghettiController() string {
	return GetString("SPAGHETTI_INTERNAL_HOST") + ":" + SpaghettiControllerPort()
}

func SpaghettiControllerPort() string {
	return GetString("SPAGHETTI_INTERNAL_PORT")
}

func CategorizerController() string {
	return GetString("CATEGORIZER_HOST") + ":" + SpaghettiControllerPort()
}

func CategorizerControllerPort() string {
	return GetString("CATEGORIZER_PORT")
}

func CategorizerPublisher() string {
	return GetString("CATEGORIZER_PUBLISHER_HOST") + ":" + SpaghettiControllerPort()
}

func CategorizerPublisherPort() string {
	return GetString("CATEGORIZER_PUBLISHER_PORT")
}

func StaticController() string {
	return GetString("STATIC_HOST") + ":" + StaticPort()
}

func StaticPort() string {
	return GetString("STATIC_PORT")
}

func Gateway() string {
	return GetString("GATEWAY") + ":" + GatewayPort()
}

func GatewayPort() string {
	return GetString("GATEWAY_PORT")
}
