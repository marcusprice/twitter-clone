package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestPostLike(t *testing.T) {
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
