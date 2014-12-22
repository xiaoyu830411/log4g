package logger

import (
	"errors"
	"strconv"
	"strings"
)

var (
	OFF                                         Level = newLevel(int(^uint(0)>>1), "OFF")
	ALL, TRACE, DEBUG, INFO, WARN, ERROR, FATAL Level = newLevel(-(OFF.GetId() - 1), "ALL"),
		newLevel(0, "TRACE"),
		newLevel(1, "DEBUG"),
		newLevel(2, "INFO"),
		newLevel(3, "WARN"),
		newLevel(4, "ERROR"),
		newLevel(5, "FATAL")

	ALL_LEVELS = [...]Level{ALL, TRACE, DEBUG, INFO, WARN, ERROR, FATAL, OFF}
)

type Level struct {
	id   int
	name string
}

func newLevel(id int, name string) Level {
	return Level{
		id:   id,
		name: name,
	}
}

func (this Level) GetId() int {
	return this.id
}

func (this Level) GetName() string {
	return this.name
}

func GetLevelByName(name string) (*Level, error) {
	for _, level := range ALL_LEVELS {
		if strings.EqualFold(strings.ToUpper(name), level.GetName()) {
			return &level, nil
		}
	}

	return nil, errors.New("Not find level by name[" + name + "]!")
}

func GetLevelById(id int) (*Level, error) {
	for _, level := range ALL_LEVELS {
		if level.GetId() == id {
			return &level, nil
		}
	}

	return nil, errors.New("Not find level by ID[" + strconv.Itoa(id) + "]!")
}
