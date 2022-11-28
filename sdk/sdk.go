/*
The gosds/sdk package is the client package to interact with SDS.
The following commands are available in this SDK:

1. Subscribe - subscribe for events
2. Sign - send a transaction to the blockchain
3. AddToPool - send a transaction to the pool that will be broadcasted to the blockchain bundled.
4. Read - read a smartcontract information

# Requrements

1. GATEWAY_HOST environment variable
2. GATEWAY_PORT environment variable
3. GATEWAY_BROADCAST_HOST environment variable
4. GATEWAY_BROADCAST_PORT environment variable

# Usage

----------------------------------------------------------------
example of reading smartcontract data

	   import (
		"github.com/blocklords/gosds/sdk"
		"github.com/blocklords/gosds/topic"
	   )

	   func test() {
		// returns sdk.Reader
		reader := sdk.NewReader("address", "gateway repUrl")
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
	"errors"

	"github.com/blocklords/gosds/env"
	"github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/sdk/env"
	"github.com/blocklords/gosds/sdk/reader"
	"github.com/blocklords/gosds/sdk/subscriber"
	"github.com/blocklords/gosds/sdk/writer"
)

var Version string = "Seascape GoSDS version: 0.0.8"

// Returns a new reader.Reader.
//
// The repUrl is the link to the SDS Gateway.
// The address argument is the wallet address that is allowed to read.
func NewReader(address string) (*reader.Reader, error) {
	e, err := gatewayEnv(false)
	if err != nil {
		return nil, err
	}

	gatewaySocket := remote.TcpRequestSocketOrPanic(e)

	return reader.NewReader(gatewaySocket, address), nil
}

func NewWriter(address string) (*writer.Writer, error) {
	e, err := gatewayEnv(false)
	if err != nil {
		return nil, err
	}

	gatewaySocket := remote.TcpRequestSocketOrPanic(e)

	return writer.NewWriter(gatewaySocket, address), nil
}

// Returns a new subscriber
func NewSubscriber(address string) (*subscriber.Subscriber, error) {
	e, err := gatewayEnv(true)
	if err != nil {
		return nil, err
	}

	gatewaySocket := remote.TcpRequestSocketOrPanic(e)

	return subscriber.NewSubscriber(gatewaySocket, address), nil
}

// Returns the gateway environment variable
// If the broadcast argument set true, then Gateway will require the broadcast to be set as well.
func gatewayEnv(broadcast bool) (*env.Env, error) {
	e := env.Gateway()
	if !e.UrlExist() {
		return nil, errors.New("missing 'GATEWAY_HOST' and/or 'GATEWAY_PORT' environment variables")
	}

	if broadcast && !e.BroadcastExist() {
		return nil, errors.New("missing 'GATEWAY_BROADCAST_HOST' and/or 'GATEWAY_BROADCAST_PORT' environment variables")
	}

	return e, nil
}
