package model

import (
	"database/sql"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testutil"
	"github.com/marcusprice/twitter-clone/internal/util"
)

func TestPostCreate(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		postModel := NewPostModel(db)
		userInput := dtypes.UserInput{
			Email:       "estecat42069@yahoo.com",
			Username:    "estecat",
			FirstName:   "esteban",
			LastName:    "price",
			DisplayName: "Bubba",
			Password:    "password",
		}
		userID := insertUser(userInput, db)

		postInput := dtypes.PostInput{
			UserID:  userID,
			Content: "Cats are cool",
			Image:   "cool-cats.png",
		}

		beforeAction := time.Now().UTC().Add(-1 * time.Minute)
		postID, err := postModel.Create(postInput)
		afterAction := time.Now().UTC().Add(time.Minute)

		postData := queryPost(postID, db)
		createdAt := util.ParseTime(postData.CreatedAt)
		updatedAt := util.ParseTime(postData.UpdatedAt)

		tu.AssertEqual("Cats are cool", postData.Content)
		tu.AssertEqual("cool-cats.png", postData.Image)
		tu.AssertErrorNil(err)
		tu.AssertEqual(1, postID)
		tu.AssertEqual(postID, postData.ID)
		tu.AssertEqual(userID, postData.UserID)
		tu.AssertEqual(0, postData.LikeCount)
		tu.AssertEqual(0, postData.RetweetCount)
		tu.AssertEqual(0, postData.BookmarkCount)
		tu.AssertEqual(0, postData.Impressions)
		tu.AssertTrue(createdAt.After(beforeAction))
		tu.AssertTrue(createdAt.Before(afterAction))
		tu.AssertTrue(updatedAt.After(beforeAction))
		tu.AssertTrue(updatedAt.Before(afterAction))
	})
}

func queryPost(id int, db *sql.DB) PostData {
	query := `
		SELECT
			id,
			user_id,
			content,
			like_count,
			retweet_count,
			bookmark_count,
			impressions,
			image,
			created_at,
			updated_at
		FROM Post
		WHERE id = $1;
	`
	var postID int
	var userID int
	var content string
	var likeCount int
	var retweetCount int
	var bookmarkCount int
	var impressions int
	var image string
	var createdAt string
	var updatedAt string

	db.
		QueryRow(query, id).
		Scan(
			&postID, &userID, &content, &likeCount, &retweetCount,
			&bookmarkCount, &impressions, &image, &createdAt, &updatedAt)

	postData := PostData{
		ID:            postID,
		UserID:        userID,
		Content:       content,
		LikeCount:     likeCount,
		RetweetCount:  retweetCount,
		BookmarkCount: bookmarkCount,
		Impressions:   impressions,
		Image:         image,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return postData
}
