// This package defines the data types, and methods that interact with a remote SDS service.
//
// The request reply socket follows the Lazy Pirate pattern.
//
// Example using pebbe/zmq4 is here:
// https://github.com/pebbe/zmq4/blob/83013091510dd1275bbf0b9a302533cadc17d392/examples/lpclient.go
//
// The Lazy Pirate pattern is described in the ZMQ guide:
// https://zguide.zeromq.org/docs/chapter4/#Client-Side-Reliability-Lazy-Pirate-Pattern
package remote

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/blocklords/gosds/env"
	"github.com/blocklords/gosds/message"
	zmq "github.com/pebbe/zmq4"
)

// Over the socket the remote call is happening.
// This is the wrapper of zeromq socket. Wrapper enables to create larger network patterns.
type Socket struct {
	// The name of remote SDS service and its URL
	// its used as a clarification
	remoteService *env.Env
	socket        *zmq.Socket
}

type SDS_Message interface {
	*message.Request | *message.ServiceRequest

	CommandName() string
	ToString() string
}

// Request-Reply checks the internet connection after this amount of time.
const (
	REQUEST_TIMEOUT = 60 * time.Second //  msecs, (> 1000!)
)

// Close the remote connection
func (socket *Socket) Close() error {
	socket_type, err := socket.socket.GetType()
	if err != nil {
		return err
	}
	url := ""
	if socket_type == zmq.SUB {
		url = socket.remoteService.BroadcastUrl()
	} else {
		url = socket.remoteService.Url()
	}
	err = socket.socket.Disconnect("tcp://" + url)
	if err != nil {
		return err
	}
	wait := time.Duration(5) * time.Second
	time.Sleep(wait)

	return socket.socket.Close()
}

// Broadcaster URL of the SDS Service
func (socket *Socket) RemoteBroadcastUrl() string {
	return socket.remoteService.BroadcastUrl()
}

// Broadcaster Port of the SDS Service
func (socket *Socket) RemoteBroadcastPort() (uint, error) {
	port, err := strconv.Atoi(socket.remoteService.BroadcastPort())
	if err != nil {
		return 0, err
	}
	return uint(port), nil
}

// Returns the HOST envrionment parameters of the socket.
//
// Use it if you want to create another socket from this socket.
func (socket *Socket) RemoteEnv() *env.Env {
	return socket.remoteService
}

// Send a command to the remote SDS service.
// Note that it converts the failure reply into an error. Rather than replying reply itself back to user.
// In case of successful request, the function returns reply parameters.
func (socket *Socket) RequestRemoteService(request *message.Request) (map[string]interface{}, error) {
	poller := zmq.NewPoller()
	poller.Add(socket.socket, zmq.POLLIN)

	//  We send a request, then we work to get a reply
	if _, err := socket.socket.SendMessage(request.ToString()); err != nil {
		return nil, fmt.Errorf("failed to send the command '%s' to '%s'. socket error: %w", request.Command, socket.remoteService.ServiceName(), err)
	}

	request_timeout := REQUEST_TIMEOUT
	if env.Exists("SDS_REQUEST_TIMEOUT") {
		env_timeout := env.GetNumeric("SDS_REQUEST_TIMEOUT")
		if env_timeout != 0 {
			request_timeout = time.Duration(env_timeout) * time.Second
		}
	}

	// we attempt requests for an infinite amount of time.
	for {
		//  Poll socket for a reply, with timeout
		sockets, err := poller.Poll(request_timeout)
		if err != nil {
			return nil, fmt.Errorf("failed to to send the command '%s' to '%s'. poll error: %w", request.Command, socket.remoteService.ServiceName(), err)
		}

		//  Here we process a server reply and exit our loop if the
		//  reply is valid. If we didn't a reply we close the client
		//  socket and resend the request. We try a number of times
		//  before finally abandoning:

		if len(sockets) > 0 {
			// Wait for reply.
			r, err := socket.socket.RecvMessage(0)
			if err != nil {
				return nil, fmt.Errorf("failed to receive the command '%s' message from '%s'. socket error: %w", request.Command, socket.remoteService.ServiceName(), err)
			}

			reply, err := message.ParseReply(r)
			if err != nil {
				return nil, fmt.Errorf("failed to parse the command '%s' reply from '%s'. gosds error %w", request.Command, socket.remoteService.ServiceName(), err)
			}

			if !reply.IsOK() {
				return nil, fmt.Errorf("the command '%s' replied with a failure by '%s'. the reply error message: %s", request.Command, socket.remoteService.ServiceName(), reply.Message)
			}

			return reply.Params, nil
		} else {
			fmt.Println("command '", request.Command, "' wasn't replied by '", socket.remoteService.ServiceName(), "' in ", request_timeout, ", retrying...")
			//  Old socket is confused; close it and open a new one
			socket.socket.Close()

			socket.socket, _ = zmq.NewSocket(zmq.REQ)
			if err := socket.socket.Connect("tcp://" + socket.remoteService.Url()); err != nil {
				panic(fmt.Errorf("error '"+socket.remoteService.ServiceName()+"' connect: %w", err))
			}

			// Recreate poller for new client
			poller = zmq.NewPoller()
			poller.Add(socket.socket, zmq.POLLIN)

			//  Send request again, on new socket
			//  We send a request, then we work to get a reply
			if _, err := socket.socket.SendMessage(request.ToString()); err != nil {
				return nil, fmt.Errorf("failed to send the command '%s' to '%s'. socket error: %w", request.Command, socket.remoteService.ServiceName(), err)
			}
		}
	}
}

