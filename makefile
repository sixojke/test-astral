build:
	go build -o app ./cmd/app

up: build
	sudo docker-compose up -d --build

down:
	sudo docker-compose down

restart: down build up

swag:
	swag init -g internal/app/app.go