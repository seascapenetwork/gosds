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

// Converts the data structs to the JSON objects (represented as a golang map) list.
func ToMapList[V SDS_Data](list []V) []map[string]interface{} {

	map_list := make([]map[string]interface{}, len(list))
	for i, element := range list {
		map_list[i] = element.ToJSON()
	}

	return map_list
}
