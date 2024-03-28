package CommonTest

import (
	"encoding/json"
	"fmt"
	. "github.com/bytedance/mockey"
	"log"
	"time"
)

type ParamStruct struct {
	Param  any `json:"param"`
	Config any `json:"config"`
}
type McokConfig struct {
	IsMock    bool   `json:"ismock"`
	MockValue string `json:"mockValue"`
}

func Exceldata2json(exceldata string, v any) {
	if exceldata == "nil" {
		v = nil
	} else if err := json.Unmarshal([]byte(exceldata), v); err != nil {
		log.Fatal(err)
	}
}
func StructInAnyData2json(s any, v any) {
	Exceldata2json(Struct2Str(s), v)
}

func Struct2Str(p any) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	s := string(bytes)
	return s
}
func TransfExcelData(exceldata string) any {
	if exceldata == "nil" {
		return nil
	} else if exceldata == "empty" {
		return ""
	} else if exceldata == "FALSE" || exceldata == "false" {
		return false
	} else if exceldata == "TRUE" || exceldata == "true" {
		return true
	}
	return exceldata

}

func Object2JSONStr(object any) string {
	jsonData, err := json.Marshal(object)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}

	// 输出JSON字符串
	return string(jsonData)
}

func MockBuildWithRetryWithSecond(builder *MockBuilder, maxRetries int, retryInterval int) *Mocker {
	retryCount := 0
	var mocker *Mocker
	for {
		// 增加重试计数器
		retryCount++

		// 使用匿名函数来捕获panic
		//func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
				// 如果达到最大重试次数，退出循环
				if retryCount >= maxRetries {
					fmt.Println("Max retries reached, exiting...")
					return
				}
				// 否则，等待后重试
				time.Sleep(time.Duration(retryInterval) * time.Second)
				return
			}
		}()

		// 调用可能引发panic的函数
		mocker = builder.Build()

		// 如果没有panic，跳出循环
		break
		//}()
	}
	log.Println("mocker build no panic!")
	// 如果没有panic，继续执行后续代码
	return mocker
}
