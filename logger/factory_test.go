package logger

import (
	"fmt"
	"strings"
	"testing"
)

func TestInitFunc(t *testing.T) {
	if _, ok := factory.appenderTypeList["std"]; !ok {
		t.Error("Not find std in Appender type list!")
	}

	if _, ok := factory.appenderTypeList["file"]; !ok {
		t.Error("Not find file in Appender type list!")
	}
}

func TestCreateAppender(t *testing.T) {
	RegisterAppenderType("testType", createTestAppenderTypeFunc)

	params := make(map[string]string)
	params["threshold"] = "info"
	params["path"] = "testAppender.log"

	config := appenderConfig{
		Id:     "testAppender",
		Type:   "testType",
		Params: params,
	}

	createAppender(config.Id, config)

	appender, ok := factory.appenderList[config.Id]
	if !ok {
		t.Error("No appender be created!")
	} else {
		if _, ok := appender.(*myTestAppender); !ok {
			t.Error(fmt.Sprintf("%s is not equals %s", appender.GetThreshold().GetName(), INFO.GetName()))
		}
	}
}

func TestCreateLog(t *testing.T) {
	appenerList := []string{"std", "testAppender"}

	config := logConfig{
		Id:           "github.com/xiaoyu830411",
		Level:        INFO,
		AppenderList: appenerList,
	}

	if _, ok := factory.appenderList["testAppender"]; !ok {
		TestCreateAppender(t)
		delete(factory.appenderTypeList, "testType")
		delete(factory.appenderList, "testAppender")
	}

	log := createLog("github.com/xiaoyu830411", config)
	if len(log.GetAppenderList()) < 2 {
		t.Error("Except 2 appenders")
	}

	for _, appener := range log.GetAppenderList() {
		if strings.EqualFold(appener.GetId(), "std") || strings.EqualFold(appener.GetId(), "testAppender") {

		} else {
			t.Error("Not match appenders[std or testAppender]")
		}
	}
}

func createTestAppenderTypeFunc(id string, params map[string]string) (Appender, error) {
	return &myTestAppender{
		id:        id,
		threshold: INFO,
	}, nil
}

type myTestAppender struct {
	id        string
	threshold Level
}

func (this myTestAppender) Fire(event Event) {
	if event.GetLevel().GetId() >= this.threshold.GetId() {
		fmt.Println(format(event))
	}
}

func (this myTestAppender) GetThreshold() Level {
	return this.threshold
}

func (this *myTestAppender) SetThreshold(threshold Level) {
	this.threshold = threshold
}

func (this myTestAppender) GetId() string {
	return this.id
}
