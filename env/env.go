/*
The environment package's file category handles loading

.env or any other environment variable that is provided by the user
*/
package env

import (
	"os"
	"strconv"
	"strings"
)

type Env struct {
	service       string
	broadcastHost string
	broadcastPort string
	port          string
	host          string
}

// Checks whether the envrionment variable exists or not
func Exists(name string) bool {
	_, exists := os.LookupEnv(name)
	return exists
}

func GetString(name string) string {
	return os.Getenv(name)
}

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

func Spaghetti() *Env   { return Get("SPAGHETTI") }
func Categorizer() *Env { return Get("CATEGORIZER") }
func Static() *Env      { return Get("STATIC") }
func Gateway() *Env     { return Get("GATEWAY") }
func Publisher() *Env   { return Get("PUBLISHER") }
func Reader() *Env      { return Get("READER") }
func Writer() *Env      { return Get("WRITER") }
func Bundle() *Env      { return Get("BUNDLE") }
func Log() *Env         { return Get("LOG") }

// for example env.Get("SPAGHETTI").
func Get(service string) *Env {
	host := GetString(service + "_HOST")
	port := GetString(service + "_PORT")
	broadcastHost := GetString(service + "_BROADCAST_HOST")
	broadcastPort := GetString(service + "_BROADCAST_PORT")

	return &Env{service: service, host: host, port: port, broadcastHost: broadcastHost, broadcastPort: broadcastPort}
}

func (e *Env) ServiceName() string {
	return "SDS " + strings.Title(strings.ToLower(e.service))
}
func (e *Env) Url() string {
	return e.host + ":" + e.port
}

func (e *Env) BroadcastUrl() string {
	return e.broadcastHost + ":" + e.broadcastPort
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

func (e *Env) BroadcastHost() string {
	return e.broadcastHost
}

func (e *Env) BroadcastHostEnv() string {
	return GetString(e.service + "_BROADCAST_HOST")
}

func (e *Env) BroadcastPort() string {
	return e.broadcastPort
}

func (e *Env) BroadcastPortEnv() string {
	return GetString(e.service + "_BROADCAST_PORT")
}

func (e *Env) UrlExist() bool {
	return len(e.port) > 0 && len(e.host) > 0
}

func (e *Env) PortExist() bool {
	return len(e.port) > 0
}

func (e *Env) BroadcastExist() bool {
	return len(e.broadcastHost) > 0 && len(e.broadcastPort) > 0
}
