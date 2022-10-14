package network

import (
	"encoding/json"
	"os"
)

func GetSupportedNetworks() map[string]string {
	env := os.Getenv("SUPPORTED_NETWORKS")
	if len(env) == 0 {
		panic("the environment variable 'SUPPORTED_NETWORKS' is not provided")
	}

	var supportedNetworks map[string]string

	parse_err := json.Unmarshal([]byte(env), &supportedNetworks)
	if parse_err != nil {
		panic("the environment variable 'SUPPORTED_NETWORKS' is not a valid JSON")
	}

	return supportedNetworks
}

func GetNetworkIds() []string {
	supportedNetworks := GetSupportedNetworks()

	ids := make([]string, 0)

	if len(supportedNetworks) == 0 {
		return ids
	}

	for networkId := range supportedNetworks {
		ids = append(ids, networkId)
	}
	return ids
}

func IsSupportedNetwork(networkId string) bool {
	supportedNetworks := GetSupportedNetworks()
	if len(supportedNetworks) == 0 {
		return false
	}

	_, ok := supportedNetworks[networkId]
	if !ok {
		return false
	}

	return true
}

func GetProvider(networkId string) string {
	if !IsSupportedNetwork(networkId) {
		return ""
	}
	supportedNetworks := GetSupportedNetworks()

	return supportedNetworks[networkId]
}
