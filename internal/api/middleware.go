package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/model"
)

func Authenticate(user *controller.User, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			statusText := http.StatusText(http.StatusUnauthorized)
			http.Error(w, statusText, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := ParseJWT(tokenString)
		if err != nil || !token.Valid {
			statusText := http.StatusText(http.StatusUnauthorized)
			http.Error(w, statusText, http.StatusUnauthorized)
			return
		}

		claims, err := GetTokenClaims(token)
		if err != nil {
			statusText := http.StatusText(http.StatusInternalServerError)
			http.Error(w, statusText, http.StatusInternalServerError)
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			statusText := http.StatusText(http.StatusUnauthorized)
			http.Error(w, statusText, http.StatusUnauthorized)
			return
		}

		userID := int(sub)
		err = user.ByID(userID)
		if err != nil || !user.IsActive {
			var userNotFoundError model.UserNotFoundError
			if err != nil && !errors.As(err, &userNotFoundError) {
				statusText := http.StatusText(http.StatusInternalServerError)
				http.Error(w, statusText, http.StatusInternalServerError)
			} else {
				statusText := http.StatusText(http.StatusUnauthorized)
				http.Error(w, statusText, http.StatusUnauthorized)
			}

			return
		}

		ctx := context.WithValue(
			r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func VerifyPostRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			statusText := http.StatusText(http.StatusMethodNotAllowed)
			http.Error(w, statusText, http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}
