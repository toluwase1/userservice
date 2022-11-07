# setup base image
FROM golang:1.17.0-alpine

WORKDIR /app

COPY ./ /app

COPY ./db/01-init.sh /docker-entrypoint-initdb.d/

RUN go mod tidy

ENTRYPOINT [ "go", "run", "main.go" ]