// Generic Type package handles the common functions of multiple SDS Data structures.
package generic_type

import (
	"github.com/blocklords/gosds/categorizer"
	"github.com/blocklords/gosds/spaghetti"
	"github.com/blocklords/gosds/static"
)

type SDS_Data interface {
	*categorizer.Log | *categorizer.Smartcontract | *categorizer.Transaction |
		*spaghetti.Log | *spaghetti.Transaction | *static.Abi |
		*static.Configuration | *static.Smartcontract

	ToJSON() map[string]interface{}
}

type SDS_String_Data interface {
	*static.SmartcontractKey
}

// Converts the data structs to the JSON objects (represented as a golang map) list.
// []map[string]interface{}
func ToMapList[V SDS_Data](list []V) []map[string]interface{} {
	map_list := make([]map[string]interface{}, len(list))
	for i, element := range list {
		map_list[i] = element.ToJSON()
	}

	return map_list
}

// Converts the data structs to the list of strings.
// []string
func ToStringList[V SDS_String_Data](list []V) []string {
	string_list := make([]string, len(list))
	for i, element := range list {
		string_list[i] = string(*element)
	}

	return string_list
}
