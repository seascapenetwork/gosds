/* The ABI file handler */
package categorizer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/static"

	"github.com/ethereum/go-ethereum/accounts/abi"
	zmq "github.com/pebbe/zmq4"
)

type Abi struct {
	staticAbi *static.Abi
	i         abi.ABI // interface
}

func (a *Abi) StringReader() *strings.Reader {
	return strings.NewReader(string(a.staticAbi.Bytes))
}

func (a *Abi) GetMethod(method string) (*abi.Method, error) {
	for _, m := range a.i.Methods {
		if m.Name == method {
			return &m, nil
		}
	}

	return nil, errors.New("no method found")
}

func (a *Abi) ParseTxInput(data string) (string, map[string]interface{}, error) {
	inputs := map[string]interface{}{}

	offset := 0
	if len(data) > 2 && data[:2] == "0x" || data[:2] == "0X" {
		offset += 2
	}

	fmt.Println("Parse tx input: ", data)

	// decode txInput method signature
	decodedSig, err := hex.DecodeString(data[offset : 8+offset])
	if err != nil {
		return "", inputs, err
	}
	fmt.Println("Function sig: ", decodedSig)
	fmt.Println(a.i.Methods)

	// recover Method from signature and ABI
	method, err := a.i.MethodById(decodedSig)
	if err != nil {
		return "", inputs, err
	}

	// decode txInput Payload
	decodedData, err := hex.DecodeString(data[8+offset:])
	if err != nil {
		return method.Name, inputs, err
	}

	err = method.Inputs.UnpackIntoMap(inputs, decodedData)
	if err != nil {
		fmt.Println("Failed to unpack tx data")
		return method.Name, inputs, err
	}

	return method.Name, inputs, nil
}

func Build(abiHash string, body interface{}) Abi {
	staticAbi := static.BuildAbi(body)
	staticAbi.AbiHash = abiHash

	abiObj := Abi{staticAbi: staticAbi}

	abiReader := abiObj.StringReader()
	i, abiErr := abi.JSON(abiReader)
	if abiErr != nil {
		fmt.Println("Failed to parse ABI", abiErr)
	}
	abiObj.i = i

	println("The abi built and return it back")
	return abiObj
}

func AbiGet(socket zmq.Socket, abiHash string) (Abi, error) {
	// Send hello.
	abiGetRequest := message.Request{
		Command: "abi_get",
		Param: map[string]interface{}{
			"abi_hash": abiHash,
		},
	}
	fmt.Println("Sending message to STATIC server to get abi. The mesage sent to server")
	fmt.Println(abiGetRequest.ToString())
	if _, err := socket.SendMessage(abiGetRequest.ToString()); err != nil {
		fmt.Println("Failed to send a command for abi getting from static controller", err.Error())
		return Abi{}, fmt.Errorf("sending: %w", err)
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller")
		return Abi{}, fmt.Errorf("receiving: %w", err)
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse abi reply")
		return Abi{}, fmt.Errorf("spaghetti block invalid Reply: %w", err)
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure")
		return Abi{}, fmt.Errorf("spaghetti block reply status is not ok: %s", reply.Message)
	}

	fmt.Println("Abi build and returned")
	return Build(abiHash, reply.Params["abi"]), nil
}
