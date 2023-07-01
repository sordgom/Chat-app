package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"chat-go/pkg/login/jwt"

	"go.uber.org/zap"
)

// TokenMiddleware is the token validation route handler
type TokenMiddleware struct {
	logger *zap.Logger
}

// NewTokenMiddleware returns a frsh Token controller
func NewTokenMiddleware(logger *zap.Logger) *TokenMiddleware {
	return &TokenMiddleware{
		logger: logger,
	}
}

func (ctrl *TokenMiddleware) TokenValidation(rw http.ResponseWriter, r *http.Request) {
	// check if token is present
	authHeader, ok := r.Header["Authorization"]
	if !ok || len(authHeader) == 0 {
		ctrl.logger.Warn("Token was not found in the header")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("false"))
		return
	}
	headerParts := strings.Split(authHeader[0], " ")
	if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
		ctrl.logger.Warn("Token was not found in the header")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("false"))
		return
	}

	token := headerParts[1]

	secret := jwt.GetSecret()
	if secret == "" {
		ctrl.logger.Error("Empty JWT secret")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("false"))
		return
	}

	err := jwt.ValidateToken(token, secret)
	if err != nil {
		errInString := fmt.Sprint(err)
		ctrl.logger.Error(errInString, zap.String("token", token))
		if errInString == jwt.CORRUPT_TOKEN || errInString == jwt.INVALID_TOKEN || errInString == jwt.EXPIRED_TOKEN {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("false"))
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Internal Server Error"))
		}
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("true"))
}

// Middleware itself returns a function that is a Handler. it is executed for each request.
// We want all our routes for REST to be authenticated. So, we validate the token
func (ctrl *TokenMiddleware) TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// check if token is present
		authHeader, ok := r.Header["Authorization"]
		if !ok || len(authHeader) == 0 {
			ctrl.logger.Warn("Token was not found in the header")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Token Missing"))
			return
		}
		headerParts := strings.Split(authHeader[0], " ")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			ctrl.logger.Warn("Token was not found in the header")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Token Missing"))
			return
		}

		token := headerParts[1]

		secret := jwt.GetSecret()
		if secret == "" {
			ctrl.logger.Error("Empty JWT secret")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Internal Server Error"))
			return
		}

		err := jwt.ValidateToken(token, secret)
		if err != nil {
			errInString := fmt.Sprint(err)
			ctrl.logger.Error(errInString, zap.String("token", token))
			if errInString == jwt.CORRUPT_TOKEN || errInString == jwt.INVALID_TOKEN || errInString == jwt.EXPIRED_TOKEN {
				rw.WriteHeader(http.StatusUnauthorized)
			} else {
				rw.WriteHeader(http.StatusInternalServerError)
			}
			rw.Write([]byte(errInString))
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Authorized Token"))

		// this calls the next function. If not included, the router wont entertain any requests
		next.ServeHTTP(rw, r)
	})
}
