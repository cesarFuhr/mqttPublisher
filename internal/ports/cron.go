package ports

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cesarFuhr/mqttPublisher/internal/app"
)

func NewCronPort(a app.Application, duration time.Duration) Cron {
	return Cron{
		app: a,
		tk:  time.NewTicker(duration),
	}
}

type Cron struct {
	app app.Application
	tk  *time.Ticker
}

func (c *Cron) Run(ctx context.Context) {
	for {
		select {
		case <-c.tk.C:
			err := c.app.Commands.NotifyStatus.Handle()
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			fmt.Println("exiting")
			return
		}
	}
}
