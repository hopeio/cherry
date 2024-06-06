package log

import "testing"

func TestLog(t *testing.T) {
	Info("test")
}

func TestLogStack(t *testing.T) {
	ErrorStack("test")
}
