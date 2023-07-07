# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

WORKDIR /api

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /api
EXPOSE 8080 8081 8082

CMD [ "/api" ]
