Log
==========
这是一个golang的log lib

#安装
```shell
go get -u github.com/xiaoyu830411/log4g
```

快速使用
==========
#配置文件(log4g.properties 放到项目下)
```shell
[appender]
myFile.type = file
myFile.threshold = info
myFile.path = /var/logs/my.log

[logger]
__root__ = debug, std
github.com/xiaoyu830411 = info, myFile

```
#代码使用

```golang
package main

import (
	"github.com/xiaoyu830411/log4g/logger"
)

var (
    log := logger.GetLog("main")
)

func main() {
    log.info("My First Log")
}
```

深入了解
==========

#注册自己的Appender类型

```golang
//实现创建的接口（第一个是被创建的Appender的Id，第二个是创建的时候所需要的其他参数），请保证线程安全
//type CreateAppender func(string, map[string]string) (Appender, error)

func MyType (id string, params[string]string) (Appender, error) {
	//...
}

//实现Appender接口, 请保证线程安全
type Appender interface {
	GetId() string
	Fire(event Event)
	GetThreshold() Level
	SetThreshold(threshold Level)
}
```

#main 文件

```golang
package main

import (
    "github.com/xiaoyu830411/log4g/logger"
)

func init() {
	logger.RegisterAppenderType("myType", MyType)
}

func main() {
    //....
}
```

#配置文件(log4g.properties 放到项目下)
```shell
[appender]
myAppender.type = myType
myAppender.threshold = info
myAppender.参数1 = 参数值1
myAppender.参数2 = 参数值2

[logger]
__root__ = debug, std
github.com/xiaoyu830411 = info, myAppender
```
#代码使用
```golang
package myPackage

import (
    "github.com/xiaoyu830411/log4g/logger"
)

var (
	log := logger.GetLog("github.com/xiaoyu830411/myPackage")
)

func method1() {
    log.info("I am %v", "Xiaoyu")
}
```