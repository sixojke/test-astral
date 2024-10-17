build:
	go build -o app ./cmd/app

up:
	sudo docker-compose up -d --build

down:
	sudo docker-compose down

restart: down up