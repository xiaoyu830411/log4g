package logger

import (
	"fmt"
	"sync"
	"time"
)

type Log struct {
	id           string
	level        Level
	appenderList map[string]Appender
	lock         sync.RWMutex
	inherit      string
}

func (this Log) GetId() string {
	return this.id
}

func (this *Log) GetLevel() Level {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.level
}

func (this *Log) IsTrace() bool {
	return this.GetLevel().GetId() <= TRACE.GetId()
}

func (this *Log) IsDebug() bool {
	return this.GetLevel().GetId() <= DEBUG.GetId()
}

func (this *Log) IsInfo() bool {
	return this.GetLevel().GetId() <= INFO.GetId()
}

func (this *Log) IsWarn() bool {
	return this.GetLevel().GetId() <= WARN.GetId()
}

func (this *Log) IsError() bool {
	return this.GetLevel().GetId() <= ERROR.GetId()
}

func (this *Log) IsFatal() bool {
	return this.GetLevel().GetId() <= FATAL.GetId()
}

func (this *Log) fire(level Level, message string) {
	event := newEvent(this.GetId(), level, message, time.Now())

	this.lock.RLock()
	defer this.lock.RUnlock()

	for _, appender := range this.appenderList {
		appender.Fire(event)
	}
}

func (this *Log) Trace(format string, value ...interface{}) {
	if !this.IsTrace() {
		return
	}

	message := fmt.Sprintf(format, value)
	this.fire(TRACE, message)
}

func (this *Log) Debug(format string, value ...interface{}) {
	if !this.IsDebug() {
		return
	}

	message := fmt.Sprintf(format, value)
	this.fire(DEBUG, message)
}

func (this *Log) Info(format string, value ...interface{}) {
	if !this.IsInfo() {
		return
	}

	message := fmt.Sprintf(format, value)
	this.fire(INFO, message)
}

func (this *Log) Warn(format string, value ...interface{}) {
	if !this.IsWarn() {
		return
	}

	message := fmt.Sprintf(format, value)
	this.fire(WARN, message)
}

func (this *Log) Error(format string, value ...interface{}) {
	if !this.IsError() {
		return
	}

	message := fmt.Sprintf(format, value)
	this.fire(ERROR, message)
}

func (this *Log) Fatal(format string, value ...interface{}) {
	if !this.IsFatal() {
		return
	}

	message := fmt.Sprintf(format, value)
	this.fire(FATAL, message)
}

func (this *Log) setLevel(level Level) {
	this.lock.Lock()
	this.lock.Unlock()

	this.level = level
}

func (this *Log) GetAppenderList() map[string]Appender {
	this.lock.RLock()
	defer this.lock.RUnlock()

	tmp := make(map[string]Appender)
	for k, v := range this.appenderList {
		tmp[k] = v
	}

	return tmp
}

func (this *Log) AddAppender(key string, appender Appender) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.appenderList[key] = appender
}

func (this *Log) RemoveAppender(key string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.appenderList, key)
}
