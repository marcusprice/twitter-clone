package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/logger"
	"github.com/marcusprice/twitter-clone/internal/model"
	"github.com/marcusprice/twitter-clone/internal/permissions"
)

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tweetRotHost := os.Getenv("TWEETROT_HOST")
		tweetRotPort := os.Getenv("TWEETROT_PORT")

		w.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("http://%s:%s", tweetRotHost, tweetRotPort))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Max-Age", "600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateUser(user *controller.User, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value("requestID")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, Unauthorized, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := ParseJWT(tokenString)
		if err != nil || !token.Valid {
			logger.LogWarn("failed authenticating user")
			logger.LogInfo(
				fmt.Sprintf(
					"user authenitcation failed * requestID %v",
					requestID,
				),
			)
			http.Error(w, Unauthorized, http.StatusUnauthorized)
			return
		}

		claims, err := GetTokenClaims(token)
		if err != nil {
			logger.LogError("failed processing claims")
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			logger.LogWarn("failed processing sub claim")
			http.Error(w, Unauthorized, http.StatusUnauthorized)
			return
		}

		userID := int(sub)
		err = user.ByID(userID)
		if err != nil || (!user.IsActive && user.Role != permissions.SYSTEM_ROLE) {
			if err != nil && !errors.Is(err, model.UserNotFoundError{}) {
				http.Error(w, InternalServerError, http.StatusInternalServerError)
			} else {
				http.Error(w, Unauthorized, http.StatusUnauthorized)
			}

			return
		}

		logger.LogInfo(
			fmt.Sprintf(
				"user authenticated * userID: %d * requestID %v",
				userID,
				requestID,
			),
		)
		ctx := context.WithValue(
			r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func VerifyPostMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func VerifyGetMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AllowMethods(methods []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(methods, r.Method) {
			http.Error(w, MethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var requestCounter uint64

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := atomic.AddUint64(&requestCounter, 1)
		msg := fmt.Sprintf("%s %s requestID %d", r.Method, r.URL.Path, requestID)
		logger.LogInfo(msg)
		ctx := context.WithValue(r.Context(), "requestID", fmt.Sprintf("%d", requestID))
		r = r.WithContext(ctx)

		ww := &responseWriterWrapper{
			ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		msg = fmt.Sprintf(
			"%d * total time %v * requestID %d \n",
			ww.statusCode,
			time.Since(start),
			requestID)

		switch {
		case ww.statusCode >= 500:
			logger.LogError(msg)
		case ww.statusCode >= 400:
			logger.LogWarn(msg)
		default:
			logger.LogInfo(msg)
		}
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
