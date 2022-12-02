/*
Controller package is the interface of the module.
It acts as the input receiver for other services or for external users.
*/
package controller

import (
	"database/sql"
	"fmt"

	"github.com/blocklords/gosds/env"
	"github.com/blocklords/gosds/message"

	zmq "github.com/pebbe/zmq4"
)

type CommandHandlers map[string]interface{}

/*
Creates a new Reply controller using ZeroMQ
*/
func ReplyController(db *sql.DB, commands CommandHandlers, e *env.Env) {
	if !e.PortExist() {
		panic(fmt.Errorf("missing .env variable: Please set '" + e.ServiceName() + "' port"))
	}

	// Socket to talk to clients
	socket, _ := zmq.NewSocket(zmq.REP)
	defer socket.Close()
	if err := socket.Bind("tcp://*:" + e.Port()); err != nil {
		println("error to bind socket for '"+e.ServiceName()+" - "+e.Url()+"' : ", err.Error())
		panic(err)
	}

	println("'" + e.ServiceName() + "' request-reply server runs on port " + e.Port())

	for {
		msg_raw, err := socket.RecvMessage(0)
		if err != nil {
			println(fmt.Errorf("receiving: %w", err))
			continue
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

		reply := commands[request.Command].(func(*sql.DB, message.Request) message.Reply)(db, request)

		if _, err := socket.SendMessage(reply.ToString()); err != nil {
			println(fmt.Errorf("error sending controller reply: %w", err))
		}
	}
}
