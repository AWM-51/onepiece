# OnePiece

## 一.简介
OnePiece 是杭州测试团研发的 Go 组件库，包含了丰富的组件封装，封装的组件包括但不限于断言处理、JSON 处理、日志打印、加密等内容。
### 1.1背景
公司全面转go持续进行中，不管研发和测试都应有相应的准备，所以为了提高相率，对于业务中常用工具进行go的组件化，实现通用化
## 二.快速开始
### 2.1 安装使用
目前还没有打成relase包发布,不能使用
#### 安装

```
go get git.duowan.com:opensource/hz-tester/onepiece
```

#### 使用
```
import (
    "git.duowan.com:opensource/hz-tester/onepiece/jsonutil"
)
```
### 2.2 组件列表
| 组件名 |组件功能  |贡献者  |版本  |
| --- | --- | --- | --- |
| jsonutil | 提供 JSON 数据的解析、比较、排序等功能 | zwj |v1  |
|  logutil| 提供日志打印功能。 |zwj  | v1 |
| encryptutil | 提供常见的加解密算法的实现。 | zwj | v1 |
| assertutil | 提供基于go原生testing的断言封装 |zwj  |v1|
| ... |  |  |  |
## 三.更多
### 3.1贡献者

* zwj

### 3.2许可证
MIT License. 详细内容参见 LICENSE 文件。