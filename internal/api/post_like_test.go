package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestPostLikeSimple(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)
		user := createTestUser(db)
		user.Login()
		token, _ := GenerateJWT(user.ID())
		post := createTestPost(user.ID(), db)

		req := httptest.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/v1/post/%d/like", post.ID),
			nil,
		)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		post.Sync()
		tu.AssertEqual(http.StatusNoContent, res.Code)
		tu.AssertEqual(1, post.LikeCount)

		req = httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/v1/post/%d/like", post.ID),
			nil,
		)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		post.Sync()
		tu.AssertEqual(http.StatusNoContent, res.Code)
		tu.AssertEqual(0, post.LikeCount)
	})
}

func TestPostLikeComprehensive(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		endpoint := "/api/v1/post/%d/like"
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)
		user1 := loadUserControllerByID(db, 1)
		user2 := loadUserControllerByID(db, 2)
		user3 := loadUserControllerByID(db, 3)
		user1Token := loginAndToken(user1)
		user2Token := loginAndToken(user2)
		user3Token := loginAndToken(user3)
		post := loadPostControllerByID(db, 1)

		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, post.ID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user1Token))
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		post.Sync()
		tu.AssertEqual(http.StatusNoContent, res.Code)
		tu.AssertEqual(1, post.LikeCount)

		req = httptest.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, post.ID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user2Token))
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		post.Sync()
		tu.AssertEqual(http.StatusNoContent, res.Code)
		tu.AssertEqual(2, post.LikeCount)

		req = httptest.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, post.ID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user3Token))
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		post.Sync()
		tu.AssertEqual(http.StatusNoContent, res.Code)
		tu.AssertEqual(3, post.LikeCount)

		// if user double-likes a post, API sends the same StatusNoContent
		// response but LikeCount isn't effected
		req = httptest.NewRequest(http.MethodPut, fmt.Sprintf(endpoint, post.ID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user3Token))
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		post.Sync()
		tu.AssertEqual(http.StatusNoContent, res.Code)
		tu.AssertEqual(3, post.LikeCount)

		req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf(endpoint, post.ID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user2Token))
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		post.Sync()
		tu.AssertEqual(http.StatusNoContent, res.Code)
		tu.AssertEqual(2, post.LikeCount)

		// if user double-unlikes a post, API sends the same StatusNoContent
		// response but LikeCount isn't effected
		req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf(endpoint, post.ID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user2Token))
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		post.Sync()
		tu.AssertEqual(http.StatusNoContent, res.Code)
		tu.AssertEqual(2, post.LikeCount)
	})
}

func TestLikePostMissingPost(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)
		user := createTestUser(db)
		user.Login()
		token, _ := GenerateJWT(user.ID())

		req := httptest.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/v1/post/%d/like", 42069),
			nil,
		)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		tu.AssertEqual(http.StatusNotFound, res.Code)
	})
}

func TestCreatePostLikeWrongMethod(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/post/1/like", nil)
		getRes := httptest.NewRecorder()
		postReq := httptest.NewRequest(http.MethodPost, "/api/v1/post/1/like", nil)
		postRes := httptest.NewRecorder()
		patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/post/1/like", nil)
		patchRes := httptest.NewRecorder()
		headReq := httptest.NewRequest(http.MethodHead, "/api/v1/post/1/like", nil)
		headRes := httptest.NewRecorder()
		optionReq := httptest.NewRequest(http.MethodOptions, "/api/v1/post/1/like", nil)
		optionRes := httptest.NewRecorder()
		traceReq := httptest.NewRequest(http.MethodTrace, "/api/v1/post/1/like", nil)
		traceRes := httptest.NewRecorder()
		connectReq := httptest.NewRequest(http.MethodConnect, "/api/v1/post/1/like", nil)
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

func TestPostLikeUnauthorized(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		noAuthHeaderReq := httptest.NewRequest(http.MethodPut, "/api/v1/post/1/like", nil)
		noAuthHeaderRes := httptest.NewRecorder()
		handler.ServeHTTP(noAuthHeaderRes, noAuthHeaderReq)

		headerNoTokenReq := httptest.NewRequest(http.MethodDelete, "/api/v1/post/1/like", nil)
		headerNoTokenReq.Header.Set("Authorization", "Bearer ")
		headerNoTokenRes := httptest.NewRecorder()
		handler.ServeHTTP(headerNoTokenRes, headerNoTokenReq)

		headerWrongKeywordReq := httptest.NewRequest(http.MethodPut, "/api/v1/post/1/like", nil)
		headerWrongKeywordReq.Header.Set("Authorization", "Esteban ")
		headerWrongKeywordRes := httptest.NewRecorder()
		handler.ServeHTTP(headerWrongKeywordRes, headerWrongKeywordReq)

		badToken := generateBadToken()
		badTokenReq := httptest.NewRequest(http.MethodDelete, "/api/v1/post/1/like", nil)
		badTokenReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", badToken))
		badTokenRes := httptest.NewRecorder()
		handler.ServeHTTP(badTokenRes, badTokenReq)

		tu.AssertEqual(http.StatusUnauthorized, noAuthHeaderRes.Code)
		tu.AssertEqual(http.StatusUnauthorized, headerNoTokenRes.Code)
		tu.AssertEqual(http.StatusUnauthorized, headerWrongKeywordRes.Code)
		tu.AssertEqual(http.StatusUnauthorized, badTokenRes.Code)
	})
}

func createTestPost(userID int, db *sql.DB) *controller.Post {
	post := controller.NewPostController(db)
	postInput := dtypes.PostInput{
		UserID:  userID,
		Content: "Cats are cool",
		Image:   "smiley-cat.png",
	}

	post.New(postInput)

	return post
}

func loadUserControllerByID(db *sql.DB, userID int) *controller.User {
	if db == nil {
		panic("db conn cannot be nil")
	}

	user := controller.NewUserController(db)
	err := user.ByID(userID)
	if err != nil {
		log.Fatal("error loading user controller by ID:", err)
	}

	return user
}

func loadPostControllerByID(db *sql.DB, postID int) *controller.Post {
	if db == nil {
		panic("db conn cannot be nil")
	}

	post := controller.NewPostController(db)
	err := post.ByID(postID)
	if err != nil {
		log.Fatal("error loading user controller by ID:", err)
	}

	return post
}

func loginAndToken(user *controller.User) (token string) {
	user.Login()
	token, err := GenerateJWT(user.ID())
	if err != nil {
		log.Fatal(err)
	}

	return token
}
