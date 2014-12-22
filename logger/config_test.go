package logger

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t1 *testing.T) {
	path := "log4g_test.properties"

	c, err := loadConfig(path)
	if err != nil {
		t1.Fatal(err)
	}

	fmt.Println(c)
}

func TestGetCurrentWorkDirectory(t *testing.T) {
	_, err := getCurrentWorkDirectory()
	if err != nil {
		t.Error(err)
	}
}
