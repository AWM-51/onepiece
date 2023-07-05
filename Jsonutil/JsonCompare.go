package Jsonutil

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// compareJSON 比较两个 JSON 字符串是否相同，忽略 ignoreFields 中的字段
func CompareJSONWithignore(json1, json2 string, ignoreFields []string) (bool, map[string]interface{}) {
	var obj1, obj2 interface{}
	err := json.Unmarshal([]byte(json1), &obj1)
	if err != nil {
		return false, nil
	}
	err = json.Unmarshal([]byte(json2), &obj2)
	if err != nil {
		return false, nil
	}

	return CompareObjects(obj1, obj2, "", ignoreFields)
}

// compareObjects 比较两个对象是否相同，忽略 ignoreFields 中的字段
func CompareObjects(obj1, obj2 interface{}, currentPath string, ignoreFields []string) (bool, map[string]interface{}) {
	value1 := reflect.ValueOf(obj1)
	value2 := reflect.ValueOf(obj2)

	result := make(map[string]interface{})
	if value1.Kind() != value2.Kind() {
		result["{The type is not equal}"] = fmt.Sprintf("the one kind is %d,the other kind is %d", value1.Kind(), value2.Kind())
		return false, result
	}

	switch value1.Kind() {
	case reflect.Map:
		return CompareMaps(value1, value2, currentPath, ignoreFields)
	case reflect.Slice:
		return CompareSlices(value1, value2, currentPath, ignoreFields)
	default:
		return CompareValues(value1, value2, ignoreFields)
	}
}

// compareMaps 比较两个 Map 是否相同，忽略 ignoreFields 中的字段
func CompareMaps(map1, map2 reflect.Value, currentPath string, ignoreFields []string) (bool, map[string]interface{}) {
	result := make(map[string]interface{})
	if map1.Len() != map2.Len() {
		result["{The JSON lengths are inconsistent}"] = fmt.Sprintf("the one is %d,the other is %d", map1.Len(), map2.Len())
		return false, result
	}

	for _, key := range map1.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())

		// 拼接当前字段的路径到父节点的路径上，形成完整的字段路径
		nextPath := keyStr
		if currentPath != "" {
			nextPath = fmt.Sprintf("%s.%s", currentPath, keyStr)
		}

		// 判断当前字段路径是否
		if ContainsString(ignoreFields, nextPath) {
			continue
		}
		// 比较两个字段的值
		value1 := map1.MapIndex(key)
		value2 := map2.MapIndex(key)
		equal, diff := CompareObjects(value1.Interface(), value2.Interface(), nextPath, ignoreFields)
		if !equal {
			result[keyStr] = diff
		}
	}

	if len(result) == 0 {
		return true, nil
	}

	return false, result
}

// compareSlices 比较两个 Slice 是否相同，忽略 ignoreFields 中的字段
func CompareSlices(slice1, slice2 reflect.Value, currentPath string, ignoreFields []string) (bool, map[string]interface{}) {
	result := make(map[string]interface{})
	if slice1.Len() != slice2.Len() {
		result["{The slices lengths are inconsistent}"] = nil
		return false, nil
	}

	for i := 0; i < slice1.Len(); i++ {
		// 拼接当前字段的路径到父节点的路径上，形成完整的字段路径
		nextPath := fmt.Sprintf("%s.[%d]", currentPath, i)

		// 判断当前字段路径是否需要忽略
		if ContainsString(ignoreFields, nextPath) {
			continue
		}

		// 比较两个字段的值
		value1 := slice1.Index(i)
		value2 := slice2.Index(i)
		equal, diff := CompareObjects(value1.Interface(), value2.Interface(), nextPath, ignoreFields)
		if !equal {
			result[fmt.Sprintf("[%d]", i)] = diff
		}
	}

	if len(result) == 0 {
		return true, nil
	}

	return false, result
}

//compareValues 比较两个值是否相同，忽略 ignoreFields 中的字段
func CompareValues(value1, value2 reflect.Value, ignoreFields []string) (bool, map[string]interface{}) {
	// 获取当前字段的路径
	currentPath := ""

	// 判断当前字段路径是否需要忽略
	if ContainsString(ignoreFields, currentPath) {
		return true, nil
	}

	result := make(map[string]interface{})
	if !value1.IsValid() && !value2.IsValid() {
		return true, nil
	}
	if value1.IsValid() && !value2.IsValid() {
		result["{Inconsistent values}"] = nil
		return false, result
	}

	if !value1.IsValid() && value2.IsValid() {
		result["{Inconsistent values}"] = nil
		return false, result
	}

	if reflect.DeepEqual(value1.Interface(), value2.Interface()) {
		return true, nil
	}

	result["{Inconsistent values}"] = value1.Interface()

	return false, result
}

// containsString 判断字符串切片中是否包含指定字符串
func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// containsJSON 判断一个JSON对象是否包含另一个JSON对象
func ContainsJSON(expected, actual interface{}) bool {
	if reflect.DeepEqual(expected, actual) {
		return true
	}

	switch exp := expected.(type) {
	case map[string]interface{}:
		actMap, ok := actual.(map[string]interface{})
		if !ok {
			return false
		}
		for key, expVal := range exp {
			actVal, ok := actMap[key]
			if !ok {
				return false
			}
			if !ContainsJSON(expVal, actVal) {
				return false
			}
		}
		return true
	case []interface{}:
		actSlice, ok := actual.([]interface{})
		if !ok {
			return false
		}
		for i, expVal := range exp {
			if i >= len(actSlice) {
				return false
			}
			actVal := actSlice[i]
			if !ContainsJSON(expVal, actVal) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func GetNestedKeys(m map[string]interface{}) []string {
	//var res = make(map[string]interface{})
	var keys []string
	for key, value := range m {
		switch value.(type) {
		case map[string]interface{}:
			subKeys := GetNestedKeys(value.(map[string]interface{}))
			for _, subKey := range subKeys {
				keys = append(keys, fmt.Sprintf("%s.%s", key, subKey))
			}
		default:
			keys = append(keys, key)
		}
	}
	return keys
}

func GetKeys(m map[string]interface{}) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[j] = k
		j++
	}
	return keys
}
