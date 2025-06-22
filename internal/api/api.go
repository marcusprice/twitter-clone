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
	comment := controller.NewCommentController(db)
	userAPI := NewUserAPI(user)
	postAPI := NewPostAPI(post)
	commentAPI := NewCommentAPI(comment)
	timelineAPI := NewTimelineAPI(db)

	mux := http.NewServeMux()

	if util.InDevContext() {
		projectRoot, err := util.ProjectRoot()
		if err != nil {
			panic(err)
		}

		fs := http.FileServer(http.Dir(projectRoot + "/static/swagger-ui/"))

		mux.Handle(
			"/docs/",
			Logger(
				http.StripPrefix("/docs/", fs)))

		mux.Handle(
			"/swagger.yaml",
			Logger(
				http.FileServer(http.Dir("."))))
	}

	mux.Handle(
		"/api/v1/timeline",
		Logger(
			VerifyGetMethod(
				ValidateUser(
					user,
					http.HandlerFunc(timelineAPI.Get)))),
	)

	mux.Handle(
		"/api/v1/user/create",
		Logger(
			VerifyPostMethod(
				http.HandlerFunc(userAPI.CreateUser))),
	)

	mux.Handle(
		"/api/v1/user/authenticate",
		Logger(
			VerifyPostMethod(
				http.HandlerFunc(userAPI.Authenticate))),
	)

	mux.Handle(
		"/api/v1/user/{username}/follow",
		Logger(
			AllowMethods(
				[]string{http.MethodPut, http.MethodDelete},
				ValidateUser(
					user,
					http.HandlerFunc(userAPI.Follow)))),
	)

	mux.Handle(
		"/api/v1/post/create",
		Logger(
			VerifyPostMethod(
				ValidateUser(
					user,
					http.HandlerFunc(postAPI.CreatePost)))),
	)

	mux.Handle(
		"/api/v1/post/{id}/like",
		Logger(
			AllowMethods(
				[]string{http.MethodPut, http.MethodDelete},
				ValidateUser(
					user,
					http.HandlerFunc(postAPI.Like)))),
	)

	mux.Handle(
		"/api/v1/post/{id}/retweet",
		Logger(
			AllowMethods(
				[]string{http.MethodPut, http.MethodDelete},
				ValidateUser(
					user,
					http.HandlerFunc(postAPI.Retweet)))),
	)

	mux.Handle(
		"/api/v1/post/{id}/bookmark",
		Logger(
			AllowMethods(
				[]string{http.MethodPut, http.MethodDelete},
				ValidateUser(
					user,
					http.HandlerFunc(postAPI.Bookmark)))),
	)

	mux.Handle(
		"/api/v1/comment/create",
		Logger(
			VerifyPostMethod(
				ValidateUser(
					user,
					http.HandlerFunc(commentAPI.Create)))),
	)

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
