package status

import (
	"time"
)

func NewStatusService(id string, p Publisher) *StatusService {
	return &StatusService{
		ID:  id,
		Pub: p,
	}
}

type StatusService struct {
	ID  string
	Pub Publisher
}

type Status struct {
	At     time.Time
	Status bool
}

type Publisher interface {
	Publish(string, Status) error
}

func (s *StatusService) Send() error {
	return s.Pub.Publish(s.ID, Status{
		Status: true,
		At:     time.Now(),
	})
}
