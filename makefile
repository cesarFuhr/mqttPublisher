# local development environments
SERVER_PORT=5000
MQTT_BROKER_HOST=localhost
MQTT_BROKER_PORT=17083
MQTT_BROKER_USER=
MQTT_BROKER_PASSWORD=
APP_WORKERS_NUMBER=1

APP_ENV_STRING = SERVER_PORT=$(SERVER_PORT) \
	MQTT_BROKER_HOST=$(MQTT_BROKER_HOST) \
	MQTT_BROKER_PORT=$(MQTT_BROKER_PORT) \
	APP_WORKERS_NUMBER=$(APP_WORKERS_NUMBER)

build:
	go build -o main ./cmd/main.go

install:
	go mod tidy
	go mod vendor

run: build
	./main

run-dev: build
	env $(APP_ENV_STRING) ./main

watch-dev: build
	env $(APP_ENV_STRING) air -c air.toml

start-local-broker:
	docker run --detach -p 127.0.0.1:$(MQTT_BROKER_PORT):$(MQTT_BROKER_PORT) \
		-p 127.0.0.1:18083:18083 \
		--env EMQX_LISTENER__TCP__EXTERNAL=$(MQTT_BROKER_PORT) \
		--name mqttbroker \
		emqx/emqx:4.2.11-alpine-amd64

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