/*
The environment package's file category handles loading

.env or any other environment variable that is provided by the user
*/
package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Environment variables for each SDS Service
type Env struct {
	service              string // Service name
	broadcast_host       string // Broadcasting host
	broadcast_port       string // Broadcasting port
	host                 string // request-reply host
	port                 string // request-reply port
	public_key           string // The Curve key of the service
	secret_key           string // The Curve secret key of the service
	broadcast_public_key string
	broadcast_secret_key string
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

// Returns the envrionment variable for the SDS Developer Gateway
func DeveloperGateway() *Env { return Get("DEVELOPER_GATEWAY") }

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

func NewDeveloper(public_key string, secret_key string) *Env {
	return &Env{
		service:        "developer",
		host:           "",
		port:           "",
		broadcast_host: "",
		broadcast_port: "",
		public_key:     public_key,
		secret_key:     secret_key,
	}
}

func Developer() *Env {
	public_key := GetString("DEVELOPER_PUBLIC_KEY")
	secret_key := GetString("DEVELOPER_SECRET_KEY")

	return NewDeveloper(public_key, secret_key)
}

// for example env.Get("SPAGHETTI").
func Get(service string) *Env {
	host := GetString(service + "_HOST")
	port := GetString(service + "_PORT")
	broadcast_host := GetString(service + "_BROADCAST_HOST")
	broadcast_port := GetString(service + "_BROADCAST_PORT")
	public_key := GetString(service + "_PUBLIC_KEY")
	secret_key := GetString(service + "_SECRET_KEY")
	broadcast_public_key := GetString(service + "_BROADCAST_PUBLIC_KEY")
	broadcast_secret_key := GetString(service + "_BROADCAST_SECRET_KEY")

	return &Env{
		service:              service,
		host:                 host,
		port:                 port,
		broadcast_host:       broadcast_host,
		broadcast_port:       broadcast_port,
		public_key:           public_key,
		secret_key:           secret_key,
		broadcast_public_key: broadcast_public_key,
		broadcast_secret_key: broadcast_secret_key,
	}
}

// Returns the service environment parameters by its Public Key
func GetByPublicKey(public_key string) (*Env, error) {
	services := []string{
		"SPAGHETTI",
		"CATEGORIZER",
		"STATIC",
		"GATEWAY",
		"PUBLISHER",
		"READER",
		"WRITER",
		"BUNDLE",
		"LOG",
		"DEVELOPER_GATEWAY",
	}

	for _, service := range services {
		service_env := Get(service)
		if service_env != nil && service_env.public_key == public_key {
			return service_env, nil
		}
	}

	return nil, errors.New("the service wasn't found for a given public key")
}

func (e *Env) SecretKey() string {
	return e.secret_key
}

func (e *Env) PublicKey() string {
	return e.public_key
}

func (e *Env) BroadcastSecretKey() string {
	return e.broadcast_secret_key
}

func (e *Env) BroadcastPublicKey() string {
	return e.broadcast_public_key
}

func (e *Env) DomainName() string {
	return e.service
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

func (e *Env) BroadcastPortExists() bool {
	return len(e.broadcast_port) > 0 && len(e.broadcast_public_key) > 0 && len(e.broadcast_secret_key) > 0
}

// Returns the broadcast port environment variable
func (e *Env) BroadcastPortEnv() string {
	return GetString(e.service + "_BROADCAST_PORT")
}

// Checks whether the request-reply's host and port exists
func (e *Env) UrlExist() bool {
	return len(e.port) > 0 && len(e.host) > 0 && len(e.public_key) > 0
}

// Checks whether the port exists
func (e *Env) PortExist() bool {
	return len(e.port) > 0 && len(e.public_key) > 0 && len(e.secret_key) > 0
}

// Checks whether the broadcast host and port exists
func (e *Env) BroadcastExist() bool {
	return len(e.broadcast_host) > 0 && len(e.broadcast_port) > 0 && len(e.broadcast_public_key) > 0
}

// Necessary environment variables, to subscribe to the SDS Service
func (e *Env) ToSubscribe() *Env {
	if !e.BroadcastExist() {
		service := e.DomainName()
		broadcast_host := service + "_BROADCAST_HOST"
		broadcast_port := service + "_BROADCAST_PORT"
		public_key := service + "_BROADCAST_PUBLIC_KEY"

		panic(fmt.Sprintf("the '%s' service couldn't be built. missing: '%s', '%s', '%s'",
			e.ServiceName(), broadcast_host, broadcast_port, public_key))
	}

	return e
}

// Necessary environment variables, to request to SDS Service.
func (e *Env) ToRequest() *Env {
	if !e.UrlExist() {
		service := e.DomainName()
		host := service + "_HOST"
		port := service + "_PORT"
		public_key := service + "_PUBLIC_KEY"

		panic(fmt.Sprintf("the '%s' service couldn't be built. missing: '%s', '%s', '%s'",
			e.ServiceName(), host, port, public_key))
	}

	return e
}

// Necessary environment variables, to broadcast by SDS Service.
// Panic otherwise.
func (e *Env) ToBroadcast() *Env {
	if !e.BroadcastPortExists() {
		service := e.DomainName()
		port := service + "_BROADCAST_PORT"
		public_key := service + "_BROADCAST_PUBLIC_KEY"
		secret_key := service + "_BROADCAST_SECRET_KEY"

		panic(fmt.Sprintf("the '%s' service couldn't be built. missing: '%s', '%s', '%s'",
			e.ServiceName(), port, public_key, secret_key))
	}

	return e
}

// Necessary environment variables, to reply by SDS Service.
// Panic otherwise.
func (e *Env) ToReply() *Env {
	if !e.BroadcastPortExists() {
		service := e.DomainName()
		port := service + "_PORT"
		public_key := service + "_PUBLIC_KEY"
		secret_key := service + "_SECRET_KEY"

		panic(fmt.Sprintf("the '%s' service couldn't be built. missing: '%s', '%s', '%s'",
			e.ServiceName(), port, public_key, secret_key))
	}

	return e
}
