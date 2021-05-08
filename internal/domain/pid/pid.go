package pid

import (
	"time"
)

type PID struct {
	PID         string
	At          time.Time
	Value       string
	Description string
	Unit        string
}
