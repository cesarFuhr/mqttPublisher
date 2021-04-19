package main

import (
	"context"
	"log"
	"time"

	"github.com/cesarFuhr/mqttPublisher/internal/app/adapters"
	"github.com/cesarFuhr/mqttPublisher/internal/app/status"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/broker"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/config"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/exit"
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

	statusPublisher := adapters.NewStatusPublisher(mqttClient)
	statusService := status.NewStatusService("ISS-9811", &statusPublisher)

	tk := time.NewTicker(time.Second)

ExecutionLoop:
	for {
		select {
		case <-tk.C:
			log.Println("Status!")
			err := statusService.Send()
			if err != nil {
				log.Println(err)
			}
		case <-e:
			gracefullShutdown(mqttClient)
			log.Println("Finished")
			break ExecutionLoop
		}
	}
}

func gracefullShutdown(mqtt mqtt.Client) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	mqtt.Disconnect(1000)

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

	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return cli
}
