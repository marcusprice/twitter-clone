package api

import (
	"database/sql"
	"net/http"

	"github.com/marcusprice/twitter-clone/internal/controller"
)

func RegisterHandlers(dbConn *sql.DB) http.Handler {
	if dbConn == nil {
		panic("db conn cannot be nil")
	}

	user := controller.NewUserController(dbConn)
	userAPI := NewUserAPI(user)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/createUser", userAPI.CreateUser)
	mux.HandleFunc("/api/v1/authenticateUser", userAPI.Authenticate)

	return mux
}

func validRequestMethod(expectedMethod, actualMethod string) bool {
	return expectedMethod == actualMethod
}
