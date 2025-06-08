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
		tu := testutil.NewTestUtil(t)
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

		var payload UserPayload
		json.Unmarshal(newUserResponse.Body.Bytes(), &payload)

		tu.AssertEqual(http.StatusOK, newUserResponse.Code)
		tu.AssertEqual("estecat42069@yahoo.com", payload.Email)
		tu.AssertEqual("estecat", payload.Username)
		tu.AssertEqual("hungry boy", payload.DisplayName)
		tu.AssertEqual("Esteban", payload.FirstName)
		tu.AssertEqual("Price", payload.LastName)

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

		tu.AssertEqual(http.StatusOK, newHumanResponse.Code)
		tu.AssertEqual("marcus@yodel.com", humanPayload.Email)
		tu.AssertEqual("catdad42069", humanPayload.Username)
		tu.AssertEqual("Dude where's my car", humanPayload.DisplayName)
		tu.AssertEqual("Marcus", humanPayload.FirstName)
		tu.AssertEqual("Price", humanPayload.LastName)
	})
}

func TestCreateUserAlreadyExists(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
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

		tu.AssertEqual(http.StatusConflict, duplicateUserResponse.Code)
		tu.AssertEqual(http.StatusConflict, duplicateEmailResponse.Code)
		tu.AssertEqual(http.StatusConflict, duplicateUsernameResponse.Code)
	})
}

func TestCreateUserMissingRequiredFields(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
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

		tu.AssertEqual(http.StatusBadRequest, missingEmailResponse.Code)
		tu.AssertEqual(http.StatusBadRequest, missingUsernameResponse.Code)
		tu.AssertEqual(http.StatusBadRequest, missingDisplayNameResponse.Code)
		tu.AssertEqual(http.StatusBadRequest, missingPasswordResponse.Code)
	})
}

func TestCreateUserMalformedJSON(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	UserAPI := UserAPI{}

	malformedJSON := "alkj}"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/createUser", strings.NewReader(malformedJSON))
	res := httptest.NewRecorder()
	UserAPI.CreateUser(res, req)

	tu.AssertEqual(http.StatusBadRequest, res.Code)
}

func TestCreateUserWrongMethod(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	UserAPI := UserAPI{}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/createUser", nil)
	getRes := httptest.NewRecorder()
	putReq := httptest.NewRequest(http.MethodPut, "/api/v1/createUser", nil)
	putRes := httptest.NewRecorder()
	patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/createUser", nil)
	patchRes := httptest.NewRecorder()
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/createUser", nil)
	deleteRes := httptest.NewRecorder()
	headReq := httptest.NewRequest(http.MethodHead, "/api/v1/createUser", nil)
	headRes := httptest.NewRecorder()
	optionReq := httptest.NewRequest(http.MethodOptions, "/api/v1/createUser", nil)
	optionRes := httptest.NewRecorder()
	traceReq := httptest.NewRequest(http.MethodTrace, "/api/v1/createUser", nil)
	traceRes := httptest.NewRecorder()
	connectReq := httptest.NewRequest(http.MethodConnect, "/api/v1/createUser", nil)
	connectRes := httptest.NewRecorder()

	UserAPI.CreateUser(getRes, getReq)
	UserAPI.CreateUser(putRes, putReq)
	UserAPI.CreateUser(patchRes, patchReq)
	UserAPI.CreateUser(deleteRes, deleteReq)
	UserAPI.CreateUser(headRes, headReq)
	UserAPI.CreateUser(optionRes, optionReq)
	UserAPI.CreateUser(traceRes, traceReq)
	UserAPI.CreateUser(connectRes, connectReq)

	tu.AssertEqual(http.StatusMethodNotAllowed, getRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, putRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, patchRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, deleteRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, headRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, optionRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, traceRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, connectRes.Code)
}
