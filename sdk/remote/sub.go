// This package establishes a short living connection to the SDS Server
package remote

import (
	zmq "github.com/pebbe/zmq4"
)

// Subscriber. In the web it's called RPC.
// If the topic parameter is empty, then user has to set it up by himself.
func NewSub(host string, topic string) (*zmq.Socket, error) {
	socket, sockErr := zmq.NewSocket(zmq.SUB)
	if sockErr != nil {
		return nil, sockErr
	}

	conErr := socket.Connect("tcp://" + host)
	if conErr != nil {
		return nil, conErr
	}

	if topic != "" {
		err := socket.SetSubscribe(topic)
		if err != nil {
			return nil, err
		}
	}

	return socket, nil
}
