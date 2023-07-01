package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"chat-go/pkg/login/jwt"
	"chat-go/pkg/login/redisrepo"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	signinRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signin_total",
		Help: "Total number of signup requests",
	})
	signinSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signin_success",
		Help: "Successful signup requests",
	})
	signinFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signin_fail",
		Help: "Failed signup requests",
	})
	signinError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signin_error",
		Help: "Erroneous signup requests",
	})
)

var (
	singupRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signup_total",
		Help: "Total number of signup requests",
	})
	signupSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signup_success",
		Help: "Successful signup requests",
	})
	signupFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signup_fail",
		Help: "Failed signup requests",
	})
)

// SignupController is the Signup route handler
type SignupController struct {
	logger            *zap.Logger
	promSignupTotal   prometheus.Counter
	promSignupSuccess prometheus.Counter
	promSignupFail    prometheus.Counter
}

// NewSignupController returns a frsh Signup controller
func NewSignupController(logger *zap.Logger) *SignupController {
	return &SignupController{
		logger:            logger,
		promSignupTotal:   singupRequests,
		promSignupSuccess: signupSuccess,
		promSignupFail:    signupFail,
	}
}

// SigninController is the Signin route handler
type SigninController struct {
	logger            *zap.Logger
	promSigninTotal   prometheus.Counter
	promSigninSuccess prometheus.Counter
	promSigninFail    prometheus.Counter
	promSigninError   prometheus.Counter
}

// NewSigninController returns a frsh Signin controller
func NewSigninController(logger *zap.Logger) *SigninController {
	return &SigninController{
		logger:            logger,
		promSigninTotal:   signinRequests,
		promSigninSuccess: signinSuccess,
		promSigninFail:    signinFail,
		promSigninError:   signinError,
	}
}

type userReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

type response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int         `json:"total,omitempty"`
}

// we need this function to be private
func getSignedToken() (string, error) {
	claimsMap := jwt.ClaimsMap{
		Aud: "chat.xyz",
		Iss: "chat_app",
		Exp: fmt.Sprint(time.Now().Add(time.Hour * 3).Unix()),
	}

	secret := jwt.GetSecret()
	if secret == "" {
		return "", errors.New("empty JWT secret")
	}

	header := "HS256"
	tokenString, err := jwt.GenerateToken(header, claimsMap, secret)
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

// Handlers
func verifyContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoidng request object", http.StatusBadRequest)
		return
	}

	res := verifyContact(u.Username)
	json.NewEncoder(w).Encode(res)
}

func chatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// user1 user2
	u1 := r.URL.Query().Get("u1")
	u2 := r.URL.Query().Get("u2")

	// chat between timerange fromTS toTS
	// where TS is timestamp
	// 0 to positive infinity
	fromTS, toTS := "0", "+inf"

	if r.URL.Query().Get("from-ts") != "" && r.URL.Query().Get("to-ts") != "" {
		fromTS = r.URL.Query().Get("from-ts")
		toTS = r.URL.Query().Get("to-ts")
	}

	res := chatHistory(u1, u2, fromTS, toTS)
	json.NewEncoder(w).Encode(res)
}

func contactListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := r.URL.Query().Get("username")

	res := contactList(u)
	json.NewEncoder(w).Encode(res)
}

// This will be supplied to the MUX router. It will be called when signin request is sent
// if user not found or not validates, returns the Unauthorized error
// if found, returns the JWT back. [How to return this?]
func (ctrl *SigninController) SigninHandler(rw http.ResponseWriter, r *http.Request) {
	// increment total singin requests
	ctrl.promSigninTotal.Inc()
	rw.Header().Set("Content-Type", "application/json")

	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Error decoding request object",
		})
		return
	}

	res := login(u, rw, r)

	tokenString, err := getSignedToken()
	if err != nil {
		ctrl.logger.Error("unable to sign the token", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"status":  "error",
			"message": "Internal server error: unable to sign the token",
		})
		ctrl.promSigninError.Inc()
		return
	}

	ctrl.logger.Info("Token sign", zap.String("token", tokenString))
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Token successfully signed",
		"token":   tokenString,
		"data":    res,
	})
	ctrl.promSigninSuccess.Inc()
}

// adds the user to the database of users
func (ctrl *SignupController) SignupHandler(rw http.ResponseWriter, r *http.Request) {
	// we increment the signup request counter
	ctrl.promSignupTotal.Inc()
	rw.Header().Set("Content-Type", "application/json")

	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(rw, "error decoidng request object", http.StatusBadRequest)
		return
	}

	res := register(u)
	json.NewEncoder(rw).Encode(res)

	if r.Header["Username"] != nil {
		ctrl.logger.Info("User created", zap.String("username", r.Header["Username"][0]))
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("User Created"))
	// this will mean the request was successfully added
	ctrl.promSignupSuccess.Inc()
}

// Logic functions
func register(u *userReq) *response {
	// check if username in userset
	// return error if exist
	// create new user
	// create response for error
	res := &response{Status: true}

	status := redisrepo.IsUserExist(u.Username)
	if status {
		res.Status = false
		res.Message = "username already taken. try something else."
		return res
	}

	err := redisrepo.RegisterNewUser(u.Username, u.Password)
	if err != nil {
		res.Status = false
		res.Message = "something went wrong while registering the user. please try again after sometime."
		return res
	}

	return res
}

func login(u *userReq, w http.ResponseWriter, r *http.Request) *response {
	res := &response{Status: true}

	err := redisrepo.IsUserAuthentic(u.Username, u.Password)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return res
	}
	return res
}

func verifyContact(username string) *response {
	// if invalid username and password return error
	// if valid user create new session
	res := &response{Status: true}

	status := redisrepo.IsUserExist(username)
	if !status {
		res.Status = false
		res.Message = "invalid username"
	}

	return res
}

func chatHistory(username1, username2, fromTS, toTS string) *response {
	// if invalid usernames return error
	// if valid users fetch chats
	res := &response{}

	fmt.Println(username1, username2)
	// check if user exists
	if !redisrepo.IsUserExist(username1) || !redisrepo.IsUserExist(username2) {
		res.Message = "incorrect username"
		return res
	}

	chats, err := redisrepo.FetchChatBetween(username1, username2, fromTS, toTS)
	if err != nil {
		log.Println("error in fetch chat between", err)
		res.Message = "unable to fetch chat history. please try again later."
		return res
	}

	res.Status = true
	res.Data = chats
	res.Total = len(chats)
	return res
}

func contactList(username string) *response {
	// if invalid username return error
	// if valid users fetch chats
	res := &response{}

	// check if user exists
	if !redisrepo.IsUserExist(username) {
		res.Message = "incorrect username"
		return res
	}

	contactList, err := redisrepo.FetchContactList(username)
	if err != nil {
		log.Println("error in fetch contact list of username: ", username, err)
		res.Message = "unable to fetch contact list. please try again later."
		return res
	}

	res.Status = true
	res.Data = contactList
	res.Total = len(contactList)
	return res
}
