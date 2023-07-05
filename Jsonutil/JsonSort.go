package Jsonutil

import (
	"github.com/json-iterator/go"
	"sort"
)

func SortJSON(input string) (string, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var jmap interface{}
	err := json.UnmarshalFromString(input, &jmap)
	if err != nil {
		return "", err
	}

	sorted := sortMapKeys(jmap)

	b, err := json.Marshal(sorted)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func sortMapKeys(m interface{}) interface{} {
	switch m := m.(type) {
	case map[string]interface{}:
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		sortedMap := make(map[string]interface{})
		for _, k := range keys {
			sortedMap[k] = sortMapKeys(m[k])
		}
		return sortedMap
	case []interface{}:
		for i, v := range m {
			m[i] = sortMapKeys(v)
		}
		return m
	default:
		return m
	}
}
