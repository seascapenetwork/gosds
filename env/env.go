package env

import (
	"os"
	"strconv"
)

type Env struct {
	service       string
	publisherHost string
	publisherPort string
	port          string
	host          string
}

func GetString(name string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		println("no " + name + "' environment variable set")
	}
	return value
}

func GetNumeric(name string) uint {
	value := os.Getenv(name)
	if len(value) == 0 {
		println("no " + name + "' environment variable set")
		return 0
	}

	num, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		println("invalid number format " + err.Error())
		return 0
	}

	return uint(num)
}

func Spaghetti() *Env   { return Get("SPAGHETTI") }
func Categorizer() *Env { return Get("CATEGORIZER") }
func Static() *Env      { return Get("STATIC") }
func Gateway() *Env     { return Get("GATEWAY") }
func Publisher() *Env   { return Get("PUBLISHER") }
func Reader() *Env      { return Get("READER") }

// env.Get("SPAGHETTI").
func Get(service string) *Env {
	host := GetString(service + "_HOST")
	port := GetString(service + "_PORT")
	publisherHost := GetString(service + "_PUBLISHER_HOST")
	publisherPort := GetString(service + "_PUBLISHER_PORT")

	return &Env{service: service, host: host, port: port, publisherHost: publisherHost, publisherPort: publisherPort}
}

func (e *Env) Url() string {
	return e.host + ":" + e.port
}

func (e *Env) PublisherUrl() string {
	return e.publisherHost + ":" + e.publisherPort
}

func (e *Env) Port() string {
	return e.port
}

func (e *Env) PortEnv() string {
	return GetString(e.service + "_PORT")
}

func (e *Env) Host() string {
	return e.host
}

func (e *Env) HostEnv() string {
	return GetString(e.service + "_HOST")
}

func (e *Env) PublisherHost() string {
	return e.publisherHost
}

func (e *Env) PublisherHostEnv() string {
	return GetString(e.service + "_PUBLISHER_HOST")
}

func (e *Env) PublisherPort() string {
	return e.publisherPort
}

func (e *Env) PublisherPortEnv() string {
	return GetString(e.service + "_PUBLISHER_PORT")
}

func (e *Env) UrlExist() bool {
	return len(e.port) > 0 && len(e.host) > 0
}

func (e *Env) PortExist() bool {
	return len(e.port) > 0
}

func (e *Env) PublisherExist() bool {
	return len(e.publisherHost) > 0 && len(e.publisherPort) > 0
}
