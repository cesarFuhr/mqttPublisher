package config

import (
	"github.com/kelseyhightower/envconfig"
)

// LoadConfigs loads a configuration into an object
func LoadConfigs() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Config set of configurations needed to run the app
type Config struct {
	Server struct {
		Port string `envconfig:"SERVER_PORT"`
	}
	Broker struct {
		Host          string `envconfig:"MQTT_BROKER_HOST"`
		Port          int    `envconfig:"MQTT_BROKER_PORT"`
		User          string `envconfig:"MQTT_BROKER_USER"`
		Password      string `envconfig:"MQTT_BROKER_PASSWORD"`
		AutoReconnect bool   `envconfig:"MQTT_AUTORECONNECT"`
	}
	App struct {
		Workers struct {
			Number int `envconfig:"APP_WORKERS_NUMBER"`
		}
	}
}
