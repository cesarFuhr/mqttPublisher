package command

import (
	"context"
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
	PID         string
	Value       string
	At          time.Time
	Description string
	Unit        string
}

func (h *NotifyPIDsHandler) Handle(pids []PIDCommand) error {
	results := make(chan error, len(pids))
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer func() {
		cancel()
		close(results)
	}()

	for _, v := range pids {
		go h.publishAsync(ctx, results, v)
	}

	var err error
	for i := 0; i < len(pids); i++ {
		err = <-results
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *NotifyPIDsHandler) publishAsync(ctx context.Context, res chan error, cmd PIDCommand) {
	res <- s.Pub.Publish(s.ID, pid.PID{
		PID:         cmd.PID,
		Value:       cmd.Value,
		At:          cmd.At,
		Description: cmd.Description,
		Unit:        cmd.Unit,
	})
}
