package utils

import (
	"encoding/json"
	"reflect"
	"sort"

	"gorm.io/datatypes"
)

// ParseJSONInterface разбирает JSON в interface{}:
//   - map[string]interface{} для объектов,
//   - []interface{} для массивов,
//   - string, float64 и т.д. для примитивов.
func ParseJSONInterface(data datatypes.JSON) interface{} {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil
	}
	return v
}

func ParseJSONToMap(data []byte) map[string]interface{} {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}

	return result
}

func ParseMapToJSON(value map[string]interface{}) datatypes.JSON {
	bytes, err := json.Marshal(value)
	if err != nil {
		return nil
	}

	return datatypes.JSON(bytes)
}

func ParseToJSON(value interface{}) datatypes.JSON {
	bytes, err := json.Marshal(value)
	if err != nil {
		return nil
	}

	return datatypes.JSON(bytes)
}

func ToStringSlice(input interface{}) ([]string, bool) {
	arr, ok := input.([]interface{})
	if !ok {
		return nil, false
	}
	result := make([]string, len(arr))
	for i, val := range arr {
		str, ok := val.(string)
		if !ok {
			return nil, false
		}
		result[i] = str
	}
	return result, true
}

func EqualStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)
	return reflect.DeepEqual(a, b)
}
