package message

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Returns the parameter as an uint64
func GetUint64(parameters map[string]interface{}, name string) (uint64, error) {
	raw, exists := parameters[name]
	if !exists {
		return 0, errors.New("missing '" + name + "' parameter in the Request")
	}

	value, ok := raw.(json.Number)
	if !ok {
		return 0, errors.New("parameter '" + name + "' expected to be as a number")
	}

	number, err := strconv.ParseUint(string(value), 10, 64)

	return number, err
}

func GetFloat64(parameters map[string]interface{}, name string) (float64, error) {
	raw, exists := parameters[name]
	if !exists {
		return 0, errors.New("missing '" + name + "' parameter in the Request")
	}
	value, err := raw.(json.Number).Float64()

	return value, err
}

// Returns the paramater as a string
func GetString(parameters map[string]interface{}, name string) (string, error) {
	raw, exists := parameters[name]
	if !exists {
		return "", errors.New("missing '" + name + "' parameter in the Request")
	}
	value, ok := raw.(string)
	if !ok {
		return "", errors.New("expected string type for '" + name + "' parameter")
	}

	return value, nil
}

// Returns list of strings
func GetStringList(parameters map[string]interface{}, name string) ([]string, error) {
	raw, exists := parameters[name]
	if !exists {
		return nil, errors.New("missing '" + name + "' parameter in the Request")
	}

	values, ok := raw.([]interface{})
	if !ok {
		ready_list, ok := raw.([]string)
		if !ok {
			return nil, errors.New("expected list type for '" + name + "' parameter")
		} else {
			return ready_list, nil
		}
	}

	list := make([]string, len(values))
	for i, raw_value := range values {
		v, ok := raw_value.(string)
		if !ok {
			return nil, errors.New("one of the elements in the parameter is not a string")
		}

		list[i] = v
	}

	return list, nil
}

// Returns the parameter as a slice of map:
//
// []map[string]interface{}
func GetMapList(parameters map[string]interface{}, name string) ([]map[string]interface{}, error) {
	raw, exists := parameters[name]
	if !exists {
		return nil, errors.New("missing '" + name + "' parameter in the Request")
	}
	values, ok := raw.([]interface{})
	if !ok {
		ready_list, ok := raw.([]map[string]interface{})
		if !ok {
			return nil, errors.New("expected list type for '" + name + "' parameter")
		} else {
			return ready_list, nil
		}
	}

	list := make([]map[string]interface{}, len(values))
	for i, raw_value := range values {
		v, ok := raw_value.(map[string]interface{})
		if !ok {
			return nil, errors.New("one of the elements in the parameter is not a map")
		}

		list[i] = v
	}

	return list, nil
}

// Returns the parameter as a map:
//
// map[string]interface{}
func GetMap(parameters map[string]interface{}, name string) (map[string]interface{}, error) {
	raw, exists := parameters[name]
	if !exists {
		return nil, errors.New("missing '" + name + "' parameter in the Request")
	}
	value, ok := raw.(map[string]interface{})
	if !ok {
		return nil, errors.New("expected map type for '" + name + "' parameter")
	}

	return value, nil
}
