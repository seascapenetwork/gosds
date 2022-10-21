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
	defer socket.Close()

	conErr := socket.Connect(host)
	if conErr != nil {
		return nil, conErr
	}

	err := socket.SetSubscribe(topic)
	if err != nil {
		return nil, err
	}

	return socket, nil
}
