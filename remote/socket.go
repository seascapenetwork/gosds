// This package defines a data types, functions that interacts with a remote SDS service.
package remote

import (
	"fmt"

	"github.com/blocklords/gosds/env"
	"github.com/blocklords/gosds/message"
	zmq "github.com/pebbe/zmq4"
)

// Over the socket the remote call is happening.
// its a wrapper of zeromq socket.
type Socket struct {
	// The name of remote SDS service.
	// its used as a clarification
	RemoteServiceName string
	socket            *zmq.Socket
}

// Close the remote connection
func (socket *Socket) Close() {
	socket.socket.Close()
}

// Send a command to the remote SDS service.
// Note that it converts the failure reply into an error. Rather than replying reply itself back to user.
// In case of successful request, the function returns reply parameters.
func (socket *Socket) RequestRemoteService(request *message.Request) (map[string]interface{}, error) {
	if _, err := socket.socket.SendMessage(request.ToString()); err != nil {
		return nil, fmt.Errorf("failed to send the command '%s' to '%s'. socket error: %w", request.Command, socket.RemoteServiceName, err)
	}

	// Wait for reply.
	r, err := socket.socket.RecvMessage(0)
	if err != nil {
		return nil, fmt.Errorf("failed to receive the command '%s' message from '%s'. socket error: %w", request.Command, socket.RemoteServiceName, err)
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the command '%s' reply from '%s'. gosds error %w", request.Command, socket.RemoteServiceName, err)
	}
	if !reply.IsOK() {
		return nil, fmt.Errorf("the command '%s' failed by '%s'. the reply error message: %s", request.Command, socket.RemoteServiceName, reply.Message)
	}

	return reply.Params, nil
}

func TcpRequestSocketOrPanic(e *env.Env) *Socket {
	if !e.UrlExist() {
		panic(fmt.Errorf("missing .env variable: Please set '" + e.ServiceName() + "' host and port"))
	}

	sock, _ := zmq.NewSocket(zmq.REQ)
	if err := sock.Connect("tcp://" + e.Url()); err != nil {
		panic(fmt.Errorf("error '"+e.ServiceName()+"' connect: %w", err))
	}

	return &Socket{
		RemoteServiceName: e.ServiceName(),
		socket:            sock,
	}
}
