package main

import (
	"context"
	"time"

	"github.com/cesarFuhr/gocrypto/internal/pkg/config"
	"github.com/cesarFuhr/gocrypto/internal/pkg/exit"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/broker"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func main() {
	run()
}

func run() {
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}

	mqttClient := setupMQTTClient(cfg)

	e := make(chan struct{}, 1)
	exit.ListenToExit(e)

	go gracefullShutdown(e)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func gracefullShutdown(e chan struct{}) {
	<-e
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	cancel()
}

func setupMQTTClient(cfg config.Config) mqtt.Client {
	mqttCfg := broker.BrokerCfg{
		Host:     cfg.Broker.Host,
		Port:     cfg.Broker.Port,
		ClientID: uuid.NewString(),
	}

	cli, err := broker.NewBrokerClient(mqttCfg)
	if err != nil {
		panic(err)
	}

	return cli
}
