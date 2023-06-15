package websocket

import (
	"chat-go/model"
	"chat-go/pkg/redisrepo"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string
	Conn     *websocket.Conn
	Pool     *Pool
	Username string
}

type Message struct {
	Type string     `json:"type"`
	Body string     `json:"body"`
	User string     `json:"user,omitempty"`
	Chat model.Chat `json:"chat,omitempty"`
}

// Each client is also responsible for listening to incoming messages from its connection.
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		m := &Message{}
		err = json.Unmarshal(p, m)
		if err != nil {
			log.Println("error while unmarshaling chat", err)
			continue
		}
		fmt.Println("host", c.Conn.RemoteAddr())

		if m.Type == "bootup" {
			// do mapping on bootup
			c.Username = m.User
			fmt.Println("client successfully mapped", &c, c, c.Username)
		} else {
			fmt.Println("received message", m.Type, m.Chat)
			chat := m.Chat
			chat.Timestamp = time.Now().Unix()

			// save in redis
			id, err := redisrepo.CreateChat(&chat)
			if err != nil {
				log.Println("error while saving chat in redis", err)
				return
			}

			chat.ID = id
			c.Pool.Broadcast <- &chat
			fmt.Printf("Message Received: %+v\n", &m)
		}
	}
}
