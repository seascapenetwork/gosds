package network

import (
	"encoding/json"
	"os"
	"strings"
)

// any blockchain
const ALL = true

// filter only blockchain networks with Virtual Machine.
// "imx" will not be here.
const WITH_VM = false

func GetSupportedNetworks(all bool) map[string]string {
	env := os.Getenv("SUPPORTED_NETWORKS")
	if len(env) == 0 {
		panic("the environment variable 'SUPPORTED_NETWORKS' is not provided")
	}

	var supportedNetworks map[string]string

	parse_err := json.Unmarshal([]byte(env), &supportedNetworks)
	if parse_err != nil {
		panic("the environment variable 'SUPPORTED_NETWORKS' is not a valid JSON")
	}

	if all {
		return supportedNetworks
	}

	imx := "imx"
	for networkId := range supportedNetworks {
		if strings.ToLower(networkId) == imx {
			delete(supportedNetworks, networkId)
		}
	}

	return supportedNetworks
}

func GetNetworkIds(all bool) []string {
	supportedNetworks := GetSupportedNetworks(all)

	ids := make([]string, 0)

	if len(supportedNetworks) == 0 {
		return ids
	}

	for networkId := range supportedNetworks {
		ids = append(ids, networkId)
	}
	return ids
}

func IsSupportedNetwork(networkId string, all bool) bool {
	supportedNetworks := GetSupportedNetworks(all)
	if len(supportedNetworks) == 0 {
		return false
	}

	_, ok := supportedNetworks[networkId]
	return ok
}

func GetProvider(networkId string, all bool) string {
	if !IsSupportedNetwork(networkId, all) {
		return ""
	}
	supportedNetworks := GetSupportedNetworks(all)

	return supportedNetworks[networkId]
}
