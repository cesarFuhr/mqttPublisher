package server

import (
	"testing"

	"go.uber.org/zap"
)

type loggerStub struct {
	CalledWith []interface{}
	Called     bool
}

func (l *loggerStub) Info(msg string, args ...zap.Field) {
	l.CalledWith = append(l.CalledWith, msg)
	l.CalledWith = append(l.CalledWith, args)
	l.Called = true
}

var (
	log    = new(loggerStub)
	server = NewHTTPServer(log).Handler
)

func assertValue(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}
