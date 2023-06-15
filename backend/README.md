# Backend

## Authentication 

In order to run the auth service on 8080:

`go run main.go --server=http`

## Videocall

In order to run the videochat service on 8081:

`go run main.go --server=video`

## Chat

In order to run the chat service on 8082:

`go run main.go --server=chat`

### Docker build

Make sure you have docker installed. Then, run the following commands:

`docker build -t backend .`