// Requests a message to the remote service.
// The socket parameter is the Request socket from this service.
// The request is the message.
func RequestReply[V SDS_Message](socket *Socket, request V) (map[string]interface{}, error) {
	socket_type, err := socket.socket.GetType()
	if err != nil {
		return nil, err
	}

	if socket_type != zmq.REQ && socket_type != zmq.DEALER {
		return nil, errors.New("invalid socket type for request-reply")
	}

	poller := zmq.NewPoller()
	poller.Add(socket.socket, zmq.POLLIN)

	command_name := request.CommandName()

	//  We send a request, then we work to get a reply
	if _, err := socket.socket.SendMessage(request.ToString()); err != nil {
		return nil, fmt.Errorf("failed to send the command '%s' to '%s'. socket error: %w", command_name, socket.remoteService.ServiceName(), err)
	}

	request_timeout := REQUEST_TIMEOUT
	if env.Exists("SDS_REQUEST_TIMEOUT") {
		env_timeout := env.GetNumeric("SDS_REQUEST_TIMEOUT")
		if env_timeout != 0 {
			request_timeout = time.Duration(env_timeout) * time.Second
		}
	}

	// we attempt requests for an infinite amount of time.
	for {
		//  Poll socket for a reply, with timeout
		sockets, err := poller.Poll(request_timeout)
		if err != nil {
			return nil, fmt.Errorf("failed to to send the command '%s' to '%s'. poll error: %w", command_name, socket.remoteService.ServiceName(), err)
		}

		//  Here we process a server reply and exit our loop if the
		//  reply is valid. If we didn't a reply we close the client
		//  socket and resend the request. We try a number of times
		//  before finally abandoning:

		if len(sockets) > 0 {
			// Wait for reply.
			r, err := socket.socket.RecvMessage(0)
			if err != nil {
				return nil, fmt.Errorf("failed to receive the command '%s' message from '%s'. socket error: %w", command_name, socket.remoteService.ServiceName(), err)
			}

			reply, err := message.ParseReply(r)
			if err != nil {
				return nil, fmt.Errorf("failed to parse the command '%s' reply from '%s'. gosds error %w", command_name, socket.remoteService.ServiceName(), err)
			}

			if !reply.IsOK() {
				return nil, fmt.Errorf("the command '%s' replied with a failure by '%s'. the reply error message: %s", command_name, socket.remoteService.ServiceName(), reply.Message)
			}

			return reply.Params, nil
		} else {
			fmt.Println("command '", command_name, "' wasn't replied by '", socket.remoteService.ServiceName(), "' in ", request_timeout, ", retrying...")
			//  Old socket is confused; close it and open a new one
			socket.socket.Close()

			socket.socket, _ = zmq.NewSocket(zmq.REQ)
			if err := socket.socket.Connect("tcp://" + socket.remoteService.Url()); err != nil {
				panic(fmt.Errorf("error '"+socket.remoteService.ServiceName()+"' connect: %w", err))
			}

			// Recreate poller for new client
			poller = zmq.NewPoller()
			poller.Add(socket.socket, zmq.POLLIN)

			//  Send request again, on new socket
			//  We send a request, then we work to get a reply
			if _, err := socket.socket.SendMessage(request.ToString()); err != nil {
				return nil, fmt.Errorf("failed to send the command '%s' to '%s'. socket error: %w", command_name, socket.remoteService.ServiceName(), err)
			}
		}
	}
}

// Create a new Socket on TCP protocol otherwise exit from the program
// The socket is the wrapper over zmq.REQ
func TcpRequestSocketOrPanic(e *env.Env, client *env.Env) *Socket {
	if !e.UrlExist() {
		panic(fmt.Errorf("missing .env variable: Please set '" + e.ServiceName() + "' host and port"))
	}

	sock, _ := zmq.NewSocket(zmq.REQ)
	err := sock.ClientAuthCurve(e.PublicKey(), client.PublicKey(), client.SecretKey())
	if err != nil {
		panic(err)
	}

	if err := sock.Connect("tcp://" + e.Url()); err != nil {
		panic(fmt.Errorf("error '"+e.ServiceName()+"' connect: %w", err))
	}

	return &Socket{
		remoteService: e,
		socket:        sock,
	}
}

// Create a new Socket on TCP protocol otherwise exit from the program
// The socket is the wrapper over zmq.PULL
func TcpPullSocketOrPanic(port uint) *zmq.Socket {
	sock, _ := zmq.NewSocket(zmq.PULL)
	if err := sock.Bind(fmt.Sprintf("tcp://*:%d", port)); err != nil {
		panic(fmt.Errorf("error to create a pull socket at port %d", port))
	}

	return sock
}

// Create a new Socket on TCP protocol otherwise exit from the program
// The socket is the wrapper over zmq.PUSH
func TcpPushSocketOrPanic(port uint) *zmq.Socket {
	sock, _ := zmq.NewSocket(zmq.PUSH)
	if err := sock.Connect(fmt.Sprintf("tcp://localhost:%d", port)); err != nil {
		panic(fmt.Errorf("error to create a push socket at port %d", port))
	}

	return sock
}

// Create a new Socket on TCP protocol otherwise exit from the program
// The socket is the wrapper over zmq.SUB
func TcpSubscriberOrPanic(e *env.Env, client_env *env.Env) *Socket {
	socket, sockErr := zmq.NewSocket(zmq.SUB)
	if sockErr != nil {
		panic(sockErr)
	}

	socket.ClientAuthCurve(e.BroadcastPublicKey(), client_env.BroadcastPublicKey(), client_env.BroadcastSecretKey())

	conErr := socket.Connect("tcp://" + e.BroadcastUrl())
	if conErr != nil {
		panic(conErr)
	}

	return &Socket{
		remoteService: e,
		socket:        socket,
	}
}
