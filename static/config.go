package static

import (
	"github.com/blocklords/gosds/topic"

	zmq "github.com/pebbe/zmq4"
)

type Config struct {
}

func ByTopic(sock *zmq.Socket, t *topic.Topic) Config {
	return Config{}
}
