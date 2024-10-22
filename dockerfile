FROM golang:latest

WORKDIR /app

COPY ./app ./
COPY ./.env ./
COPY ./configs ./configs
COPY ./schema ./schema

EXPOSE 8080

CMD ["./app"]