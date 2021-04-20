package status

import (
	"time"
)

type Status struct {
	At     time.Time
	Status bool
}
