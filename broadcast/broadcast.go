/*Broadcast the new spaghetti update*/
package broadcast

import (
	"log"

	"github.com/blocklords/gosds/env"

	"github.com/blocklords/gosds/message"

	zmq "github.com/pebbe/zmq4"
)

// Run a new broadcaster
//
// It assumes that the another package is starting an authentication layer of zmq:
// ZAP.
func Run(channel chan message.Broadcast, broadcast_env *env.Env, whitelisted_users []*env.Env) {
	public_keys := make([]string, len(whitelisted_users))
	for k, v := range whitelisted_users {
		public_keys[k] = v.BroadcastPublicKey()
	}

	domain_name := broadcast_env.DomainName() + "_broadcast"

	zmq.AuthCurveAdd(domain_name, public_keys...)

	// prepare the publisher
	pub, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		panic("error while trying to create a new socket " + err.Error())
	}
	defer pub.Close()
	pub.ServerAuthCurve(domain_name, broadcast_env.BroadcastSecretKey())

	err = pub.Bind("tcp://*:" + broadcast_env.BroadcastPort())
	if err != nil {
		log.Fatalf("could not listen to publisher: %v", err)
	}

	for {
		broadcast := <-channel

		_, err = pub.SendMessage(broadcast.Topic, broadcast.ToBytes())
		if err != nil {
			log.Fatal(err)
		}
	}
}
