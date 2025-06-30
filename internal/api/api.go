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

	mux.Handle(
		"/api/v1/timeline",
		VerifyGetMethod(
			ValidateUser(
				user,
				http.HandlerFunc(timelineAPI.Get))),
	)

	mux.Handle(
		"/api/v1/user",
		ValidateUser(
			user,
			http.HandlerFunc(userAPI.Get),
		),
	)

	mux.Handle(
		"/api/v1/user/create",
		VerifyPostMethod(
			http.HandlerFunc(userAPI.Create)),
	)

	mux.Handle(
		"/api/v1/user/authenticate",
		VerifyPostMethod(
			http.HandlerFunc(userAPI.Authenticate)),
	)

	mux.Handle(
		"/api/v1/user/bookmarks",
		VerifyGetMethod(
			ValidateUser(
				user,
				http.HandlerFunc(userAPI.GetBookmarks))),
	)

	mux.Handle(
		"/api/v1/user/{username}/follow",
		AllowMethods(
			[]string{http.MethodPut, http.MethodDelete},
			ValidateUser(
				user,
				http.HandlerFunc(userAPI.Follow))),
	)

	mux.Handle(
		"/api/v1/post/{postID}",
		VerifyGetMethod(
			ValidateUser(
				user,
				http.HandlerFunc(postAPI.Get))),
	)

	mux.Handle(
		"/api/v1/post/create",
		VerifyPostMethod(
			ValidateUser(
				user,
				http.HandlerFunc(postAPI.Create))),
	)

	mux.Handle(
		"/api/v1/post/{id}/like",
		AllowMethods(
			[]string{http.MethodPut, http.MethodDelete},
			ValidateUser(
				user,
				http.HandlerFunc(postAPI.Like))),
	)

	mux.Handle(
		"/api/v1/post/{id}/retweet",
		AllowMethods(
			[]string{http.MethodPut, http.MethodDelete},
			ValidateUser(
				user,
				http.HandlerFunc(postAPI.Retweet))),
	)

	mux.Handle(
		"/api/v1/post/{id}/bookmark",
		AllowMethods(
			[]string{http.MethodPut, http.MethodDelete},
			ValidateUser(
				user,
				http.HandlerFunc(postAPI.Bookmark))),
	)

	mux.Handle(
		"/api/v1/comment/create",
		VerifyPostMethod(
			ValidateUser(
				user,
				http.HandlerFunc(commentAPI.Create))),
	)

	projectRoot, err := util.ProjectRoot()
	uploadFileServer := http.FileServer(http.Dir(projectRoot + "/upload"))

	mux.Handle(
		"/uploads/",
		VerifyGetMethod(
			http.StripPrefix(UPLOADS_PREFIX, uploadFileServer)),
	)

	if util.InDevContext() {
		if err != nil {
			panic(err)
		}

		swaggerFS := http.FileServer(http.Dir(projectRoot + "/static/swagger-ui/"))

		mux.Handle(
			"/docs/",
			http.StripPrefix("/docs/", swaggerFS))

		mux.Handle(
			"/swagger.yaml",
			http.FileServer(http.Dir(".")))
	}

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

	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
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
