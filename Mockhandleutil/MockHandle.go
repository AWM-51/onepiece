package Mockhandleutil

import (
	"fmt"
	. "github.com/bytedance/mockey"
	"log"
	"time"
)

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
