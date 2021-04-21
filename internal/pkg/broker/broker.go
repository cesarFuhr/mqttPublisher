package broker

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type BrokerCfg struct {
	Host          string
	Port          int
	ClientID      string
	AutoReconnect bool
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected to Broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connection lost: %v", err)
}

var reconCount = 0

var reconnectionHandler mqtt.ReconnectHandler = func(client mqtt.Client, opts *mqtt.ClientOptions) {
	log.Printf("Reconnecting count %d", reconCount)
	reconCount++
}

func NewBrokerClient(cfg BrokerCfg) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Host, cfg.Port))
	opts.SetClientID(cfg.ClientID)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetAutoReconnect(cfg.AutoReconnect)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.OnReconnecting = reconnectionHandler

	return mqtt.NewClient(opts), nil
}

func Connect(cli mqtt.Client) error {
	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
