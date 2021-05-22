# local development environments
SERVER_PORT=5000
MQTT_BROKER_HOST=192.46.219.17
MQTT_BROKER_PORT=1883
MQTT_AUTORECONNECT=true
MQTT_BROKER_USER=
MQTT_BROKER_PASSWORD=
PUBLISHER_QOS=1
APP_VEHICLE_LICENSE=TST-1234
APP_WORKERS_NUMBER=1

APP_ENV_STRING = SERVER_PORT=$(SERVER_PORT) \
	MQTT_BROKER_HOST=$(MQTT_BROKER_HOST) \
	MQTT_BROKER_PORT=$(MQTT_BROKER_PORT) \
	MQTT_AUTORECONNECT=$(MQTT_AUTORECONNECT) \
	PUBLISHER_QOS=$(PUBLISHER_QOS) \
	APP_VEHICLE_LICENSE=$(APP_VEHICLE_LICENSE) \
	APP_WORKERS_NUMBER=$(APP_WORKERS_NUMBER)

build:
	go build -o main ./cmd/main.go

build-docker:
	docker build --tag=mqttpub:latest -f builds/Dockerfile .

install:
	go mod tidy
	go mod vendor

run: build
	./main

run-dev: build
	env $(APP_ENV_STRING) ./main

run-docker:
	docker run --detach \
		--env SERVER_PORT=$(SERVER_PORT) \
		--env MQTT_BROKER_HOST=$(MQTT_BROKER_HOST) \
		--env MQTT_BROKER_PORT=$(MQTT_BROKER_PORT) \
		--env MQTT_AUTORECONNECT=$(MQTT_AUTORECONNECT) \
		--env PUBLISHER_QOS=$(PUBLISHER_QOS) \
		--env APP_WORKERS_NUMBER=$(APP_WORKERS_NUMBER) \
		-p 5000:5000 \
		--name mqttpub \
		mqttpub:latest

watch-dev: build
	env $(APP_ENV_STRING) air -c air.toml

start-local-broker:
	docker run --detach -p 127.0.0.1:$(MQTT_BROKER_PORT):$(MQTT_BROKER_PORT) \
		-p 127.0.0.1:18083:18083 \
		--env EMQX_LISTENER__TCP__EXTERNAL=$(MQTT_BROKER_PORT) \
		--name mqttbroker \
		emqx/emqx:4.2.11-alpine-amd64

stop-docker:
	docker stop mqttpub
	docker rm mqttpub

stop-local-broker:
	docker stop mqttbroker
	docker rm mqttbroker

test-unit:
	go test ./internal/...

test-full:
	docker-compose -f docker-compose.test.yml up -d broker
	docker-compose -f docker-compose.test.yml up --build test
	docker-compose -f docker-compose.test.yml down

watch-test:
	watcher -cmd="make test-unit" -keepalive=true

watch-test-full:
	watcher -cmd="make test-full" -keepalive=true
