package adapters

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cesarFuhr/mqttPublisher/internal/domain/pid"
	"github.com/cesarFuhr/mqttPublisher/internal/domain/status"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func NewStatusPublisher(c mqtt.Client) StatusPublisher {
	return StatusPublisher{
		client: c,
		qos:    1,
	}
}

type StatusPublisher struct {
	client mqtt.Client
	qos    int
}

type statusNotification struct {
	At     time.Time
	Status bool
}

func (p *StatusPublisher) Publish(id string, s status.Status) error {

	msg, err := json.Marshal(statusNotification{
		At:     s.At,
		Status: s.Status,
	})
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("carMon/%s/status", id)
	token := p.client.Publish(topic, byte(p.qos), false, msg)

	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func NewPIDPublisher(c mqtt.Client, qos byte) PIDPublisher {
	return PIDPublisher{
		client: c,
		qos:    qos,
	}
}

type PIDPublisher struct {
	client mqtt.Client
	qos    byte
}

type PIDNotification struct {
	EventID string
	At      time.Time
	Value   string
}

func (p *PIDPublisher) Publish(id string, pid pid.PID) error {
	msg, err := json.Marshal(PIDNotification{
		EventID: uuid.NewString(),
		At:      pid.At,
		Value:   pid.Value,
	})
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("carMon/%s/param/%s", id, pid.PID)
	token := p.client.Publish(topic, byte(p.qos), false, msg)

	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
