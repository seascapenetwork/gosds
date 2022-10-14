/*-------------------------------------------------------------

SIGN

-------------------------------------------------------------
*/
package subscribe_test

import (
	"github.com/blocklords/gosds/topic"
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/sdk"
)

func test() {
	signer := sdk.Signer.FromEnvKey("MY_PRIVATE_KEY") // subscriber.Subscriber
	signTopic := topic.ParseString("seascape-network.blocklords.11155111.core.ImportExportElastic.export")

	signParams := map[string]interface{}{
		"_greeting": "Hello and welcome",
	}

	repl, err := signer.Exec(signTopic, signParams) // repl is message.Reply
	if err != nil {
		fmt.Println("Error ", err.Error())
	}

	fmt.Println("Reply ", repl.ToString())
}

test()