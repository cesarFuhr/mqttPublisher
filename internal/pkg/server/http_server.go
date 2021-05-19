package server

import (
	"net/http"

	"github.com/cesarFuhr/mqttPublisher/internal/ports"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// HTTPLogger http server logger
type HTTPLogger interface {
	Info(string, ...zap.Field)
}

// NewHTTPServer creates a new http handler
func NewHTTPServer(l HTTPLogger, p ports.Http) *http.Server {
	router := mux.NewRouter()
	logger := newLoggerMiddleware(l)

	router.Use(logger)

	router.HandleFunc("/publish/pids", p.PublishPIDs).
		Methods(http.MethodPost)
	router.HandleFunc("/publish/dtcs", p.PublishDTCs).
		Methods(http.MethodPost)

	return &http.Server{
		Handler: router,
	}
}
