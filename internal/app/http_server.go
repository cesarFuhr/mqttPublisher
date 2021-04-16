package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// HTTPLogger http server logger
type HTTPLogger interface {
	Info(string, ...zap.Field)
}

// NewHTTPServer creates a new http handler
func NewHTTPServer(
	l HTTPLogger,
) *http.Server {
	router := mux.NewRouter()
	logger := newLoggerMiddleware(l)

	router.Use(logger)

	return &http.Server{
		Handler: router,
	}
}
