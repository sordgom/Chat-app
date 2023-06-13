package main

import (
	"chat-go/pkg/videochat"
	"chat-go/pkg/websocket"
	"flag"
	"fmt"
	// "log"
	// "github.com/joho/godotenv"
)

// func init() {
// 	// Load the environment file .env
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Unable to Load the env file.", err)
// 	}
// }

func main() {
	server := flag.String("server", "", "http,websocket")
	flag.Parse()

	if *server == "chat" {
		fmt.Println("Distributed Chat App v0.01")
		websocket.StartWebsocketServer()
	} else if *server == "video" {
		fmt.Println("Video Call server is starting on :8081")
		videochat.SetupVideoChat()
	} else {
		fmt.Println("invalid server. Available server: chat or video")
	}
}
