package logger

import (
	"errors"
	"fmt"
	"sync"
)

var (
	factory logFactory
)

func init() {
	factory = logFactory{
		logInstanceList:  make(map[string]Log),
		appenderList:     make(map[string]Appender),
		appenderTypeList: make(map[string]CreateAppender),
	}

	initAppenderTypes()
	initStdAppender()
}

func initAppenderTypes() {
	factory.appenderTypeList["std"] = CreateConsoleAppender
	factory.appenderTypeList["file"] = CreateDailyRollingFileAppender
}

func initStdAppender() {
	if createAppenderFunc, ok := factory.appenderTypeList["std"]; ok {
		appender, _ := createAppenderFunc("std", make(map[string]string))
		factory.appenderList["std"] = appender
	}
}

type logFactory struct {
	logInstanceList  map[string]Log
	lock             sync.Mutex
	once             sync.Once
	appenderList     map[string]Appender
	appenderTypeList map[string]CreateAppender
}

func GetLog(pkg string) Log {
	factory.once.Do(initAppenders)

	factory.lock.Lock()
	defer factory.lock.Unlock()

	v, ok := factory.logInstanceList[pkg]
	if ok {
		return v
	}

	logConfig := GetLogConfig(pkg)

	log := createLog(pkg, logConfig)
	factory.logInstanceList[pkg] = log

	return log
}

func createLog(pkg string, logConfig LogConfig) Log {
	log := Log{
		id:           pkg,
		level:        logConfig.Level,
		appenderList: make(map[string]Appender),
	}

	for _, v := range logConfig.AppenderList {
		if appender, ok := factory.appenderList[v]; ok {
			log.AddAppender(v, appender)
		} else {
			fmt.Println("No find the appender[" + v + "] for logger[" + pkg + "]!")
		}
	}

	return log
}

func initAppenders() {
	err := initConfig()
	if err != nil {
		panic(err)
	}

	for id, appenderConfig := range conf.appenders {
		createAppender(id, appenderConfig)
	}
}

func createAppender(id string, appenderConfig AppenderConfig) {
	createAppenderFunc, ok := factory.appenderTypeList[appenderConfig.Type]
	if !ok {
		panic(errors.New("No find appender type[" + appenderConfig.Type + "] for appender[" + id + "]. Please registher the appender type before startup!"))
	}

	appender, err := createAppenderFunc(id, appenderConfig.Params)
	if err != nil {
		panic(err)
	}

	registerAppender(id, appender)
}

func RegisterAppenderType(id string, appenderFactory CreateAppender) {
	factory.lock.Lock()
	defer factory.lock.Unlock()

	_, existed := factory.appenderList[id]
	if existed {
		panic("The appender type[" + id + "] has existed!")
	}

	factory.appenderTypeList[id] = appenderFactory
}

func registerAppender(id string, appender Appender) {
	factory.lock.Lock()
	defer factory.lock.Unlock()

	_, existed := factory.appenderList[id]
	if existed {
		panic("The appender[" + id + "] has existed!")
	}

	factory.appenderList[id] = appender
}
