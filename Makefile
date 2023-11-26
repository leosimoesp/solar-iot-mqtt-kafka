include .env

start:
	docker compose up -d

stop:
	docker compose down

create-user:
	docker compose exec mosquitto mosquitto_passwd -b /tools/mosquitto/config/password.txt $(user) $(password)
	docker compose down 
	docker compose up -d mosquitto

delete-user:
	docker compose exec mosquitto mosquitto_passwd -D /tools/mosquitto/config/password.txt $(user)
	docker compose down
	docker compose up -d mosquitto

create-mqtt-clients:
	bash scripts/init.sh

run-tests:
	go test -v -count 1 ./...

build:
	go build -o bin/simulator cmd/main.go

run: 
	./bin/simulator