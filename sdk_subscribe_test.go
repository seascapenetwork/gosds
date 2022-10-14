/*-------------------------------------------------------------

SUBSCRIBE

-------------------------------------------------------------
*/
package subscribe_test

import (
	"github.com/blocklords/gosds/topic"
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/sdk"
)

func test() {
	subscriber := sdk.Subscriber.FromEnvKey("MY_PRIVATE_KEY") // subscriber.Subscriber
	subTopic := topic.ParseString("seascape-network.blocklords")

	channel := make(chan message.Subscription)

	go subscriber.listen(channel)

	for {
		event := <-channel //message.Subscription{}
		// todo with event

		break
	}
}

test()