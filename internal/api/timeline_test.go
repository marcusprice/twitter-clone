package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testhelpers"
	"github.com/marcusprice/twitter-clone/internal/testutil"
	"github.com/marcusprice/twitter-clone/internal/util"
)

func TestTimelineGet(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)
		user1 := controller.NewUserController(db)
		user1.ByID(1)
		user1.Login()
		token, _ := GenerateJWT(user1.ID())
		user2 := controller.NewUserController(db)
		user2.ByID(2)
		user3 := controller.NewUserController(db)
		user3.ByID(3)
		user1.Follow(user2.Username)

		limit := 10
		offset := 0
		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/timeline?limit=%d&offset=%d", limit, offset),
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)
		res := httptest.NewRecorder()
		user2Posts := testhelpers.QueryUserPosts(user2.ID(), db)

		handler.ServeHTTP(res, req)

		var payload TimelinePayload
		json.NewDecoder(res.Body).Decode(&payload)
		tu.AssertEqual(http.StatusOK, res.Code)
		tu.AssertEqual(10, len(payload.Posts))
		tu.AssertEqual(len(user2Posts)-10, payload.PostsRemaining)
		tu.AssertEqual(user2Posts[0].ID, payload.Posts[0].ID)
		tu.AssertEqual(user2Posts[0].Content, payload.Posts[0].Content)
		tu.AssertEqual(user2Posts[0].Image, payload.Posts[0].Image)
		tu.AssertEqual(user2Posts[0].LikeCount, payload.Posts[0].LikeCount)
		tu.AssertEqual(user2Posts[0].BookmarkCount, payload.Posts[0].BookmarkCount)
		tu.AssertEqual(user2Posts[0].RetweetCount, payload.Posts[0].RetweetCount)
		tu.AssertEqual(user2Posts[0].Author.Username, payload.Posts[0].Author.Username)
		tu.AssertEqual(user2Posts[0].Author.DisplayName, payload.Posts[0].Author.DisplayName)
		tu.AssertEqual(user2Posts[0].Author.Avatar, payload.Posts[0].Author.Avatar)
		tu.AssertEqual(util.ParseTime(user2Posts[0].CreatedAt), payload.Posts[0].CreatedAt)
		tu.AssertEqual(util.ParseTime(user2Posts[0].UpdatedAt), payload.Posts[0].UpdatedAt)
		tu.AssertEqual(user2Posts[0].Impressions+1, payload.Posts[0].Impressions)
		tu.AssertFalse(payload.Posts[0].IsRetweet)
		tu.AssertTrue(payload.HasMore)

		tu.AssertEqual(user2Posts[9].ID, payload.Posts[9].ID)
		tu.AssertEqual(user2Posts[9].Content, payload.Posts[9].Content)
		tu.AssertEqual(user2Posts[9].Image, payload.Posts[9].Image)
		tu.AssertEqual(user2Posts[9].LikeCount, payload.Posts[9].LikeCount)
		tu.AssertEqual(user2Posts[9].BookmarkCount, payload.Posts[9].BookmarkCount)
		tu.AssertEqual(user2Posts[9].RetweetCount, payload.Posts[9].RetweetCount)
		tu.AssertEqual(user2Posts[9].Author.Username, payload.Posts[9].Author.Username)
		tu.AssertEqual(user2Posts[9].Author.DisplayName, payload.Posts[9].Author.DisplayName)
		tu.AssertEqual(user2Posts[9].Author.Avatar, payload.Posts[9].Author.Avatar)
		tu.AssertEqual(util.ParseTime(user2Posts[9].CreatedAt), payload.Posts[9].CreatedAt)
		tu.AssertEqual(util.ParseTime(user2Posts[9].UpdatedAt), payload.Posts[9].UpdatedAt)
		tu.AssertEqual(user2Posts[9].Impressions+1, payload.Posts[9].Impressions)
		tu.AssertFalse(payload.Posts[9].IsRetweet)

		limit = 10
		offset = 10
		req = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/timeline?limit=%d&offset=%d", limit, offset),
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)
		res = httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		json.NewDecoder(res.Body).Decode(&payload)
		tu.AssertFalse(payload.HasMore)
		tu.AssertEqual(user2Posts[10].ID, payload.Posts[0].ID)
		tu.AssertEqual(
			user2Posts[len(user2Posts)-1].ID,
			payload.Posts[len(payload.Posts)-1].ID,
		)

		postInput := dtypes.PostInput{
			UserID:  user3.ID(),
			Content: "Strawberry fields forever",
			Image:   "strawberries.jpeg",
		}

		// user 2 retweets new post
		retweetedPostID := testhelpers.CreatePost(postInput, db)
		testhelpers.CreateRetweet(retweetedPostID, user2.ID(), db)

		limit = 10
		offset = 0
		req = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/timeline?limit=%d&offset=%d", limit, offset),
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)
		res = httptest.NewRecorder()

		handler.ServeHTTP(res, req)
		json.NewDecoder(res.Body).Decode(&payload)
		tu.AssertEqual(retweetedPostID, payload.Posts[0].ID)
		tu.AssertTrue(payload.Posts[0].IsRetweet)
	})
}

