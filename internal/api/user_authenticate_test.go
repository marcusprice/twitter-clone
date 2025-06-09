package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestAuthenticateUser(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)
		user := controller.NewUserController(db)
		userInput := dtypes.UserInput{
			Username:    "esteban",
			Email:       "estecat42069@yahoo.com",
			Password:    "password",
			DisplayName: "yodel",
		}
		user.Set(nil, userInput)
		user.Create("password")

		authJson := `{
			"username": "esteban",
			"email": "estecat42069@yahoo.com",
			"password": "password"
		}`

		authUsernameOnlyJson := `{
			"username": "esteban",
			"password": "password"
		}`

		authEmailOnlyJson := `{
			"username": "esteban",
			"password": "password"
		}`

		authReq := httptest.NewRequest(http.MethodPost,
			"/api/v1/authenticateUser",
			strings.NewReader(authJson))
		authRes := httptest.NewRecorder()
		authWithUsernameReq := httptest.NewRequest(http.MethodPost,
			"/api/v1/authenticateUser",
			strings.NewReader(authUsernameOnlyJson))
		authWithUsernameRes := httptest.NewRecorder()
		authWithEmailReq := httptest.NewRequest(http.MethodPost,
			"/api/v1/authenticateUser",
			strings.NewReader(authEmailOnlyJson))
		authWithEmailRes := httptest.NewRecorder()

		beforeRequest := time.Now().UTC().Add(-1 * time.Minute)
		handler.ServeHTTP(authRes, authReq)
		afterRequest := time.Now().UTC().Add(time.Minute)

		handler.ServeHTTP(authWithUsernameRes, authWithUsernameReq)
		handler.ServeHTTP(authWithEmailRes, authWithEmailReq)

		var userPayload UserPayload
		json.Unmarshal(authRes.Body.Bytes(), &userPayload)

		tu.AssertEqual(http.StatusOK, authRes.Code)
		tu.AssertEqual(http.StatusOK, authWithUsernameRes.Code)
		tu.AssertEqual(http.StatusOK, authWithEmailRes.Code)
		tu.AssertEqual(userPayload.Username, "esteban")
		tu.AssertEqual(userPayload.Email, "estecat42069@yahoo.com")
		tu.AssertEqual(userPayload.DisplayName, "yodel")

		authHeader := authRes.Header().Get("Authorization")
		tokenString := strings.Split(authHeader, " ")[1]

		token, err := ParseJWT(tokenString)
		if err != nil {
			panic(err)
		}

		claims, err := GetTokenClaims(token)
		if err != nil {
			panic(err)
		}

		userID := int(claims["sub"].(float64))
		user.ByID(userID)

		tu.AssertTrue(user.LastLogin.After(beforeRequest))
		tu.AssertTrue(user.LastLogin.Before(afterRequest))
		tu.AssertTrue(token.Valid)
		tu.AssertEqual(1, userID)
	})
}

func TestAuthenticateUserWrongPassword(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)
		user := controller.NewUserController(db)
		userInput := dtypes.UserInput{
			Username:    "esteban",
			Email:       "estecat42069@yahoo.com",
			Password:    "password",
			DisplayName: "yodel",
		}
		user.Set(nil, userInput)
		user.Create("password")

		authJson := `{
			"username": "esteban",
			"email": "estecat42069@yahoo.com",
			"password": "wrong_password"
		}`

		authReq := httptest.NewRequest(http.MethodPost,
			"/api/v1/authenticateUser",
			strings.NewReader(authJson))
		authRes := httptest.NewRecorder()
		handler.ServeHTTP(authRes, authReq)

		tu.AssertEqual(http.StatusUnauthorized, authRes.Code)
	})
}

func TestAuthenticateUserMissingRequiredFields(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		missingUsernameAndEmail := `{
			"displayName": "estecat",
			"password": "password"
		}`

		missingPassword := `{
			"username": "esteban",
			"email": "estecat42069@yahoo.com",
			"displayName": "estecat",
		}`

		missingUsernameAndEmailRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/authenticateUser",
			strings.NewReader(missingUsernameAndEmail))
		missingUsernameAndEmailResponse := httptest.NewRecorder()

		missingPasswordRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/authenticateUser",
			strings.NewReader(missingPassword))
		missingPasswordResponse := httptest.NewRecorder()

		handler.ServeHTTP(missingUsernameAndEmailResponse, missingUsernameAndEmailRequest)
		handler.ServeHTTP(missingPasswordResponse, missingPasswordRequest)

		tu.AssertEqual(http.StatusBadRequest, missingUsernameAndEmailResponse.Code)
		tu.AssertEqual(http.StatusBadRequest, missingPasswordResponse.Code)
	})
}

func TestAuthenticateUserMalformedJSON(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	UserAPI := UserAPI{}

	malformedJSON := "alkj}"
	req := httptest.NewRequest(http.MethodPost, "/api/v1/authenticateUser", strings.NewReader(malformedJSON))
	res := httptest.NewRecorder()
	UserAPI.Authenticate(res, req)

	tu.AssertEqual(http.StatusBadRequest, res.Code)
}

func TestAuthenticateUserWrongMethod(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	UserAPI := UserAPI{}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/authenticateUser", nil)
	getRes := httptest.NewRecorder()
	putReq := httptest.NewRequest(http.MethodPut, "/api/v1/authenticateUser", nil)
	putRes := httptest.NewRecorder()
	patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/authenticateUser", nil)
	patchRes := httptest.NewRecorder()
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/authenticateUser", nil)
	deleteRes := httptest.NewRecorder()
	headReq := httptest.NewRequest(http.MethodHead, "/api/v1/authenticateUser", nil)
	headRes := httptest.NewRecorder()
	optionReq := httptest.NewRequest(http.MethodOptions, "/api/v1/authenticateUser", nil)
	optionRes := httptest.NewRecorder()
	traceReq := httptest.NewRequest(http.MethodTrace, "/api/v1/authenticateUser", nil)
	traceRes := httptest.NewRecorder()
	connectReq := httptest.NewRequest(http.MethodConnect, "/api/v1/authenticateUser", nil)
	connectRes := httptest.NewRecorder()

	UserAPI.Authenticate(getRes, getReq)
	UserAPI.Authenticate(putRes, putReq)
	UserAPI.Authenticate(patchRes, patchReq)
	UserAPI.Authenticate(deleteRes, deleteReq)
	UserAPI.Authenticate(headRes, headReq)
	UserAPI.Authenticate(optionRes, optionReq)
	UserAPI.Authenticate(traceRes, traceReq)
	UserAPI.Authenticate(connectRes, connectReq)

	tu.AssertEqual(http.StatusMethodNotAllowed, getRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, putRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, patchRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, deleteRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, headRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, optionRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, traceRes.Code)
	tu.AssertEqual(http.StatusMethodNotAllowed, connectRes.Code)
}
