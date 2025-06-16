package api

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/model"
)

func Authenticate(user *controller.User, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, Unauthorized, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := ParseJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, Unauthorized, http.StatusUnauthorized)
			return
		}

		claims, err := GetTokenClaims(token)
		if err != nil {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, Unauthorized, http.StatusUnauthorized)
			return
		}

		userID := int(sub)
		err = user.ByID(userID)
		if err != nil || !user.IsActive {
			var userNotFoundError model.UserNotFoundError
			if err != nil && !errors.As(err, &userNotFoundError) {
				http.Error(w, InternalServerError, http.StatusInternalServerError)
			} else {
				http.Error(w, Unauthorized, http.StatusUnauthorized)
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
