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
	RegisterAppenderType("testType", CreateTestAppenderTypeFunc)

	params := make(map[string]string)
	params["threshold"] = "info"
	params["path"] = "testAppender.log"

	appenderConfig := AppenderConfig{
		Id:     "testAppender",
		Type:   "testType",
		Params: params,
	}

	createAppender(appenderConfig.Id, appenderConfig)

	appender, ok := factory.appenderList[appenderConfig.Id]
	if !ok {
		t.Error("No appender be created!")
	} else {
		if _, ok := appender.(*TestAppender); !ok {
			t.Error(fmt.Sprintf("%s is not equals %s", appender.GetThreshold().GetName(), INFO.GetName()))
		}
	}
}

func TestCreateLog(t *testing.T) {
	appenerList := []string{"std", "testAppender"}

	logConfig := LogConfig{
		Id:           "github.com/xiaoyu830411",
		Level:        INFO,
		AppenderList: appenerList,
	}

	if _, ok := factory.appenderList["testAppender"]; !ok {
		TestCreateAppender(t)
		delete(factory.appenderTypeList, "testType")
		delete(factory.appenderList, "testAppender")
	}

	log := createLog("github.com/xiaoyu830411", logConfig)
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

func CreateTestAppenderTypeFunc(id string, params map[string]string) (Appender, error) {
	return &TestAppender{
		id:        id,
		threshold: INFO,
	}, nil
}

type TestAppender struct {
	id        string
	threshold Level
}

func (this TestAppender) Fire(event Event) {
	if event.GetLevel().GetId() >= this.threshold.GetId() {
		fmt.Println(format(event))
	}
}

func (this TestAppender) GetThreshold() Level {
	return this.threshold
}

func (this *TestAppender) SetThreshold(threshold Level) {
	this.threshold = threshold
}

func (this TestAppender) GetId() string {
	return this.id
}
