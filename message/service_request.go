package message

import (
	"encoding/json"
	"fmt"

	"github.com/blocklords/gosds/env"
)

// The SDS Service will accepts a request from another request
type ServiceRequest struct {
	Service    *env.Env               // The service parameters
	Command    string                 // Command type
	Parameters map[string]interface{} // Parameters of the request
}

// Convert ServiceRequest to JSON
func (request *ServiceRequest) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"service_name": request.Service.DomainName(),
		"public_key":   request.Service.PublicKey(),
		"command":      request.Command,
		"parameters":   request.Parameters,
	}
}

// ServiceRequest message as a  sequence of bytes
func (request *ServiceRequest) ToBytes() []byte {
	interfaces := request.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		fmt.Println("error while converting json into bytes", err)
		return []byte{}
	}

	return byt
}

// Convert ServiceRequest message to the string
func (request *ServiceRequest) ToString() string {
	return string(request.ToBytes())
}

// Parse the messages from zeromq into the ServiceRequest
func ParseServiceRequest(msgs []string) (ServiceRequest, error) {
	msg := ""
	for _, v := range msgs {
		msg += v
	}

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		return ServiceRequest{}, err
	}

	command, err := GetString(dat, "command")
	if err != nil {
		return ServiceRequest{}, err
	}
	parameters, err := GetMap(dat, "parameters")
	if err != nil {
		return ServiceRequest{}, err
	}

	public_key, err := GetString(dat, "public_key")
	if err != nil {
		return ServiceRequest{}, err
	}

	// The developers or smartcontract developer public keys are not in the environment variable
	// as a servie.
	service_env, err := env.GetByPublicKey(public_key)
	if err != nil {
		return ServiceRequest{}, err
	}

	request := ServiceRequest{
		Service:    service_env,
		Command:    command,
		Parameters: parameters,
	}

	return request, nil
}
