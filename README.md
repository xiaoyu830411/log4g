Log
==========
这是一个golang的log lib

#安装
```shell
go get -u github.com/xiaoyu830411/log4g
```

#快速使用
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
