package broker

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type BrokerCfg struct {
	Host     string
	Port     int
	ClientID string
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func NewBrokerClient(cfg BrokerCfg) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Host, cfg.Port))
	opts.SetClientID(cfg.ClientID)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	return mqtt.NewClient(opts), nil
}
