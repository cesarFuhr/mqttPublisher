package adapters

import (
	"fmt"

	"github.com/cesarFuhr/mqttPublisher/internal/domain/pid"
	"github.com/cesarFuhr/mqttPublisher/internal/domain/status"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	statusNotification := &StatusNotification{
		Status: s.Status,
		At:     timestamppb.New(s.At),
	}

	msg, err := proto.Marshal(statusNotification)
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

func (p *PIDPublisher) Publish(id string, pid pid.PID) error {

	pidNotification := &PIDNotification{
		EventID:     uuid.NewString(),
		Value:       pid.Value,
		At:          timestamppb.New(pid.At),
		Description: pid.Description,
		Unit:        pid.Unit,
	}

	msg, err := proto.Marshal(pidNotification)
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