func TestTimelineGetBadRequest(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)
		user1 := controller.NewUserController(db)
		user1.ByID(1)
		token, _ := GenerateJWT(user1.ID())
		user1.Login()

		limitStr := "ljkahkljhas"
		offset := 0
		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/timeline?limit=%s&offset=%d", limitStr, offset),
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		tu.AssertEqual(http.StatusBadRequest, res.Code)

		limit := 200000000000
		offset = 0
		req = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/timeline?limit=%d&offset=%d", limit, offset),
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		body := res.Body.String()

		tu.AssertEqual(http.StatusBadRequest, res.Code)
		tu.AssertEqual(
			fmt.Sprintf("Too small of a limit, max limit: %d\n", MIN_LIMIT),
			body,
		)

		limit = -420
		offset = 0
		req = httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/v1/timeline?limit=%d&offset=%d", limit, offset),
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)
		res = httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		body = res.Body.String()

		tu.AssertEqual(http.StatusBadRequest, res.Code)
		tu.AssertEqual(
			fmt.Sprintf("Too large of a limit, max limit: %d\n", MAX_LIMIT),
			body,
		)
	})
}

func TestTimelineGetUnauthorized(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		noAuthHeaderReq := httptest.NewRequest(http.MethodGet, "/api/v1/timeline", nil)
		noAuthHeaderRes := httptest.NewRecorder()
		handler.ServeHTTP(noAuthHeaderRes, noAuthHeaderReq)

		headerNoTokenReq := httptest.NewRequest(http.MethodGet, "/api/v1/timeline", nil)
		headerNoTokenReq.Header.Set("Authorization", "Bearer ")
		headerNoTokenRes := httptest.NewRecorder()
		handler.ServeHTTP(headerNoTokenRes, headerNoTokenReq)

		headerWrongKeywordReq := httptest.NewRequest(http.MethodGet, "/api/v1/timeline", nil)
		headerWrongKeywordReq.Header.Set("Authorization", "Esteban ")
		headerWrongKeywordRes := httptest.NewRecorder()
		handler.ServeHTTP(headerWrongKeywordRes, headerWrongKeywordReq)

		badToken := generateBadToken()
		badTokenReq := httptest.NewRequest(http.MethodGet, "/api/v1/timeline", nil)
		badTokenReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", badToken))
		badTokenRes := httptest.NewRecorder()
		handler.ServeHTTP(badTokenRes, badTokenReq)

		tu.AssertEqual(http.StatusUnauthorized, noAuthHeaderRes.Code)
		tu.AssertEqual(http.StatusUnauthorized, headerNoTokenRes.Code)
		tu.AssertEqual(http.StatusUnauthorized, headerWrongKeywordRes.Code)
		tu.AssertEqual(http.StatusUnauthorized, badTokenRes.Code)
	})
}

func TestTimelineGetWrongMethod(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		handler := RegisterHandlers(db)

		postReq := httptest.NewRequest(http.MethodPost, "/api/v1/timeline", nil)
		postRes := httptest.NewRecorder()
		putReq := httptest.NewRequest(http.MethodPut, "/api/v1/timeline", nil)
		putRes := httptest.NewRecorder()
		patchReq := httptest.NewRequest(http.MethodPatch, "/api/v1/timeline", nil)
		patchRes := httptest.NewRecorder()
		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/timeline", nil)
		deleteRes := httptest.NewRecorder()
		headReq := httptest.NewRequest(http.MethodHead, "/api/v1/timeline", nil)
		headRes := httptest.NewRecorder()
		optionReq := httptest.NewRequest(http.MethodOptions, "/api/v1/timeline", nil)
		optionRes := httptest.NewRecorder()
		traceReq := httptest.NewRequest(http.MethodTrace, "/api/v1/timeline", nil)
		traceRes := httptest.NewRecorder()
		connectReq := httptest.NewRequest(http.MethodConnect, "/api/v1/timeline", nil)
		connectRes := httptest.NewRecorder()

		handler.ServeHTTP(postRes, postReq)
		handler.ServeHTTP(putRes, putReq)
		handler.ServeHTTP(patchRes, patchReq)
		handler.ServeHTTP(deleteRes, deleteReq)
		handler.ServeHTTP(headRes, headReq)
		handler.ServeHTTP(optionRes, optionReq)
		handler.ServeHTTP(traceRes, traceReq)
		handler.ServeHTTP(connectRes, connectReq)

		tu.AssertEqual(http.StatusMethodNotAllowed, postRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, putRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, patchRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, deleteRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, headRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, optionRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, traceRes.Code)
		tu.AssertEqual(http.StatusMethodNotAllowed, connectRes.Code)
	})
}
