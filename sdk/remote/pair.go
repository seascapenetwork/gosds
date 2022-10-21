// This package establishes a short living connection to the SDS Server
package remote

import (
	zmq "github.com/pebbe/zmq4"
)

// Request Reply pattern. In the web it's called RPC.
func NewPair(host string, sender bool) (*zmq.Socket, error) {
	//  Bind inproc socket before starting step1
	socket, err := zmq.NewSocket(zmq.PAIR)
	if err != nil {
		return nil, err
	}
	defer socket.Close()
	if sender {
		socket.Bind("inproc://pair_" + host)
	} else {
		socket.Connect("inproc://pair_" + host)
	}
	return socket, nil
}
