package login

import (
	"fmt"
	"net/http"

	auth "chat-go/pkg/AuthMiddleware"
	"chat-go/pkg/redisrepo"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func StartHTTPServer() {
	// initialise redis
	redisClient := redisrepo.InitialiseRedis()
	defer redisClient.Close()

	// create indexes
	redisrepo.CreateFetchChatBetweenIndex()

	r := mux.NewRouter()
	r.Use(auth.AuthMiddleware)

	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	}).Methods(http.MethodGet)

	r.HandleFunc("/register", registerHandler).Methods(http.MethodPost)
	r.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	r.HandleFunc("/verify-contact", verifyContactHandler).Methods(http.MethodPost)
	r.HandleFunc("/chat-history", chatHistoryHandler).Methods(http.MethodGet)
	r.HandleFunc("/contact-list", contactListHandler).Methods(http.MethodGet)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // change to your frontend url
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true, // this allows cookies to be sent
	})
	// Use default options
	handler := c.Handler(r)
	http.ListenAndServe(":8080", handler)
}
