package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/model"
)

type UserPayload struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
}

type UserAPI struct {
	user *controller.User
}

func (userAPI UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInput dtypes.UserInput

	if !validRequestMethod(http.MethodPost, r.Method) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&userInput)

	if err != nil || !validUserFields(userInput, true) {
		statusMessage := http.StatusText(http.StatusBadRequest)
		http.Error(w, statusMessage, http.StatusBadRequest)
		return
	}

	user := userAPI.user
	user.Set(nil, userInput)

	err = user.Create(userInput.Password)
	if err != nil {
		var identifierError dtypes.IdentifierAlreadyExistsError

		if errors.As(err, &identifierError) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else if dbutils.IsConstraintError(err) {
			statusText := http.StatusText(http.StatusBadRequest)
			http.Error(w, statusText, http.StatusBadRequest)
		} else {
			statusText := http.StatusText(http.StatusInternalServerError)
			http.Error(w, statusText, http.StatusInternalServerError)
		}

		return
	}

	payload := generateUserPayload(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func (userAPI UserAPI) Authenticate(w http.ResponseWriter, r *http.Request) {
	if !validRequestMethod(http.MethodPost, r.Method) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userInput dtypes.UserInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	username := userInput.Username
	email := userInput.Email
	pwd := userInput.Password

	if err != nil || pwd == "" || (username == "" && email == "") {
		statusMessage := http.StatusText(http.StatusBadRequest)
		http.Error(w, statusMessage, http.StatusBadRequest)
		return
	}

	user := userAPI.user
	user.Set(nil, userInput)
	authenticated, err := user.AuthenticateAndSet(userInput.Password)
	if err != nil {
		var notFoundError model.UserNotFoundError
		if errors.As(err, &notFoundError) {
			statusText := http.StatusText(http.StatusUnauthorized)
			http.Error(w, statusText, http.StatusUnauthorized)
		} else {
			statusText := http.StatusText(http.StatusInternalServerError)
			http.Error(w, statusText, http.StatusInternalServerError)
		}

		return
	}

	if !authenticated {
		statusText := http.StatusText(http.StatusUnauthorized)
		http.Error(w, statusText, http.StatusUnauthorized)
		return
	}

	if err := user.SetLastLogin(); err != nil {
		statusText := http.StatusText(http.StatusInternalServerError)
		http.Error(w, statusText, http.StatusInternalServerError)
		return
	}

	payload := generateUserPayload(user)
	token, err := GenerateJWT(user.ID())
	if err != nil {
		statusText := http.StatusText(http.StatusInternalServerError)
		http.Error(w, statusText, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func NewUserAPI(user *controller.User) UserAPI {
	return UserAPI{user}
}

func validUserFields(userInput dtypes.UserInput, pwdRequired bool) bool {
	if pwdRequired && userInput.Password == "" {
		return false
	}

	if userInput.Username == "" {
		return false
	}

	if userInput.Email == "" {
		return false
	}

	if userInput.DisplayName == "" {
		return false
	}

	return true
}

func generateUserPayload(user *controller.User) UserPayload {
	return UserPayload{
		user.Email, user.Username, user.FirstName,
		user.LastName, user.DisplayName}
}
