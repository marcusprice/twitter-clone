package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestCreateUser(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		handler := RegisterHandlers(db)
		newUserJson := `{
			"email": "estecat42069@yahoo.com",
			"username": "estecat",
			"displayName": "hungry boy",
			"password": "password",
			"firstName": "Esteban",
			"lastName": "Price"
		}`

		newUserRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/createUser",
			strings.NewReader(newUserJson))
		newUserResponse := httptest.NewRecorder()

		handler.ServeHTTP(newUserResponse, newUserRequest)

		if newUserResponse.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, newUserResponse.Code)
		}

		var payload UserPayload
		json.Unmarshal(newUserResponse.Body.Bytes(), &payload)

		if payload.Email != "estecat42069@yahoo.com" {
			t.Errorf("expected email %s, got %s", "estecat42069@yahoo.com", payload.Email)
		}

		if payload.Username != "estecat" {
			t.Errorf("expected username %s, got %s", "estecat", payload.Username)
		}

		if payload.DisplayName != "hungry boy" {
			t.Errorf("expected displayName %s, got %s", "hungry boy", payload.DisplayName)
		}

		if payload.FirstName != "Esteban" {
			t.Errorf("expected firstName %s, got %s", "Esteban", payload.FirstName)
		}

		if payload.LastName != "Price" {
			t.Errorf("expected lastName %s, got %s", "Price", payload.LastName)
		}

		newHumanJson := `{
			"email": "marcus@yodel.com",
			"username": "catdad42069",
			"displayName": "Dude where's my car",
			"password": "password",
			"firstName": "Marcus",
			"lastName": "Price"
		}`

		newHumanRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/createUser",
			strings.NewReader(newHumanJson))
		newHumanResponse := httptest.NewRecorder()

		handler.ServeHTTP(newHumanResponse, newHumanRequest)

		if newHumanResponse.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, newHumanResponse.Code)
		}

		var humanPayload UserPayload
		json.Unmarshal(newHumanResponse.Body.Bytes(), &humanPayload)

		if humanPayload.Email != "marcus@yodel.com" {
			t.Errorf("expected email %s, got %s", "marcus@yodel.com", humanPayload.Email)
		}

		if humanPayload.Username != "catdad42069" {
			t.Errorf("expected username %s, got %s", "catdad42069", humanPayload.Username)
		}

		if humanPayload.DisplayName != "Dude where's my car" {
			t.Errorf("expected displayName %s, got %s", "Dude where's my car", humanPayload.DisplayName)
		}

		if humanPayload.FirstName != "Marcus" {
			t.Errorf("expected firstName %s, got %s", "Esteban", humanPayload.FirstName)
		}

		if humanPayload.LastName != "Price" {
			t.Errorf("expected lastName %s, got %s", "Price", humanPayload.LastName)
		}
	})
}

func TestCreateUserAlreadyExists(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		handler := RegisterHandlers(db)
		user := controller.NewUserController(db)
		existingUser := dtypes.UserInput{
			Email:       "estecat42069@yahoo.com",
			Username:    "estecat",
			DisplayName: "estecat",
		}
		user.Set(nil, existingUser)
		user.Create("password")

		duplicateUserJson := `{
			"email": "estecat42069@yahoo.com",
			"username": "estecat",
			"displayName": "estecat",
			"password": "password"
		}`

		duplicateUserJsonSameEmail := `{
			"email": "estecat42069@yahoo.com",
			"username": "skateboarder_cat",
			"displayName": ":)",
			"password": "password"
		}`

		duplicateUserJsonSameUsername := `{
			"email": "skateboarder_cat@yahoo.com",
			"username": "estecat",
			"displayName": ":)",
			"password": "password"
		}`

		duplicateUserRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/createUser",
			strings.NewReader(duplicateUserJson))
		duplicateUserResponse := httptest.NewRecorder()

		duplicateEmailRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/createUser",
			strings.NewReader(duplicateUserJsonSameEmail))
		duplicateEmailResponse := httptest.NewRecorder()

		duplicateUsernameRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/createUser",
			strings.NewReader(duplicateUserJsonSameUsername))
		duplicateUsernameResponse := httptest.NewRecorder()

		handler.ServeHTTP(duplicateUserResponse, duplicateUserRequest)
		handler.ServeHTTP(duplicateEmailResponse, duplicateEmailRequest)
		handler.ServeHTTP(duplicateUsernameResponse, duplicateUsernameRequest)

		if duplicateUserResponse.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, duplicateUserResponse.Code)
		}

		if duplicateEmailResponse.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, duplicateEmailResponse.Code)
		}

		if duplicateUsernameResponse.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, duplicateUsernameResponse.Code)
		}
	})
}

