package command

import (
	"time"

	"github.com/cesarFuhr/mqttPublisher/internal/domain/pid"
)

func NewPIDsHandler(id string, p pid.Publisher) NotifyPIDsHandler {
	return NotifyPIDsHandler{
		ID:  id,
		Pub: p,
	}
}

type NotifyPIDsHandler struct {
	ID  string
	Pub pid.Publisher
}

type PIDCommand struct {
	PID   string
	Value string
	At    time.Time
}

func (s *NotifyPIDsHandler) Handle(pids []PIDCommand) error {
	var err error
	for _, v := range pids {
		err = s.Pub.Publish(s.ID, pid.PID{
			PID:   v.PID,
			Value: v.Value,
			At:    v.At,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
