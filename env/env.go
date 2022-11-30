/*
The environment package's file category handles loading

.env or any other environment variable that is provided by the user
*/
package env

import (
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Environment variables for each SDS Service
type Env struct {
	service        string // Service name
	broadcast_host string // Broadcasting host
	broadcast_port string // Broadcasting port
	host           string // request-reply host
	port           string // request-reply port
}

// Checks whether the envrionment variable exists or not
func Exists(name string) bool {
	_, exists := os.LookupEnv(name)
	return exists
}

// Returns an envrionment variable as a string
func GetString(name string) string {
	return os.Getenv(name)
}

// Returns an envrionment variable as a number
func GetNumeric(name string) uint {
	value := os.Getenv(name)
	if len(value) == 0 {
		return 0
	}

	num, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0
	}

	return uint(num)
}

// Returns the envrionment variable for the SDS Spaghetti
func Spaghetti() *Env { return Get("SPAGHETTI") }

// Returns the envrionment variable for the SDS Categorizer
func Categorizer() *Env { return Get("CATEGORIZER") }

// Returns the envrionment variable for the SDS Static
func Static() *Env { return Get("STATIC") }

// Returns the envrionment variable for the SDS Gateway
func Gateway() *Env { return Get("GATEWAY") }

// Returns the envrionment variable for the SDS Publisher
func Publisher() *Env { return Get("PUBLISHER") }

// Returns the envrionment variable for the SDS Reader
func Reader() *Env { return Get("READER") }

// Returns the envrionment variable for the SDS Writer
func Writer() *Env { return Get("WRITER") }

// Returns the envrionment variable for the SDS Bundler
func Bundle() *Env { return Get("BUNDLE") }

// Returns the envrionment variable for the SDS Log
func Log() *Env { return Get("LOG") }

// for example env.Get("SPAGHETTI").
func Get(service string) *Env {
	host := GetString(service + "_HOST")
	port := GetString(service + "_PORT")
	broadcast_host := GetString(service + "_BROADCAST_HOST")
	broadcast_port := GetString(service + "_BROADCAST_PORT")

	return &Env{service: service, host: host, port: port, broadcast_host: broadcast_host, broadcast_port: broadcast_port}
}

// Returns the Service Name
func (e *Env) ServiceName() string {
	caser := cases.Title(language.AmericanEnglish)
	return "SDS " + caser.String(strings.ToLower(e.service))
}

// Returns the request-reply url as a host:port
func (e *Env) Url() string {
	return e.host + ":" + e.port
}

// Returns the broadcast url as a host:port
func (e *Env) BroadcastUrl() string {
	return e.broadcast_host + ":" + e.broadcast_port
}

// returns the request-reply port
func (e *Env) Port() string {
	return e.port
}

// Returns the request-reply port environment variable
func (e *Env) PortEnv() string {
	return GetString(e.service + "_PORT")
}

// Returns the request-reply host
func (e *Env) Host() string {
	return e.host
}

// Returns the request-reply host environment variable
func (e *Env) HostEnv() string {
	return GetString(e.service + "_HOST")
}

// Returns the broadcast host
func (e *Env) BroadcastHost() string {
	return e.broadcast_host
}

// Returns the broadcast host environment variable
func (e *Env) BroadcastHostEnv() string {
	return GetString(e.service + "_BROADCAST_HOST")
}

// Returns the broadcast port
func (e *Env) BroadcastPort() string {
	return e.broadcast_port
}

// Returns the broadcast port environment variable
func (e *Env) BroadcastPortEnv() string {
	return GetString(e.service + "_BROADCAST_PORT")
}

// Checks whether the request-reply's host and port exists
func (e *Env) UrlExist() bool {
	return len(e.port) > 0 && len(e.host) > 0
}

// Checks whether the port exists
func (e *Env) PortExist() bool {
	return len(e.port) > 0
}

// Checks whether the broadcast host and port exists
func (e *Env) BroadcastExist() bool {
	return len(e.broadcast_host) > 0 && len(e.broadcast_port) > 0
}
