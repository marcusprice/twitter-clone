package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testhelpers"
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
			http.MethodPost, "/api/v1/user/create",
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
			http.MethodPost, "/api/v1/user/create",
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
			http.MethodPost, "/api/v1/user/create",
			strings.NewReader(duplicateUserJson))
		duplicateUserResponse := httptest.NewRecorder()

		duplicateEmailRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/user/create",
			strings.NewReader(duplicateUserJsonSameEmail))
		duplicateEmailResponse := httptest.NewRecorder()

		duplicateUsernameRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/user/create",
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
			http.MethodPost, "/api/v1/user/create",
			strings.NewReader(missingUsername))
		missingUsernameResponse := httptest.NewRecorder()

		missingEmailRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/user/create",
			strings.NewReader(missingEmail))
		missingEmailResponse := httptest.NewRecorder()

		missingDisplayNameRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/user/create",
			strings.NewReader(missingDisplayName))
		missingDisplayNameResponse := httptest.NewRecorder()

		missingPasswordRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/user/create",
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
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		malformedJSON := "alkj}"
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/create", strings.NewReader(malformedJSON))
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)

		tu.AssertEqual(http.StatusBadRequest, res.Code)
	})
}

func TestCreateUserWrongMethod(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/user/create", nil)
		getRes := httptest.NewRecorder()
		putReq := httptest.NewRequest(http.MethodPut, "/api/v1/user/create", nil)
		putRes := httptest.NewRecorder()
		patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/user/create", nil)
		patchRes := httptest.NewRecorder()
		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/user/create", nil)
		deleteRes := httptest.NewRecorder()
		headReq := httptest.NewRequest(http.MethodHead, "/api/v1/user/create", nil)
		headRes := httptest.NewRecorder()
		optionReq := httptest.NewRequest(http.MethodOptions, "/api/v1/user/create", nil)
		optionRes := httptest.NewRecorder()
		traceReq := httptest.NewRequest(http.MethodTrace, "/api/v1/user/create", nil)
		traceRes := httptest.NewRecorder()
		connectReq := httptest.NewRequest(http.MethodConnect, "/api/v1/user/create", nil)
		connectRes := httptest.NewRecorder()

		handler.ServeHTTP(getRes, getReq)
		handler.ServeHTTP(putRes, putReq)
		handler.ServeHTTP(patchRes, patchReq)
		handler.ServeHTTP(deleteRes, deleteReq)
		handler.ServeHTTP(headRes, headReq)
		handler.ServeHTTP(optionRes, optionReq)
		handler.ServeHTTP(traceRes, traceReq)
		handler.ServeHTTP(connectRes, connectReq)

		tu.AssertEqual(http.StatusMethodNotAllowed, getRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, putRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, patchRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, deleteRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, headRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, optionRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, traceRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, connectRes.Code)
	})
}

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
			"/api/v1/user/authenticate",
			strings.NewReader(authJson))
		authRes := httptest.NewRecorder()
		authWithUsernameReq := httptest.NewRequest(http.MethodPost,
			"/api/v1/user/authenticate",
			strings.NewReader(authUsernameOnlyJson))
		authWithUsernameRes := httptest.NewRecorder()
		authWithEmailReq := httptest.NewRequest(http.MethodPost,
			"/api/v1/user/authenticate",
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
		tu.AssertEqual("esteban", userPayload.Username)
		tu.AssertEqual("estecat42069@yahoo.com", userPayload.Email)
		tu.AssertEqual("yodel", userPayload.DisplayName)

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
			"/api/v1/user/authenticate",
			strings.NewReader(authJson))
		authRes := httptest.NewRecorder()
		handler.ServeHTTP(authRes, authReq)

		tu.AssertEqual(http.StatusUnauthorized, authRes.Code)
	})
}

func TestAuthenticateUserWrongUsernmae(t *testing.T) {
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
			"username": "esteba",
			"password": "password"
		}`

		authReq := httptest.NewRequest(http.MethodPost,
			"/api/v1/user/authenticate",
			strings.NewReader(authJson))
		authRes := httptest.NewRecorder()
		handler.ServeHTTP(authRes, authReq)

		tu.AssertEqual(http.StatusUnauthorized, authRes.Code)
	})
}

func TestAuthenticateUserWrongEmail(t *testing.T) {
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
			"email": "whispers_from_wallface@freakseasy.com",
			"password": "password"
		}`

		authReq := httptest.NewRequest(http.MethodPost,
			"/api/v1/user/authenticate",
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
			http.MethodPost, "/api/v1/user/authenticate",
			strings.NewReader(missingUsernameAndEmail))
		missingUsernameAndEmailResponse := httptest.NewRecorder()

		missingPasswordRequest := httptest.NewRequest(
			http.MethodPost, "/api/v1/user/authenticate",
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
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/authenticate", strings.NewReader(malformedJSON))
	res := httptest.NewRecorder()
	UserAPI.Authenticate(res, req)

	tu.AssertEqual(http.StatusBadRequest, res.Code)
}

