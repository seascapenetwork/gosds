/*Controller package is the interface of the module.
It acts as the input receiver for other services.*/
package controller

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/blocklords/gosds/message"

	zmq "github.com/pebbe/zmq4"
)

type CommandHandlers map[string]interface{}

/*
Creates a new Reply controller using ZeroMQ
*/
func ReplyController(db *sql.DB, commands CommandHandlers, port string) {
	// Socket to talk to clients
	socket, _ := zmq.NewSocket(zmq.REP)
	defer socket.Close()
	if err := socket.Bind("tcp://*:" + port); err != nil {
		println(fmt.Errorf("listening: %w", err))
	}

	println("Waiting for commands on port `3013`")

	for {
		msg_raw, err := socket.RecvMessage(0)
		if err != nil {
			println(fmt.Errorf("receiving: %w", err))
		}
		request, err := message.ParseRequest(msg_raw)
		if err != nil {
			fail := message.Fail("invalid request " + err.Error())
			reply := fail.ToString()
			if _, err := socket.SendMessage(reply); err != nil {
				println(fmt.Errorf("sending reply: %w", err))
			}
			continue
		}

		if commands[request.Command] == nil {
			fail := message.Fail("invalid command " + request.Command)
			reply := fail.ToString()
			if _, err := socket.SendMessage(reply); err != nil {
				println(fmt.Errorf(" reply: %w", err))
			}
			continue
		}

		reply := commands[request.Command].(func(message.Request) message.Reply)(request)

		// Do some 'work'
		time.Sleep(time.Second)

		if _, err := socket.SendMessage(reply.ToString()); err != nil {
			println(fmt.Errorf("sending reply: %w", err))
		}
	}
}
