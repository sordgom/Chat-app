package chat

import (
	"fmt"
	"net/http"
	"strconv"
)

func serveWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	client := &Client{
		ID:   strconv.FormatInt(int64(len(pool.Clients)), 10),
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	go client.Read() // Keep listening for new messages from the client in a separate goroutine
}

func setupRoutes() {
	pool := NewPool()
	go pool.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func StartWebsocketServer() {
	setupRoutes()
	http.ListenAndServe(":8082", nil)
}
