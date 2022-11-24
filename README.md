# SeascapeSDS & SeascapeSDK
> this is the golang package of the SeascapeSDS

This is an SDK to interact with SeascapeSDS. Also, this module contains the common data types of all SeascapeSDS.

SeascaeSDS incapsulates reading, writing data to the blockchain behind the API.

The exposed API is developer friendly.

Forget about remembering smartcontract address, smartcontract ABI or worrying about the safety of the private keys.

With this kind of solution, a backend developer could interact with the blockchain without knowing blockchain parts.

---

# Prerequirements

* a developer address derived from an assymetric key. An address whitelisted in SeascapeSDS (in the future it will be public).
* URL of the SDS Gateway (request-reply server).
* URL of the SDS Publisher (broadcaster of the transaction events).

*for example:*
```js
BACKEND_DEVELOPER_ACCOUNT=0x5bDed8f6BdAE766C361EDaE25c5DC966BCaF8f43
SDS_GATEWAY_URL=sds-gateway.seascape.network:3000
SDS_PUBLISHER_URL=sds-gateway.seascape.network:3001
```

---

# Understanding how SeascapeSDS does work

Using this SDK, a backend developer can interact with the smartcontract through SeascapeSDS. Two kind of operations users can do with the smartcontract:

* Read the data from smartcontract.
* Write the data to the smartcontract.

## Reading smartcontract
* Users can read the smartcontract's method.
* Users can listen for any transaction that occured with a smartcontract.

*For example, there is an NFT*
* *A backend developer can read the owner of the certain token*
* *A backend developer can subscribe for a minting of a new token*

## Writing smartcontract
* Users can send a transaction to the blockchain that invokes one of the smartcontract methods.
* Users can send the transaction to the pool. Once the pool is full or pull is timed out, the SeascapeSDS will send the batch as a single transaction.

*For example, there is an NFT*
* *A backend developer could transfer a token*
* *A backend developer could do a mass airdrop of the tokens to the customers*

## Setting up the smartcontract
But what kind of smartcontracts the backend developer could interact with? Who sets the smartcontracts? 
The smartcontracts that a backend developer could interact with are set up by a smartcontract developer.

