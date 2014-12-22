package logger

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CreateAppender func(string, map[string]string) (Appender, error)

type Appender interface {
	GetId() string
	Fire(event Event)
	GetThreshold() Level
	SetThreshold(threshold Level)
}

type consoleAppender struct {
	id        string
	threshold Level
	lock      sync.RWMutex
}

func CreateConsoleAppender(id string, params map[string]string) (Appender, error) {
	var threshold Level = ALL
	if v, ok := params["threshold"]; ok {
		if level, err := GetLevelByName(v); err == nil {
			threshold = *level
		}
	}

	return &consoleAppender{id: id, threshold: threshold}, nil
}

func (this consoleAppender) Fire(event Event) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if canbeOutputted(event.GetLevel(), this.threshold) {
		fmt.Println(format(event))
	}
}

func (this consoleAppender) GetThreshold() Level {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.threshold
}

func (this *consoleAppender) SetThreshold(threshold Level) {
	this.lock.Lock()
	this.lock.Unlock()

	this.threshold = threshold
}

func (this consoleAppender) GetId() string {
	return this.id
}

//Event Format

func format(event Event) string {
	return event.GetTime().Format("Jan 02 15:04:05") + " " +
		event.GetLevel().GetName() + " " +
		event.GetLogId() + ": " +
		event.GetMessage()
}

//Daily Rolling File

type dailyRollingFileAppender struct {
	id        string
	threshold Level
	path      string
	date      int
	file      *os.File
	lock      sync.RWMutex
}

func CreateDailyRollingFileAppender(id string, params map[string]string) (Appender, error) {
	var threshold Level = ALL
	if v, ok := params["threshold"]; ok {
		if level, err := GetLevelByName(v); err == nil {
			threshold = *level
		}
	}

	var path string
	if v, ok := params["path"]; ok {
		if strings.EqualFold(v[0:1], string(os.PathSeparator)) {
			path = v
		} else {
			if dir, err := os.Getwd(); err != nil {
				return nil, err
			} else {
				path = dir + string(os.PathSeparator) + v
			}
		}
	} else {
		return nil, errors.New("The path parameter is required! Please specify its value in the configuration.")
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 666)
	if err != nil {
		return nil, err
	}

	appender := dailyRollingFileAppender{
		id:        id,
		threshold: threshold,
		path:      path,
		date:      GetDate(),
		file:      file,
	}

	return &appender, nil
}

func (this dailyRollingFileAppender) Fire(event Event) {
	this.lock.Lock()
	defer this.lock.Lock()

	if GetDate() > this.date {
		if err := this.roll(); err != nil {
			fmt.Println(err)
		}
	}

	if event.GetLevel().GetId() >= this.threshold.GetId() {
		this.output(event)
	}
}

func (this *dailyRollingFileAppender) roll() error {
	current := GetDate()
	if current > this.date {
		file := this.file
		file.Close()

		path := this.path

		os.Rename(path, this.path+"."+strconv.Itoa(this.date))
		file, err := os.OpenFile(this.path, os.O_APPEND|os.O_CREATE, 666)
		if err != nil {
			return err
		}

		this.file = file
		this.date = current
	}

	return nil
}

func (this *dailyRollingFileAppender) output(event Event) {
	this.file.WriteString(format(event) + "\n")
}

func (this dailyRollingFileAppender) GetPath() string {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.path
}

func (this *dailyRollingFileAppender) SetPath(path string) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	file := this.file
	file.Close()

	file, err := os.OpenFile(this.path, os.O_APPEND|os.O_CREATE, 666)
	if err != nil {
		return err
	}

	this.path = path
	this.file = file

	return nil
}

func (this dailyRollingFileAppender) GetThreshold() Level {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.threshold
}

func (this *dailyRollingFileAppender) SetThreshold(threshold Level) {
	this.lock.Lock()
	this.lock.Unlock()

	this.threshold = threshold
}

func (this dailyRollingFileAppender) GetId() string {
	return this.id
}

func GetDate() int {
	year, month, day := time.Now().Date()
	return (year*10000 + int(month)*100 + day)
}

func canbeOutputted(src Level, target Level) bool {
	return src.GetId() >= target.GetId()
}