func TestCreateUserMissingRequiredFields(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		handler := RegisterHandlers(db)

		missingUsername := `{
			"email": "estecat42069@yahoo.com",
			"displayName": "estecat",
			"password": "password"
		}`

		missingEmail := `{
			"displayName": "estecat",
			"username": "esteban",
			"password": "password"
		}`

		missingDisplayName := `{
			"email": "estecat42069@yahoo.com",
			"username": "esteban",
			"password": "password"
		}`

		missingPassword := `{
			"username": "esteban",
			"email": "estecat42069@yahoo.com",
			"displayName": "estecat",
		}`

		missingUsernameRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/createUser",
			strings.NewReader(missingUsername))
		missingUsernameResponse := httptest.NewRecorder()

		missingEmailRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/createUser",
			strings.NewReader(missingEmail))
		missingEmailResponse := httptest.NewRecorder()

		missingDisplayNameRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/createUser",
			strings.NewReader(missingDisplayName))
		missingDisplayNameResponse := httptest.NewRecorder()

		missingPasswordRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/createUser",
			strings.NewReader(missingPassword))
		missingPasswordResponse := httptest.NewRecorder()

		handler.ServeHTTP(missingEmailResponse, missingEmailRequest)
		handler.ServeHTTP(missingUsernameResponse, missingUsernameRequest)
		handler.ServeHTTP(missingDisplayNameResponse, missingDisplayNameRequest)
		handler.ServeHTTP(missingPasswordResponse, missingPasswordRequest)

		if missingEmailResponse.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, missingPasswordResponse.Code)
		}

		if missingUsernameResponse.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, missingPasswordResponse.Code)
		}

		if missingDisplayNameResponse.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, missingPasswordResponse.Code)
		}

		if missingPasswordResponse.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, missingPasswordResponse.Code)
		}
	})
}

func TestCreateUserMalformedJSON(t *testing.T) {
	UserAPI := UserAPI{}

	malformedJSON := "alkj}"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/createUser", strings.NewReader(malformedJSON))
	res := httptest.NewRecorder()
	UserAPI.CreateUser(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestCreateUserWrongMethod(t *testing.T) {
	UserAPI := UserAPI{}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/createUser", nil)
	getRes := httptest.NewRecorder()
	UserAPI.CreateUser(getRes, getReq)

	if getRes.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, getRes.Code)
	}

	putReq := httptest.NewRequest(http.MethodPut, "/api/v1/createUser", nil)
	putRes := httptest.NewRecorder()
	UserAPI.CreateUser(putRes, putReq)

	if putRes.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, putRes.Code)
	}

	patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/createUser", nil)
	patchRes := httptest.NewRecorder()
	UserAPI.CreateUser(patchRes, patchReq)

	if patchRes.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, patchRes.Code)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/createUser", nil)
	deleteRes := httptest.NewRecorder()
	UserAPI.CreateUser(deleteRes, deleteReq)

	if deleteRes.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, deleteRes.Code)
	}

	headReq := httptest.NewRequest(http.MethodHead, "/api/v1/createUser", nil)
	headRes := httptest.NewRecorder()
	UserAPI.CreateUser(headRes, headReq)

	if headRes.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, headRes.Code)
	}

	optionReq := httptest.NewRequest(http.MethodOptions, "/api/v1/createUser", nil)
	optionRes := httptest.NewRecorder()
	UserAPI.CreateUser(optionRes, optionReq)

	if optionRes.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, optionRes.Code)
	}

	traceReq := httptest.NewRequest(http.MethodTrace, "/api/v1/createUser", nil)
	traceRes := httptest.NewRecorder()
	UserAPI.CreateUser(traceRes, traceReq)

	if traceRes.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, traceRes.Code)
	}

	connectReq := httptest.NewRequest(http.MethodConnect, "/api/v1/createUser", nil)
	connectRes := httptest.NewRecorder()
	UserAPI.CreateUser(connectRes, connectReq)

	if connectRes.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, connectRes.Code)
	}
}