func TestAuthenticateUserWrongMethod(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/user/authenticate", nil)
		getRes := httptest.NewRecorder()
		putReq := httptest.NewRequest(http.MethodPut, "/api/v1/user/authenticate", nil)
		putRes := httptest.NewRecorder()
		patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/user/authenticate", nil)
		patchRes := httptest.NewRecorder()
		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/user/authenticate", nil)
		deleteRes := httptest.NewRecorder()
		headReq := httptest.NewRequest(http.MethodHead, "/api/v1/user/authenticate", nil)
		headRes := httptest.NewRecorder()
		optionReq := httptest.NewRequest(http.MethodOptions, "/api/v1/user/authenticate", nil)
		optionRes := httptest.NewRecorder()
		traceReq := httptest.NewRequest(http.MethodTrace, "/api/v1/user/authenticate", nil)
		traceRes := httptest.NewRecorder()
		connectReq := httptest.NewRequest(http.MethodConnect, "/api/v1/user/authenticate", nil)
		connectRes := httptest.NewRecorder()

		handler.ServeHTTP(getRes, getReq)
		handler.ServeHTTP(putRes, putReq)
		handler.ServeHTTP(patchRes, patchReq)
		handler.ServeHTTP(deleteRes, deleteReq)
		handler.ServeHTTP(headRes, headReq)
		handler.ServeHTTP(optionRes, optionReq)
		handler.ServeHTTP(traceRes, traceReq)
		handler.ServeHTTP(connectRes, connectReq)

		tu.AssertEqual(http.StatusMethodNotAllowed, getRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, putRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, patchRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, deleteRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, headRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, optionRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, traceRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, connectRes.Code)
	})
}

func TestFollowUser(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)
		user1 := controller.NewUserController(db)
		user2 := controller.NewUserController(db)
		user3 := controller.NewUserController(db)
		user1.ByID(1)
		user2.ByID(2)
		user3.ByID(3)
		user1.Login()
		user3.Login()
		user1Token, _ := GenerateJWT(user1.ID())
		user3Token, _ := GenerateJWT(user3.ID())

		req := httptest.NewRequest(
			http.MethodPut, fmt.Sprintf("/api/v1/user/%s/follow", user2.Username), nil)
		req.Header.Set("Authorization", "Bearer "+user1Token)
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)

		req = httptest.NewRequest(
			http.MethodPut, fmt.Sprintf("/api/v1/user/%s/follow", user2.Username), nil)
		req.Header.Set("Authorization", "Bearer "+user3Token)
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		tu.AssertEqual(http.StatusNoContent, res.Code)

		userFollowers := testhelpers.QueryUserFollowers(user2.ID(), db)
		tu.AssertEqual(2, len(userFollowers))
		tu.AssertEqual(user1.ID(), userFollowers[0].ID)
		tu.AssertEqual(user3.ID(), userFollowers[1].ID)

		// unfollow
		req = httptest.NewRequest(
			http.MethodDelete, fmt.Sprintf("/api/v1/user/%s/follow", user2.Username), nil)
		req.Header.Set("Authorization", "Bearer "+user3Token)
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)

		userFollowers = testhelpers.QueryUserFollowers(user2.ID(), db)
		tu.AssertEqual(http.StatusNoContent, res.Code)
		tu.AssertEqual(1, len(userFollowers))
		tu.AssertEqual(user1.ID(), userFollowers[0].ID)

		// duplicate requests okay
		req = httptest.NewRequest(
			http.MethodDelete, fmt.Sprintf("/api/v1/user/%s/follow", user2.Username), nil)
		req.Header.Set("Authorization", "Bearer "+user3Token)
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)

		userFollowers = testhelpers.QueryUserFollowers(user2.ID(), db)
		tu.AssertEqual(http.StatusNoContent, res.Code)
		tu.AssertEqual(1, len(userFollowers))
		tu.AssertEqual(user1.ID(), userFollowers[0].ID)

		req = httptest.NewRequest(
			http.MethodPut, fmt.Sprintf("/api/v1/user/%s/follow", "made-up-user-name"), nil)
		req.Header.Set("Authorization", "Bearer "+user3Token)
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)

		tu.AssertEqual(http.StatusNotFound, res.Code)

		req = httptest.NewRequest(
			http.MethodDelete, fmt.Sprintf("/api/v1/user/%s/follow", "made-up-user-name"), nil)
		req.Header.Set("Authorization", "Bearer "+user3Token)
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)

		tu.AssertEqual(http.StatusNotFound, res.Code)
	})
}

func TestFollowUserWrongMethod(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/user/esteban/follow", nil)
		getRes := httptest.NewRecorder()
		postReq := httptest.NewRequest(http.MethodPost, "/api/v1/user/esteban/follow", nil)
		postRes := httptest.NewRecorder()
		patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/user/esteban/follow", nil)
		patchRes := httptest.NewRecorder()
		headReq := httptest.NewRequest(http.MethodHead, "/api/v1/user/esteban/follow", nil)
		headRes := httptest.NewRecorder()
		optionReq := httptest.NewRequest(http.MethodOptions, "/api/v1/user/esteban/follow", nil)
		optionRes := httptest.NewRecorder()
		traceReq := httptest.NewRequest(http.MethodTrace, "/api/v1/user/esteban/follow", nil)
		traceRes := httptest.NewRecorder()
		connectReq := httptest.NewRequest(http.MethodConnect, "/api/v1/user/esteban/follow", nil)
		connectRes := httptest.NewRecorder()

		handler.ServeHTTP(getRes, getReq)
		handler.ServeHTTP(postRes, postReq)
		handler.ServeHTTP(patchRes, patchReq)
		handler.ServeHTTP(headRes, headReq)
		handler.ServeHTTP(optionRes, optionReq)
		handler.ServeHTTP(traceRes, traceReq)
		handler.ServeHTTP(connectRes, connectReq)

		tu.AssertEqual(http.StatusMethodNotAllowed, getRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, postRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, patchRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, headRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, optionRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, traceRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, connectRes.Code)
	})
}