Setting up of the smartcontracts are done using [SeascapeCLI](https://github.com/blocklords/seascape-cli). Check the documentation if you are a smartcontract developer.

*The reason why a smartcontract developer is handling the smartcontract setup is for two reasons.*
*Among all developers of the dapp, he is understandably the one who has the knowledge of how blockchain works. This knowledge is not required from the backend developer that uses SeascapeSDK*
*Secondly, the setup of the smartcontract is automated as much as possible. As a result, when a smartcontract deployed on the blockchain, **SeascapeCLI** will register the smartcontract on SDS automatically.*

## Topics
But how the backend developer do know what kind of smartcontract is setted up?

> In the future we will have a specific sector on the SDS webpage that keeps the list of the smartcontracts for the developer.

When a smartcontract is registered via **SeascapeCLI**, the smartcontract developer sets the *Topic* of the smartcontract.

Then, the smartcontract developer will share the topic with the backend developer.

Using the topic, backend developer could know what kind of smartcontract he can interact with.

### Topic structure
Topic has the following parameters:

* `organization` or `o` the name of the organization that handles the account whitelisting. *It's the name of the community, company that writes the dapp.*
* `project` or `p` the name of the project that holds the smartcontract. *Usually its the name of the dapp.*
* `network_id` or `n` the chain id where the smartcontract is deployed. *For example, it would be `"1"` for ethereum, `"imx"` if the smartcontract is deployed on the Immutable X.*
* `group` or `g` the name of the smartcontract group within the project. When a project has multiple smartcontracts, its better to group them. *If the smartcontract is the token, then the group name would be `"ERC20"`. If the smartcontract is an NFT, then the group name would be `"nft"`*
* `smartcontract` or `s` the name of the smartcontract. It should be identical to the filename of the smartcontract.
* `method` or `m` the name of the smartcontract method.
* `log` or `l` the name of the smartcontract event log.

*For example, to interact with the ScapeNFT from [Seascape Network](https://seascape.network) the following topics can be used*

```js
var1 = "o:seascape-network;p:core;n:1;g:nft;s:ScapeNFT;m:ownerOf"

var2 = "o:seascape-network;p:core;n:1;g:nft;s:ScapeNFT;l:Transfer"

```

the `var1` topic indicates the `ownerOf` method of the `ScapeNFT` smartcontract which is classified as `nft`, deployed on network `1`. The smartcontract is part of the `core` project of the `seascape-network` organization.

the `var2` topic indicates the `Transfer` event of the same `ScapeNFT` as in the `var1`.

### Topic String and Topic Object
The examples `var1` and `var2` showed above are examples of the *Topic String*. We represent the topic a string line.

However, topics can be also represented as the JSON object. The Topic Object structure is:

```js
{
    "o": "seascape-network", // organization
    "p": "core",             // project
    "n": "1",                // network id
    "g": "nft",              // group
    "s": "ScapeNFT",         // smartcontract name
    "m": "ownerOf",          // method
    "l": "Transfer"          // event log
}
```

Note that `"n"` value despite being a digit, its represented in the string format. **All values of the Topic are considered as a string**.

Note also, we represented both method and log in the topic. However, depending on the use case, you would use either method or log. Rarely you would need both of them.

### Filtering the Topic Path
The choice of the topics to represent the smartcontracts are for the reason.
It'ts a powerful tool to filter smartcontract.

What I mean is that you don't have to write topic string as we did in `var1` and `var2` above.

The `"o"`, `"m"` and other parts of the topic are called paths. The topic is power in the regard that you don't have to write full path. Some parts of the topic could be omitted.

*For example:*

```js
var3 = "o:seascape-network;p:core;g:nft;s:ScapeNFT"
```
Note that in the `var3` example we omitted `"n"`, `"m"`, and `"l"`. This means, this topic string will work with any method, log of the *ScapeNFT* smartcontract deployed on any network.
This is useful, if your smartcontract is multichain.

*Another example:*

```js
var4 = "g:nft;l:Transfer"
```

The topic string above means the `"Transfer"` log of any smartcontract from any project or organization that is classified as an `"nft"`.

*With this kind of topics for example, you can build data analytical tools or apply a machine learning on the blockchain data.*

As you see, the topic paths can be omitted. If the topic path is not written in the topic string, then it means, any path.

Question, what does the following valid topic string mean?

```js
var5 = "" // empty
```

### Topic Filter
However, SeascapeSDS provides another form of the Topic Strings and Topic Objects which is called Topic Filter.

The Topic Filter is almost identical to the Topic Strings, except that the path values are a list:

Here is the Topic Filter's JSON object:

```js
{
    "o": ["seascape-network"], // organization
    "p": ["core"],             // project
    "n": ["1"],                // network id
    "g": ["nft"],              // group
    "s": ["ScapeNFT"],         // smartcontract name
    "m": ["ownerOf"],          // method
    "l": ["Transfer"]          // event log
}
```

And here is the topic filter strings:

```js
var6 = "o:seascape-network;g:ERC20,nft;l:Mint,Burn,Transfer"
```

In the topic string, the list elements of each path are separated by comma.

In the example `var6` notice the `"g"` group and `"l"` paths.
The `"g"` path includes two elements: *ERC20* and *nft*.
The `"l"` path includes three elements: *Mint*, *Burn* and *Transfer*.

The example above means, that the topic filter does the operation with *Mint*, *Burn* and *Transfer* logs of any ERC20 or nft smartcontracts.

The Topic Filter allow to choose selected group of path.

### Message and Commands
Before we show the coding part, we need to talk two more things.

The SeascapeSDS composed of the various independent services all with `SDS` prefix.
Whether the SDS services interact to each other or external forces interact with it, user of the service interacts in the `Request-Reply` or `Pub-Sub` manners. 

The external force interacts with the SDS Gateway. And since this is one of the services, the manner of interaction applies to the gateway as well.

In a *Request-Reply* user sends a `command` including the parameters of the command. In exchange it gets the reply from the service.

Each of the services have a lit of commands that it supports. In the coding part of this documentation, I will show the list of all commands that is supported by `SDS Gateway`.

The `Pub-Sub` interact manner doesn't accept any commands. The service user is the `SUB`criber, while the SDS serive is `PUB`lisher.

In any case, those interaction manners always exchange the messages.

A message is an object with the parameters of the service or user needs.

There are three message types. `Request`, `Reply` and `Broadcast`.

Broadcast messages are the message that is send by `Pub` service to clients that are connected to the publisher. The `Broadcast` message is the encapsulation of the `Reply` message and the broadcasting `topic`. Don't worry about understanding broadcast messages. 
If you are using SeascapeSDK, you don't have to work with broadcast message or with broadcast topics. Consider them as the internal types. But its nice to know that they exist.

#### Request message
The request message is an JSON object with the two fields. Here is the reference how it looks:

```js
{
    "command": "command string",
    "parameters": {}
}
```

The `command` string accepts the name of the command. The `parameters` keeps the properties of the command. If the command doesn't have any property, then `parameters` field will keep an empty object.

The Request message is send to the `Request-Reply` service.

### Reply message
The `Request-Reply` service always return accepts `Request` messages. Then the service will return a `Reply` message back to the requester.

Here is the reference how it looks:

```js
{
    "status": "OK"   // status of reply
    "message": "",
    "parameters": {}
}
```

The `status` field value is either `"OK"` or `"fail"`.
If the status is a fail, then that means the requester didn't get what he was expected. If the status is OK, then the requester will get what he was expecting.

The `message` field contains the error message. If the status of the reply is a fail, then `message` field will contain a string explaining the what went wrong. If the status is OK, then the `message` field will be an empty string.

The `parameters` field contains the data that service replied back. If the status is OK, then `parameters` will contain the desired data that requester wants. If the reply is a failure, then the `parameters` will be an empty object.

---

# SDK interaction

Now, we come to the most interesting part. The following section lists all the SeascapeSDK functions + their examples to interact with the `SeascapeSDS`.

---

# SDS Gateway command list reference
Here is the list of the all commands that are supported by `SDS Gateway`.

You don't have to worry about them. Because `SeascapeSDK` or `SeascapeCLI` will create or read the messages for the user.

### command "smartcontract_read"
The command to be called from SeascapeSDK.
The command calls the smartcontract's read-only method.

Parameters of the `smartcontract_read` command:

`topic_string` string with the full path till the method. Full path means, that it should have any omitted string.
> Later we might support it if its needed.

`arguments` the function arguments to pass to the blockchain. The `arguments` is an object. If the smartcontract method doesn't have any method, then the arguments will be an empty object.
The argument names should be identical to the method argument names as defined in the source code of the smartcontract.

*For example the request*

```js
{
    "command": "smartcontract_read",
    "parameters": {
        "topic_string": "o:seascape-network;p:core;n:1;g:nft;s:ScapeNFT;m:ownerOf",
        "arguments": {
            "tokenId": 1
        }
    }
}
```

*A Reply from SDS Gateway:*
```js
{
    "status": "OK",
    "message": "",
    "parameters": {
        "result": {}
    }
}
```
If the smartcontract read returns a single data, then the `parameters.result` will be that single data. Otherwise it will contain an array of the smartcontract call response.

### command  "smartcontract_write"
The command to be called from SeascapeSDK.
The command calls the smartcontract's public methods that updates the smartcontract data.

Parameters of the `smartcontract_write` command:


`topic_string` string with the full path till the method. Full path means, that it should have any omitted string.
> Later we might support it if its needed.

`arguments` the function arguments to pass to the blockchain. The `arguments` is an object. If the smartcontract method doesn't have any method, then the arguments will be an empty object.
The argument names should be identical to the method argument names as defined in the source code of the smartcontract.

*For example the request*

```js
{
    "command": "smartcontract_write",
    "parameters": {
        "topic_string": "o:seascape-network;p:core;n:1;g:nft;s:ScapeNFT;m:transfer",
        "arguments": {
            "from": "",
            "to": "",
            "value": 1
        }
    }
}
```

*A Reply from SDS Gateway:*
```js
{
    "status": "OK",
    "message": "",
    "parameters": {
        "tx_id": "",
        "arguments": {}
    }
}
```

In a successful writing, the `SDS Gateway` will return the transaction id `tx_id`, the passed topic string and arguments from the requester.

> The returning of the transaction id, doesn't mean that transaction was confirmed. If it was confirmed, consider getting it using a `subscribe` command described below.

### command "subscribe"
The command to be called from SeascapeSDK.
The command indicates that the requester already connected to the `SDS Publisher` broadcaster and now ready to accept the messages that satisfies the topic filter.

Parameters of the `subscribe` command:

`topic_filter` a Topic Filter object, not a topic filter string. 

`subscriber` the backend developer's whitelisted account address.


*For example the request*

```js
{
    "command": "smartcontract_write",
    "parameters": {
        "topic_filter": {
            "from": "",
            "to": "",
            "value": 1
        },
        "subscriber": ""
    }
}
```

*A Reply from SDS Gateway:*
```js
{
    "status": "OK",
    "message": "",
    "parameters": {
        "block_timestamp": 1
    }
}
```

If the reply is successful, then the SDS Gateway will return a block timestamp. The SDS Publisher will send any event that matches to the `topic_filter` parameter of `Request` message that happened after the `block_timestamp`.

When a subscriber subscribes for the first time to the topic filter, the block timestamp will be the block timestamp of the earliest deployment of the smartcontract.
If the subscriber subscribes for the second time and on, then the `block_timestamp` will contain the last timestamp when a subscriber was online.

### command "heartbeat"
The command to be called from SeascapeSDK.
The command is called regularly after call of the `subscribe` command. The `heartbeat` command extends the connection to the `SDS Publisher`.

If the `heartbeat` command is not called within the FIVE seconds since the last `heartbeat` or `subscribe` command, then SeascapeSDS considers the Subscriber as offline and will stop broadcasting data.

The `heartbeat` command has a one parameter:

`subscriber` which should be identical as `subscriber` parameter of the `subscribe` command.

*For example the request*

```js
{
    "command": "heartbeat",
    "parameters": {
        "subscriber": ""
    }
}
```

*A Reply from SDS Gateway:*
```js
{
    "status": "OK",
    "message": "",
    "parameters": {
        "subscriber": ""
    }
}
```


### command "pool_add"
The command to be called from SeascapeSDK.
The command calls the smartcontract's public methods that updates the smartcontract data.

If the `smartcontract_write` is universal, then `pool_add` is for a method has an array argument.
The `pool_add` will store internally the elements of the array argument. Then after reaching to the pool limit or a timeout time, all the elements are packed as an array and send to the blockchain as a transaction.

> Requires smartcontract's bundler to be enabled by a smartcontract developer.

Parameters of the `pool_add` command:

`writer` string is the account of the backend developer that is whitelisted on SeascapeSDS.

`topic_string` string with the full path till the method. Full path means, that it should have any omitted string.
> Later we might support it if its needed.

`arguments` the function arguments to pass to the blockchain. The `arguments` is an object. If the smartcontract method doesn't have any method, then the arguments will be an empty object.
The argument names should be identical to the method argument names as defined in the source code of the smartcontract.

*For example the request*

```js
{
    "command": "pool_add",
    "parameters": {
        "writer": "",
        "topic_string": "o:seascape-network;p:core;n:1;g:nft;s:ScapeNFT;m:transfer",
        "arguments": {
            "from": "",
            "to": "",
            "value": 1
        }
    }
}
```

*A Reply from SDS Gateway:*
```js
{
    "status": "OK",
    "message": "",
    "parameters": {
        "topic_string":        "",
		"pool":                [{}],
		"pool_length":         3,
		"pool_limit":          5,
		"timer_left":          5,
		"next_execution_time": "123",
    }
}
```

In a successful writing, the `SDS Gateway` will return the pool information along with the topic string with a full path.

`pool` parameter of the reply is an array of function arguments.
`pool_length` is the size of the pool.
`pool_limit` is the limit of the pool array. Upon reaching to the pool limit, the pool is cleared, and all data in the pool is send to the blockchain.
`timer_left` an integer indicating the seconds till the automatic execution of the pool, even when a pool is not full.
`next_execution_time` is the unix timestamp in seconds when the automatic pool execution will occur.

> To listen the success of the pool execution, use the `subscribe` command described above.
