// The network package is used to get the blockchain network information.
package network

import (
	"encoding/json"
	"os"
	"strings"
)

type Network struct {
	Id       string
	Provider string
	Flag     int8 // With VM or Without VM
}

// any blockchain
const ALL int8 = 0

// filter only blockchain networks with Virtual Machine.
// "imx" will not be here.
const WITH_VM int8 = 1

const WITHOUT_VM int8 = -1

// Returns list of the supported networks by this SDS Service
func GetSupportedNetworks(flag int8) map[string]string {
	env := os.Getenv("SUPPORTED_NETWORKS")
	if len(env) == 0 {
		panic("the environment variable 'SUPPORTED_NETWORKS' is not provided")
	}

	var supportedNetworks map[string]string

	parse_err := json.Unmarshal([]byte(env), &supportedNetworks)
	if parse_err != nil {
		panic("the environment variable 'SUPPORTED_NETWORKS' is not a valid JSON")
	}

	if flag == ALL {
		return supportedNetworks
	}

	// without VM
	imx := "imx"

	for networkId := range supportedNetworks {
		if strings.ToLower(networkId) == imx {
			if flag == WITH_VM {
				delete(supportedNetworks, networkId)
			}
		} else if flag == WITHOUT_VM {
			delete(supportedNetworks, networkId)
		}
	}

	return supportedNetworks
}

// Returns list of support network IDs
func GetNetworkIds(flag int8) []string {
	supportedNetworks := GetSupportedNetworks(flag)

	ids := make([]string, 0)

	if len(supportedNetworks) == 0 {
		return ids
	}

	for networkId := range supportedNetworks {
		ids = append(ids, networkId)
	}
	return ids
}

// Whether the given network id is supported by this SDS Service.
func IsSupportedNetwork(networkId string, flag int8) bool {
	supportedNetworks := GetSupportedNetworks(flag)
	if len(supportedNetworks) == 0 {
		return false
	}

	_, ok := supportedNetworks[networkId]
	return ok
}

// Returns the Blockchain Network access provider
func GetProvider(networkId string, flag int8) string {
	if !IsSupportedNetwork(networkId, flag) {
		return ""
	}
	supportedNetworks := GetSupportedNetworks(flag)

	return supportedNetworks[networkId]
}
