package ports

import (
	"encoding/json"
	"net/http"

	"github.com/cesarFuhr/mqttPublisher/internal/app"
)

func NewHttpPort(a app.Application) Http {
	return Http{
		app: a,
	}
}

type Http struct {
	app app.Application
}

func (h *Http) PublishPIDs(w http.ResponseWriter, r *http.Request) {
	replyJSON(w, http.StatusOK, struct{ test string }{
		test: "response",
	})
}

func replyJSON(w http.ResponseWriter, c int, o interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(c)
	json.NewEncoder(w).Encode(o)
}
