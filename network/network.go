package network

import (
	"encoding/json"
	"os"
	"strings"
)

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

func GetNetworkIds(fullSupport bool) []string {
	supportedNetworks := GetSupportedNetworks(fullSupport)

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
	supportedNetworks := GetSupportedNetworks(true)
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
	supportedNetworks := GetSupportedNetworks(true)

	return supportedNetworks[networkId]
}
