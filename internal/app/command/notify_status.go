package command

import (
	"time"

	"github.com/cesarFuhr/mqttPublisher/internal/domain/status"
)

func NewStatusHandler(id string, p status.Publisher) NotifyStatusHandler {
	return NotifyStatusHandler{
		ID:  id,
		Pub: p,
	}
}

type NotifyStatusHandler struct {
	ID  string
	Pub status.Publisher
}

func (s *NotifyStatusHandler) Handle() error {
	return s.Pub.Publish(s.ID, status.Status{
		Status: true,
		At:     time.Now(),
	})
}
