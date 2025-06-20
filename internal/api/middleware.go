package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/logger"
	"github.com/marcusprice/twitter-clone/internal/model"
)

func ValidateUser(user *controller.User, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, Unauthorized, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := ParseJWT(tokenString)
		if err != nil || !token.Valid {
			logger.LogWarn("failed authenticating user")
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
		if err != nil || !user.IsActive {
			var userNotFoundError model.UserNotFoundError
			if err != nil && !errors.As(err, &userNotFoundError) {
				logger.LogError("user not found?")
				http.Error(w, InternalServerError, http.StatusInternalServerError)
			} else {
				http.Error(w, Unauthorized, http.StatusUnauthorized)
			}

			return
		}

		logger.LogInfo(fmt.Sprintf("user authenticated, request userID: %d", userID))
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

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriterWrapper{
			ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		msg := fmt.Sprintf(
			"%s %s | status=%d | %v\n",
			r.Method, r.URL.Path, ww.statusCode, time.Since(start),
		)

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
