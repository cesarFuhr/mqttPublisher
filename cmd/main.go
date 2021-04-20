package main

import (
	"context"
	"time"

	"github.com/cesarFuhr/mqttPublisher/internal/adapters"
	"github.com/cesarFuhr/mqttPublisher/internal/app"
	"github.com/cesarFuhr/mqttPublisher/internal/app/command"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/broker"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/config"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/exit"
	"github.com/cesarFuhr/mqttPublisher/internal/ports"
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

	ctx := context.Background()
	e := make(chan struct{}, 1)
	exit.ListenToExit(e)

	application, teardown := newApplication(cfg)

	cron := ports.NewCronPort(application, time.Second)
	go cron.Run(ctx)

	gracefullShutdown(ctx, e, teardown)
}

func gracefullShutdown(ctx context.Context, e chan struct{}, teardown func()) {
	<-e
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	teardown()
}

func newApplication(cfg config.Config) (app.Application, func()) {

	mqttClient := setupMQTTClient(cfg)
	statusPublisher := adapters.NewStatusPublisher(mqttClient)

	return app.Application{
			Commands: app.Commands{
				NotifyStatus: command.NewStatusHandler("ISS-1312", &statusPublisher),
			},
		}, func() {
			mqttClient.Disconnect(1000)
		}
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
