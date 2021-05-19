package command

import (
	"context"
	"time"

	"github.com/cesarFuhr/mqttPublisher/internal/domain/dtc"
)

func NewDTCHandler(id string, p dtc.Publisher) NotifyDTCsHandler {

	return NotifyDTCsHandler{
		ID:  id,
		Pub: p,
	}
}

type NotifyDTCsHandler struct {
	ID  string
	Pub dtc.Publisher
}

type DTCCommand struct {
	DTC         string
	At          time.Time
	Description string
}

func (h *NotifyDTCsHandler) Handle(pids []DTCCommand) error {
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

func (s *NotifyDTCsHandler) publishAsync(ctx context.Context, res chan error, cmd DTCCommand) {
	res <- s.Pub.Publish(s.ID, dtc.DTC{
		DTC:         cmd.DTC,
		At:          cmd.At,
		Description: cmd.Description,
	})
}
