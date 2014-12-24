package logger

import (
	"fmt"
	"testing"
)

func TestGetFormattedMessage(t *testing.T) {
	format := "I am xiaoyu"
	message := getFormattedMessage(format)
	fmt.Println(message)
}
