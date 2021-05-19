package ports

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/cesarFuhr/mqttPublisher/internal/app"
	"github.com/cesarFuhr/mqttPublisher/internal/app/command"
)

func NewHttpPort(a app.Application) Http {
	return Http{
		app: a,
	}
}

type Http struct {
	app app.Application
}

func (h *Http) PublishDTCs(w http.ResponseWriter, r *http.Request) {
	var o []DTC

	err := decodeJSONBody(r, &o)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			replyJSON(w, mr.status, HTTPError{
				Message: mr.msg,
			})
			return
		}
		replyJSON(w, http.StatusInternalServerError, HTTPError{
			Message: err.Error(),
		})
		return
	}

	if err := h.app.Commands.NotifyDTCs.Handle(httpDTCToCommand(o)); err != nil {
		replyJSON(w, http.StatusInternalServerError, HTTPError{
			Message: err.Error(),
		})
		return
	}

	replyJSON(w, http.StatusOK, o)
}

func (h *Http) PublishPIDs(w http.ResponseWriter, r *http.Request) {
	var o []PID

	err := decodeJSONBody(r, &o)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			replyJSON(w, mr.status, HTTPError{
				Message: mr.msg,
			})
			return
		}
		replyJSON(w, http.StatusInternalServerError, HTTPError{
			Message: err.Error(),
		})
		return
	}

	if err := h.app.Commands.NotifyPIDs.Handle(httpPIDToCommand(o)); err != nil {
		replyJSON(w, http.StatusInternalServerError, HTTPError{
			Message: err.Error(),
		})
		return
	}

	replyJSON(w, http.StatusOK, o)
}

type DTC struct {
	DTC         string    `json:"dtc"`
	At          time.Time `json:"at"`
	Description string    `json:"description"`
}

type PID struct {
	PID         string    `json:"pid"`
	At          time.Time `json:"at"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	Unit        string    `json:"unit"`
}

func httpDTCToCommand(dtcs []DTC) []command.DTCCommand {
	commands := []command.DTCCommand{}
	for _, v := range dtcs {
		commands = append(commands, command.DTCCommand{
			DTC:         v.DTC,
			At:          v.At,
			Description: v.Description,
		})
	}
	return commands
}

func httpPIDToCommand(pids []PID) []command.PIDCommand {
	commands := []command.PIDCommand{}
	for _, v := range pids {
		commands = append(commands, command.PIDCommand{
			PID:         v.PID,
			At:          v.At,
			Value:       v.Value,
			Description: v.Description,
			Unit:        v.Unit,
		})
	}
	return commands
}

type HTTPError struct {
	Message string `json:"message"`
}

func replyJSON(w http.ResponseWriter, c int, o interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(c)
	json.NewEncoder(w).Encode(o)
}
