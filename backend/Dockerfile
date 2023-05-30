FROM golang:1.12-alpine

RUN apk add --no-cache bash ca-certificates git gcc g++ libc-dev

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go mod download
RUN go build -v

