/*
Controller package is the interface of the module.
It acts as the input receiver for other services or for external users.
*/
package controller

import (
	"database/sql"
	"errors"

	"github.com/blocklords/gosds/account"
	"github.com/blocklords/gosds/argument"
	"github.com/blocklords/gosds/env"
	"github.com/blocklords/gosds/message"

	zmq "github.com/pebbe/zmq4"
)

type CommandHandlers map[string]interface{}

// Creates a new Reply controller using ZeroMQ
// The requesters is the list of curve public keys that are allowed to connect to the socket.
func ReplyController(db *sql.DB, commands CommandHandlers, e *env.Env, accounts account.Accounts) error {
	if !e.PortExist() {
		return errors.New("missing necessary environment variables. Please set '" + e.ServiceName() + "_PORT' and/or '" + e.ServiceName() + "_PUBLIC_KEY', '" + e.ServiceName() + "_SECRET_KEY'")
	}

	exist, err := argument.Exist(argument.PLAIN)
	if err != nil {
		return err
	}

	if !exist {
		// only whitelisted users are allowed
		zmq.AuthCurveAdd("*", accounts.PublicKeys()...)
	}

	// Socket to talk to clients
	socket, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		return err
	} else {
		defer socket.Close()
	}

	if !exist {
		err = socket.ServerAuthCurve(e.DomainName(), e.SecretKey())
		if err != nil {
			return err
		}
	}

	if err := socket.Bind("tcp://*:" + e.Port()); err != nil {
		return errors.New("error to bind socket for '" + e.ServiceName() + " - " + e.Url() + "' : " + err.Error())
	}

	println("'" + e.ServiceName() + "' request-reply server runs on port " + e.Port())

	for {
		// msg_raw, metadata, err := socket.RecvMessageWithMetadata(0, "pub_key")
		msg_raw, err := socket.RecvMessage(0)
		if err != nil {
			fail := message.Fail("socket error to receive message " + err.Error())
			reply := fail.ToString()
			if _, err := socket.SendMessage(reply); err != nil {
				return errors.New("failed to reply: %w" + err.Error())
			}
			continue
		}

		// All request types derive from the basic request.
		// We first attempt to parse basic request from the raw message
		request, err := message.ParseRequest(msg_raw)
		if err != nil {
			fail := message.Fail("invalid json request: " + err.Error())
			reply := fail.ToString()
			if _, err := socket.SendMessage(reply); err != nil {
				return errors.New("failed to reply: %w" + err.Error())
			}
			continue
		}

		// Any request types is compatible with the Request.
		if commands[request.Command] == nil {
			fail := message.Fail("unsupported command " + request.Command)
			reply := fail.ToString()
			if _, err := socket.SendMessage(reply); err != nil {
				return errors.New("failed to reply: %w" + err.Error())
			}
			continue
		}

		var reply message.Reply

		// The command might be from a smartcontract developer.
		command_handler, ok := commands[request.Command].(func(*sql.DB, message.SmartcontractDeveloperRequest, *account.SmartcontractDeveloper) message.Reply)
		if ok {
			smartcontract_developer_request, err := message.ParseSmartcontractDeveloperRequest(msg_raw)
			if err != nil {
				fail := message.Fail("invalid smartcontract developer request " + err.Error())
				reply := fail.ToString()
				if _, err := socket.SendMessage(reply); err != nil {
					return errors.New("failed to reply: %w" + err.Error())
				}
				continue
			}

			smartcontract_developer, err := account.NewSmartcontractDeveloper(&smartcontract_developer_request)
			if err != nil {
				println(smartcontract_developer_request.NonceTimestamp)
				fail := message.Fail("reply controller error as invalid smartcontract developer request: " + err.Error())
				reply := fail.ToString()
				if _, err := socket.SendMessage(reply); err != nil {
					return errors.New("failed to reply: %w" + err.Error())
				}
				continue
			}

			reply = command_handler(db, smartcontract_developer_request, smartcontract_developer)
		} else {
			// The command might be from another SDS Service
			service_handler, ok := commands[request.Command].(func(*sql.DB, message.ServiceRequest, *account.Account) message.Reply)
			if ok {
				service_request, err := message.ParseServiceRequest(msg_raw)
				if err != nil {
					fail := message.Fail("invalid service request " + err.Error())
					reply := fail.ToString()
					if _, err := socket.SendMessage(reply); err != nil {
						return errors.New("failed to reply: %w" + err.Error())
					}
					continue
				}

				service_account := account.NewService(service_request.Service)

				reply = service_handler(db, service_request, service_account)
			} else {
				// The command is from a developer.
				reply = commands[request.Command].(func(*sql.DB, message.Request) message.Reply)(db, request)
			}
		}

		if _, err := socket.SendMessage(reply.ToString()); err != nil {
			return errors.New("failed to reply: %w" + err.Error())
		}
	}
}
