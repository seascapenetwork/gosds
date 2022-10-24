/*The gosds/sdk package is the client package to interact with SDS.
The following commands are available in this SDK:

1. Subscribe - subscribe for events
2. Sign - send a transaction to the blockchain
3. AddToPool - send a transaction to the pool that will be broadcasted to the blockchain bundled.
4. Read - read a smartcontract information


Usage

----------------------------------------------------------------
example of reading smartcontract data

   import (
	"github.com/blocklords/gosds/sdk"
	"github.com/blocklords/gosds/topic"
   )

   func test() {
	// returns sdk.Reader
	reader := sdk.NewReader("address", "gateway host")
	// gosds.topic.Topic
	importAddressTopic := topic.ParseString("metaking.blocklords.11155111.transfer.ImportExportManager.accountHodlerOf")
	args := ["user address"]

	// returns gosds.message.Reply
	reply := reader.Read(importAddressTopic, args)

	if !reply.IsOk() {
		panic(fmt.Errorf("failed to read smartcontract data: %w", reply.Message))
	}

	fmt.Println("The user's address is: ", reply.Params["result"].(string))
   }
*/
package sdk

import (
	"github.com/blocklords/gosds/sdk/reader"
	"github.com/blocklords/gosds/sdk/subscriber"
	"github.com/blocklords/gosds/sdk/writer"
)

var Version string = "Seascape GoSDS version: 0.0.8"

// Returns a new reader.Reader.
//
// The host is the link to the SDS Gateway.
// The address argument is the wallet address that is allowed to read.
func NewReader(host string, address string) *reader.Reader {
	return reader.NewReader(host, address)
}

func NewWriter(host string, address string) *writer.Writer {
	return writer.NewWriter(host, address)
}

// Returns a new subscriber
func NewSubscriber(host string, sub string, address string) *subscriber.Subscriber {
	return subscriber.NewSubscriber(host, sub, address)
}
