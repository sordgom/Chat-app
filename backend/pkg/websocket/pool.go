package websocket

import (
	"chat-go/model"
	"fmt"
	"log"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan *model.Chat
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *model.Chat),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: "bootup", Body: "New User Joined..."})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: "bootup", Body: "User Disconnected..."})
			}
			break
		case message := <-pool.Broadcast:
			// send to every client that is currently connected
			fmt.Println("new message", message)
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
				if client.Username == message.From || client.Username == message.To {
					err := client.Conn.WriteJSON(message)
					if err != nil {
						log.Printf("Websocket error: %s", err)
						client.Conn.Close()
						delete(pool.Clients, client)
					}
				}
			}
		}
	}
}
