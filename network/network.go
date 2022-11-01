package network

import (
	"encoding/json"
	"os"
	"strings"
)

// all networks including IMX
const ALL = true

// exclude IMX from the supported networks
const NOT_ALL = false

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

func IsSupportedNetwork(networkId string) bool {
	supportedNetworks := GetSupportedNetworks(ALL)
	if len(supportedNetworks) == 0 {
		return false
	}

	_, ok := supportedNetworks[networkId]
	return ok
}

func GetProvider(networkId string) string {
	if !IsSupportedNetwork(networkId) {
		return ""
	}
	supportedNetworks := GetSupportedNetworks(ALL)

	return supportedNetworks[networkId]
}
