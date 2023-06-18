package login

import (
	"net/http"

	"chat-go/pkg/login/redisrepo"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func StartHTTPServer() {
	// initialise redis
	redisClient := redisrepo.InitialiseRedis()
	defer redisClient.Close()

	// create indexes
	redisrepo.CreateFetchChatBetweenIndex()

	//logger
	log, _ := zap.NewProduction()
	defer log.Sync()

	r := mux.NewRouter()

	suc := NewSignupController(log)
	sic := NewSigninController(log)

	r.HandleFunc("/signup", suc.SignupHandler).Methods(http.MethodPost)
	r.HandleFunc("/signin", sic.SigninHandler).Methods(http.MethodPost)

	r.HandleFunc("/register", registerHandler).Methods(http.MethodPost)
	r.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	r.HandleFunc("/verify-contact", verifyContactHandler).Methods(http.MethodPost)
	r.HandleFunc("/chat-history", chatHistoryHandler).Methods(http.MethodGet)
	r.HandleFunc("/contact-list", contactListHandler).Methods(http.MethodGet)
	r.Path("/metrics").Handler(promhttp.Handler())

	//Middleware
	//tm := middleware.NewTokenMiddleware(log)
	// r.Use(tm.TokenValidationMiddleware)

	cors := gohandlers.CORS(
		gohandlers.AllowedOrigins([]string{"http://localhost:3000"}), // your front end url
		gohandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		gohandlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
		gohandlers.AllowCredentials(),
	)
	// Use default options
	handler := cors(r)
	http.ListenAndServe(":8080", handler)
}
