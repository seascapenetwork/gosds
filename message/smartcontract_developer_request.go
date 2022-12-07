package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/blocklords/gosds/account"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// The SDS Service will accepts the SmartcontractDeveloperRequest message.
type SmartcontractDeveloperRequest struct {
	Address        string                 // The whitelisted address of the user
	NonceTimestamp uint64                 // Nonce as a unix timestamp in seconds
	Signature      string                 // Command, nonce, address and parameters signed together
	Command        string                 // Command type
	Parameters     map[string]interface{} // Parameters of the request
}

// Convert SmartcontractDeveloperRequest to JSON
func (request *SmartcontractDeveloperRequest) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"address":         request.Address,
		"nonce_timestamp": request.NonceTimestamp,
		"signature":       request.Signature,
		"command":         request.Command,
		"parameters":      request.Parameters,
	}
}

// SmartcontractDeveloperRequest message as a  sequence of bytes
func (request *SmartcontractDeveloperRequest) ToBytes() []byte {
	interfaces := request.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		fmt.Println("error while converting json into bytes", err)
		return []byte{}
	}

	return byt
}

// Convert SmartcontractDeveloperRequest message to the string
func (request *SmartcontractDeveloperRequest) ToString() string {
	return string(request.ToBytes())
}

// Gets the message without a prefix.
// The message is a JSON represantion of the Request but without "signature" parameter.
// Converted into the hash using Keccak32.
//
// The request parameters are oredered in an alphanumerical order.
func (request *SmartcontractDeveloperRequest) message_hash() []byte {
	json_object := request.ToJSON()
	delete(json_object, "signature")
	bytes, err := json.Marshal(json_object)
	if err != nil {
		fmt.Println("error while converting json into bytes", err)
		return []byte{}
	}

	hash := crypto.Keccak256Hash(bytes)

	return hash.Bytes()
}

// Gets the digested message with a prefix
// For ethereum the prefix is "\x19Ethereum Signed Message:\n"
func (request *SmartcontractDeveloperRequest) digested_message() []byte {
	message_hash := request.message_hash()
	prefix := []byte("\x19Ethereum Signed Message:\n32")
	digested_hash := crypto.Keccak256Hash(append(prefix, message_hash...))
	return digested_hash.Bytes()
}

// Get the account who did the request.
// Account is verified first using the signature parameter of the request.
// If the signature is not a valid, then returns an error.
//
// For now it supports ECDSA addresses only. Therefore verification automatically assumes that address
// is for the ethereum network.
func (request *SmartcontractDeveloperRequest) GetAccount() (*account.SmartcontractDeveloper, error) {
	// without 0x prefix
	signature, err := hexutil.Decode(request.Signature)
	if err != nil {
		return nil, err
	}
	digested_hash := request.digested_message()

	if len(signature) != 65 {
		return nil, errors.New("the ECDSA signature length is invalid. It should be 64 bytes long. Signature length: ")
	}
	if signature[64] != 27 && signature[64] != 28 {
		return nil, errors.New("invalid Ethereum signature (V is not 27 or 28)")
	}
	signature[64] -= 27 // Transform yellow paper V from 27/28 to 0/1

	ecdsa_public_key, err := crypto.SigToPub(digested_hash, signature)
	if err != nil {
		return nil, err
	}

	address := crypto.PubkeyToAddress(*ecdsa_public_key).Hex()
	if !strings.EqualFold(address, request.Address) {
		return nil, errors.New("the request 'address' parameter mismatches to the account derived from signature. Account derived from the signature: " + address + "...")
	}

	return account.NewEcdsaPublicKey(ecdsa_public_key), nil
}

// Parse the messages from zeromq into the SmartcontractDeveloperRequest
func ParseSmartcontractDeveloperRequest(msgs []string) (SmartcontractDeveloperRequest, error) {
	msg := ""
	for _, v := range msgs {
		msg += v
	}

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		return SmartcontractDeveloperRequest{}, err
	}

	command, err := GetString(dat, "command")
	if err != nil {
		return SmartcontractDeveloperRequest{}, err
	}
	parameters, err := GetMap(dat, "parameters")
	if err != nil {
		return SmartcontractDeveloperRequest{}, err
	}

	address, err := GetString(dat, "address")
	if err != nil {
		return SmartcontractDeveloperRequest{}, err
	}

	nonce_timestamp, err := GetUint64(dat, "nonce_timestamp")
	if err != nil {
		return SmartcontractDeveloperRequest{}, err
	}

	signature, err := GetString(dat, "signature")
	if err != nil {
		return SmartcontractDeveloperRequest{}, err
	}

	request := SmartcontractDeveloperRequest{
		Address:        address,
		NonceTimestamp: nonce_timestamp,
		Signature:      signature,
		Command:        command,
		Parameters:     parameters,
	}

	return request, nil
}
