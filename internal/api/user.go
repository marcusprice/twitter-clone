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

type UserAPI struct {
	user *controller.User
}

func (userAPI UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInput dtypes.UserInput
	err := json.NewDecoder(r.Body).Decode(&userInput)

	if err != nil || !validUserFields(userInput, true) {
		http.Error(w, BadRequest, http.StatusBadRequest)
		return
	}

	user := userAPI.user
	user.Set(nil, userInput)

	err = user.Create(userInput.Password)
	if err != nil {
		var identifierError dtypes.IdentifierAlreadyExistsError

		if errors.As(err, &identifierError) {
			http.Error(w, Conflict, http.StatusConflict)
		} else if dbutils.IsConstraintError(err) {
			http.Error(w, BadRequest, http.StatusBadRequest)
		} else {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
		}

		return
	}

	payload := generateUserPayload(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func (userAPI UserAPI) Follow(w http.ResponseWriter, r *http.Request) {
	followerID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	followeeUsername := r.PathValue("username")
	if followeeUsername == "" {
		http.Error(w, BadRequest, http.StatusBadRequest)
	}

	follower := userAPI.user
	err := follower.ByID(followerID)
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPut {
		err = follower.Follow(followeeUsername)
	} else {
		err = follower.UnFollow(followeeUsername)
	}

	if err != nil {
		if errors.Is(err, model.UserNotFoundError{}) {
			http.Error(w, NotFound, http.StatusNotFound)
			return
		}

		http.Error(w, InternalServerError, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (userAPI UserAPI) Authenticate(w http.ResponseWriter, r *http.Request) {
	var userInput dtypes.UserInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	username := userInput.Username
	email := userInput.Email
	pwd := userInput.Password

	if err != nil || pwd == "" || (username == "" && email == "") {
		http.Error(w, BadRequest, http.StatusBadRequest)
		return
	}

	user := userAPI.user
	user.Set(nil, userInput)
	authenticated, err := user.AuthenticateAndSet(userInput.Password)
	if err != nil {
		var notFoundError model.UserNotFoundError
		if errors.As(err, &notFoundError) {
			http.Error(w, Unauthorized, http.StatusUnauthorized)
		} else {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
		}

		return
	}

	if !authenticated {
		http.Error(w, Unauthorized, http.StatusUnauthorized)
		return
	}

	if err := user.Login(); err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	payload := generateUserPayload(user)
	token, err := GenerateJWT(user.ID())
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func NewUserAPI(user *controller.User) *UserAPI {
	return &UserAPI{user}
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
