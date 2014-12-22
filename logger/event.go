package logger

import (
	"time"
)

type Event struct {
	logId    string
	level    Level
	time     time.Time
	message  string
}

func newEvent(logId string, level Level, message string, time time.Time) Event {
	return Event {
				logId: logId,
				level: level,
				message: message,
				time: time,
	}
}

func (this Event) GetLogId() string {
	return this.logId
}

func (this Event) GetLevel() Level {
	return this.level
}

func (this Event) GetTime() time.Time {
	return this.time
}

func (this Event) GetMessage() string {
	return this.message
}
