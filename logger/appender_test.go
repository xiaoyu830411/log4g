package logger

import (
	"fmt"
	"os"
	"testing"
)

func TestCanbeOutputted(t *testing.T) {
	if canbeOutputted(INFO, DEBUG) {
	} else {
		t.Error("INFO can be outputted [DEBUG]!")
	}

	if canbeOutputted(DEBUG, INFO) {
		t.Error("DEBUG can not be outputted [INFO]!")
	}

	if canbeOutputted(INFO, INFO) {
	} else {
		t.Error("INFO can be outputted [INFO]!")
	}
}

func TestCreateConsoleAppender(t *testing.T) {
	params := make(map[string]string)
	params["threshold"] = "ERROR"
	appender, err := CreateConsoleAppender("myTest", params)
	if err != nil {
		t.Error(err)
	}

	if _, ok := appender.(*consoleAppender); !ok {
		t.Error(fmt.Sprintf("Except a consoleAppender , But it is %v", appender))
	}

	if appender.GetThreshold() != ERROR {
		t.Error(fmt.Sprintf("Except true, But false[%v, %v]", appender.GetThreshold(), ERROR))
	}
}

func TestCreateDailyRollingFileAppender(t *testing.T) {
	params := make(map[string]string)
	params["threshold"] = "INFO"
	params["path"] = "myTest2.log"

	defer os.Remove("myTest2.log")
	appender, err := CreateDailyRollingFileAppender("myTest2", params)

	if err != nil {
		t.Error(err)
	}

	if _, ok := appender.(*dailyRollingFileAppender); !ok {
		t.Error(fmt.Sprintf("Except a dailyRollingFileAppender , But it is %v", appender))
	}

	if appender.GetThreshold() != INFO {
		t.Error(fmt.Sprintf("Except true, But false[%v, %v]", appender.GetThreshold(), ERROR))
	}
}
