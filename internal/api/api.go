package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/util"
)

func RegisterHandlers(db *sql.DB) http.Handler {
	if db == nil {
		panic("db conn cannot be nil")
	}

	user := controller.NewUserController(db)
	post := controller.NewPostController(db)
	userAPI := NewUserAPI(user)
	postAPI := NewPostAPI(post)

	mux := http.NewServeMux()

	mux.Handle(
		"/api/v1/createUser",
		VerifyPostRequest(
			http.HandlerFunc(userAPI.CreateUser)),
	)

	mux.Handle(
		"/api/v1/authenticateUser",
		VerifyPostRequest(
			http.HandlerFunc(userAPI.Authenticate)),
	)

	mux.Handle(
		"/api/v1/createPost",
		VerifyPostRequest(
			Authenticate(
				user,
				http.HandlerFunc(postAPI.CreatePost))))

	return mux
}

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
	})

	secretKey := os.Getenv("JWT_KEY")
	if secretKey == "" {
		if util.InDevContext() {
			panic("JWT_KEY environment variable required")
		} else {
			return "", errors.New("missing jwt key")
		}
	}

	return token.SignedString([]byte(secretKey))
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return &jwt.Token{}, errors.New("Missing token string")
	}

	secretKey := os.Getenv("JWT_KEY")
	if secretKey == "" {
		if util.InDevContext() {
			panic("JWT_KEY environment variable required")
		} else {
			return &jwt.Token{}, errors.New("missing jwt key")
		}
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
}

func GetTokenClaims(token *jwt.Token) (jwt.MapClaims, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return jwt.MapClaims{}, errors.New("No claims")
	}
}
