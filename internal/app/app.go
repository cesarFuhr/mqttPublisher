package app

import "github.com/cesarFuhr/mqttPublisher/internal/app/command"

type Application struct {
	Commands Commands
}

type Commands struct {
	NotifyStatus command.NotifyStatusHandler
	NotifyPIDs   command.NotifyPIDsHandler
	NotifyDTCs   command.NotifyDTCsHandler
}
