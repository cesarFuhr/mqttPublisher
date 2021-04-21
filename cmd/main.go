package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/cesarFuhr/mqttPublisher/internal/adapters"
	"github.com/cesarFuhr/mqttPublisher/internal/app"
	"github.com/cesarFuhr/mqttPublisher/internal/app/command"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/broker"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/config"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/exit"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/logger"
	"github.com/cesarFuhr/mqttPublisher/internal/pkg/server"
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

	application, appTeardown := newApplication(cfg)

	server := newServer(cfg, application)
	cron := ports.NewCronPort(application, time.Second)

	go gracefullShutdown(ctx, e, appTeardown, server)

	go cron.Run(ctx)
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Goodbye...")
			return
		}
		log.Fatalf("could not listen on port %s ... %v", cfg.Server.Port, err)
	}
}

func gracefullShutdown(ctx context.Context, e chan struct{}, teardown func(), server *http.Server) {
	<-e
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println(err)
	}
	teardown()
}

func newServer(cfg config.Config, a app.Application) *http.Server {
	logger := logger.NewLogger()
	server := server.NewHTTPServer(logger, ports.NewHttpPort(a))
	server.Addr = ":" + cfg.Server.Port
	return server
}

func newApplication(cfg config.Config) (app.Application, func()) {

	mqttClient := setupMQTTClient(cfg)
	statusPublisher := adapters.NewStatusPublisher(mqttClient)
	pidPublisher := adapters.NewPIDPublisher(mqttClient)

	license := "ISS-1312"

	return app.Application{
			Commands: app.Commands{
				NotifyStatus: command.NewStatusHandler(license, &statusPublisher),
				NotifyPIDs:   command.NewPIDsHandler(license, &pidPublisher),
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
