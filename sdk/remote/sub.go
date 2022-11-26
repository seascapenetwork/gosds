// This package establishes a short living connection to the SDS Server
package remote

import (
	zmq "github.com/pebbe/zmq4"
)

// Request Reply pattern. In the web it's called RPC.
func NewSub(host string, topic string) (*zmq.Socket, error) {
	socket, sockErr := zmq.NewSocket(zmq.SUB)
	if sockErr != nil {
		return nil, sockErr
	}

	conErr := socket.Connect("tcp://" + host)
	if conErr != nil {
		return nil, conErr
	}

	err := socket.SetSubscribe(topic)
	if err != nil {
		return nil, err
	}

	return socket, nil
}

func NewSocket(host string) (*zmq.Socket, error) {
	socket, sockErr := zmq.NewSocket(zmq.REQ)
	if sockErr != nil {
		return nil, sockErr
	}

	conErr := socket.Connect("tcp://" + host)
	if conErr != nil {
		return nil, conErr
	}

	return socket, nil
}
