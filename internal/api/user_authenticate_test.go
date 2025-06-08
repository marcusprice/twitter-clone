package api

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/marcusprice/twitter-clone/internal/testutil"
)

// func TestAuthenticateUser(t *testing.T) {
// 	testutil.WithTestDB(t, func(db *sql.DB) {

// 	})
// }

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
