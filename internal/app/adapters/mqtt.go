package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/cesarFuhr/mqttPublisher/internal/app/status"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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

func (p *StatusPublisher) Publish(id string, s status.Status) error {

	msg, err := json.Marshal(s)
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
