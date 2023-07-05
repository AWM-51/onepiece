package Assertutil

import (
	"encoding/json"
	"fmt"
	jsonschema "github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	jsonHandle "onepiece/Jsonutil"
	"reflect"
	"strings"
	"testing"
)

// Assert 断言结构体
type Assert struct {
	t *testing.T
}

// NewAssert 创建一个新的断言实例
func NewAssert(t *testing.T) *Assert {
	return &Assert{t: t}
}

// Equal 判断两个值是否相等
func (a *Assert) Equal(expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		a.t.Errorf("Expected: %v, but got: %v", expected, actual)
	}
}

// NotEqual 判断两个值是否不相等
func (a *Assert) NotEqual(expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		a.t.Errorf("Expected: %v to not equal: %v", expected, actual)
	}
}

// GreaterThan 判断一个值是否大于另一个值
func (a *Assert) GreaterThan(expected, actual interface{}) {
	res, error := CompareValues(expected, actual)
	if error != nil {
		a.t.Errorf(error.Error())
	}
	if error == nil && !(res == 1) {
		a.t.Errorf("Expected not GreaterThan got;Expected : %#v, but got: %#v", expected, actual)
	}
}

// LessThan 判断一个值是否小于另一个值
func (a *Assert) LessThan(expected, actual interface{}) {
	res, error := CompareValues(expected, actual)
	if error != nil {
		a.t.Errorf(error.Error())
	}
	if error == nil && !(res == -1) {
		a.t.Errorf("Expected not LessThan gotExpected: %#v, but got: %#v", expected, actual)
	}
}

// Contains 判断一个值是否包含在另一个值中
func (a *Assert) Contains(expected, actual interface{}) {
	act := reflect.ValueOf(actual)
	exp := reflect.ValueOf(expected)

	if act.Kind() == reflect.String && exp.Kind() == reflect.String {
		if !strings.Contains(act.String(), exp.String()) {
			a.t.Errorf("Expected: %v to be contained in %v", expected, actual)
		}
	} else if act.Kind() == reflect.Slice || act.Kind() == reflect.Array {
		for i := 0; i < act.Len(); i++ {
			if reflect.DeepEqual(act.Index(i).Interface(), expected) {
				return
			}
		}
		a.t.Errorf("Expected: %v to be contained in %v", expected, actual)
	} else {
		a.t.Errorf("Expected: %v, but got: %v", expected, actual)
	}
}

// JSONEqual 判断两个JSON字符串
func (a *Assert) JSONEqual(expected, actual string) {
	var expJSON interface{}
	var actJSON interface{}

	err := json.Unmarshal([]byte(expected), &expJSON)
	if err != nil {
		a.t.Errorf("Failed to unmarshal expected JSON: %v", err)
	}

	err = json.Unmarshal([]byte(actual), &actJSON)
	if err != nil {
		a.t.Errorf("Failed to unmarshal actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expJSON, actJSON) {
		a.t.Errorf("Expected: %s, but got: %s", expected, actual)
	}
}

// JSONContains 判断一个JSON字符串是否包含另一个JSON字符串
// 后者包含前者 为正确
func (a *Assert) JSONContains(expected, actual string) {
	var expJSON interface{}
	var actJSON interface{}

	err := json.Unmarshal([]byte(expected), &expJSON)
	if err != nil {
		a.t.Errorf("Failed to unmarshal expected JSON: %v", err)
	}

	err = json.Unmarshal([]byte(actual), &actJSON)
	if err != nil {
		a.t.Errorf("Failed to unmarshal actual JSON: %v", err)
	}

	if !jsonHandle.ContainsJSON(expJSON, actJSON) {
		a.t.Errorf("Expected: %s to be contained in %s", expected, actual)
	}
}

// jsonschema 断言
// 输入的参数是文件地址和目标json
func (a *Assert) JsonschemaAssert(schemaFile string, postResult string) {
	b1, err := ioutil.ReadFile(schemaFile) //读取文件
	if err != nil {
		panic(err.Error())
	}

	a.JsonschemaStrAssert(string(b1), postResult)
}

// 输入的参数为 jsonschema String值 和目标json
func (a *Assert) JsonschemaStrAssert(schemaStr string, postResult string) {
	schemaLoader := jsonschema.NewStringLoader(schemaStr)    // jsonschema格式
	documentLoader := jsonschema.NewStringLoader(postResult) // 待校验的json数据

	result, err := jsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		//测试成功可以写点东西，但是我不想写
		//a.t.Log("The document is valid")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			a.t.Errorf("jsonschema assert result ：%s", desc)
		}
	}
}

// ObjectEqual 判断两个对象是否相等
func (a *Assert) ObjectEqual(expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		a.t.Errorf("Expected: %v, but got: %v", expected, actual)
	}
}

//JsonAssertWithignore 忽略判断两个json某些字段后再校验
func (a *Assert) JsonAssertWithignore(jsonStr1, jsonStr2 string, ignoreFields []string) {
	equal, diffFields := jsonHandle.CompareJSONWithignore(jsonStr1, jsonStr2, ignoreFields)
	if equal {
		a.t.Errorf("The two JSON strings are equal.")
	} else {
		a.t.Errorf("The two JSON strings are not equal.Different fields and results:%v", jsonHandle.GetNestedKeys(diffFields))
	}
}

func CompareValues(val1, val2 interface{}) (int, error) {
	// 获取 val1 和 val2 的 reflect.Value
	v1 := reflect.ValueOf(val1)
	v2 := reflect.ValueOf(val2)

	// 判断 val1 和 val2 是否为可比较的类型
	if !v1.Type().Comparable() || !v2.Type().Comparable() {
		return -9, fmt.Errorf("unsupported type")
	}

	// 比较 val1 和 val2
	switch v1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 处理 int 类型比较
		if v1.Int() > v2.Int() {
			return 1, nil
		} else if v1.Int() < v2.Int() {
			return -1, nil
		} else {
			return 0, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 处理 uint 类型比较
		if v1.Uint() > v2.Uint() {
			return 1, nil
		} else if v1.Uint() < v2.Uint() {
			return -1, nil
		} else {
			return 0, nil
		}
	case reflect.Float32, reflect.Float64:
		// 处理 float 类型比较
		if v1.Float() > v2.Float() {
			return 1, nil
		} else if v1.Float() < v2.Float() {
			return -1, nil
		} else {
			return 0, nil
		}
	case reflect.String:
		// 处理 float 类型比较
		if v1.String() > v2.String() {
			return 1, nil
		} else if v1.String() < v2.String() {
			return -1, nil
		} else {
			return 0, nil
		}
	default:
		return -8, fmt.Errorf("unsupported type")
	}
}
