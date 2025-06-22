package controller

import (
	"database/sql"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/constants"
	"github.com/marcusprice/twitter-clone/internal/testhelpers"
	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestTimelineSetID(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	timeline := &Timeline{}
	timeline2 := timeline.SetID(42069)
	tu.AssertEqual(timeline.userID, timeline2.userID)
	tu.AssertEqual(timeline.userID, timeline2.userID)
}

func TestTimelineGetPosts(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		timeline := NewTimelineController(db)
		user1 := NewUserController(db)
		user2 := NewUserController(db)
		user1.ByID(1)
		user2.ByID(2)
		user1.Follow(user2.Username)

		posts, postsRemaining, err := timeline.GetPosts(10, 0)
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(0, len(posts))
		tu.AssertEqual(-1, postsRemaining)

		timeline.userID = 42069
		posts, postsRemaining, err = timeline.GetPosts(10, 0)
		tu.AssertErrorNil(err)
		tu.AssertEqual(0, len(posts))
		tu.AssertEqual(0, postsRemaining)

		user2Posts := testhelpers.QueryUserPosts(user2.ID(), db)
		timeline.SetID(user1.ID())
		posts, postsRemaining, err = timeline.GetPosts(10, 0)
		tu.AssertTrue(len(posts) <= 10)
		tu.AssertEqual(user2Posts[0].ID, posts[0].ID)
		tu.AssertEqual(user2Posts[0].Content, posts[0].Content)
		tu.AssertEqual(user2Posts[0].Image, posts[0].Image)
		tu.AssertEqual(user2Posts[0].Impressions+1, posts[0].Impressions)
		tu.AssertEqual(user2Posts[0].BookmarkCount, posts[0].BookmarkCount)
		tu.AssertEqual(user2Posts[0].RetweetCount, posts[0].RetweetCount)
		tu.AssertEqual(user2Posts[0].LikeCount, posts[0].LikeCount)
		tu.AssertEqual(user2Posts[0].Author.DisplayName, posts[0].Author.DisplayName)
		tu.AssertEqual(user2Posts[0].Author.Username, posts[0].Author.Username)
		tu.AssertEqual(user2Posts[0].Author.Avatar, posts[0].Author.Avatar)
		tu.AssertEqual(user2Posts[0].CreatedAt, posts[0].CreatedAt.Format(constants.TIME_LAYOUT))
		tu.AssertEqual(user2Posts[0].UpdatedAt, posts[0].UpdatedAt.Format(constants.TIME_LAYOUT))

		tu.AssertEqual(user2Posts[9].ID, posts[9].ID)
		tu.AssertEqual(user2Posts[9].Content, posts[9].Content)
		tu.AssertEqual(user2Posts[9].Image, posts[9].Image)
		tu.AssertEqual(user2Posts[9].Impressions+1, posts[9].Impressions)
		tu.AssertEqual(user2Posts[9].BookmarkCount, posts[9].BookmarkCount)
		tu.AssertEqual(user2Posts[9].RetweetCount, posts[9].RetweetCount)
		tu.AssertEqual(user2Posts[9].LikeCount, posts[9].LikeCount)
		tu.AssertEqual(user2Posts[9].Author.DisplayName, posts[9].Author.DisplayName)
		tu.AssertEqual(user2Posts[9].Author.Username, posts[9].Author.Username)
		tu.AssertEqual(user2Posts[9].Author.Avatar, posts[9].Author.Avatar)
		tu.AssertEqual(user2Posts[9].CreatedAt, posts[9].CreatedAt.Format(constants.TIME_LAYOUT))
		tu.AssertEqual(user2Posts[9].UpdatedAt, posts[9].UpdatedAt.Format(constants.TIME_LAYOUT))

		posts, postsRemaining, err = timeline.GetPosts(10, 10)
		tu.AssertTrue(len(posts) <= 10)
		tu.AssertEqual(user2Posts[10].ID, posts[0].ID)
		tu.AssertEqual(user2Posts[10].Content, posts[0].Content)
		tu.AssertEqual(user2Posts[10].Image, posts[0].Image)
		tu.AssertEqual(user2Posts[10].Impressions+1, posts[0].Impressions)
		tu.AssertEqual(user2Posts[10].BookmarkCount, posts[0].BookmarkCount)
		tu.AssertEqual(user2Posts[10].RetweetCount, posts[0].RetweetCount)
		tu.AssertEqual(user2Posts[10].LikeCount, posts[0].LikeCount)
		tu.AssertEqual(user2Posts[10].Author.DisplayName, posts[0].Author.DisplayName)
		tu.AssertEqual(user2Posts[10].Author.Username, posts[0].Author.Username)
		tu.AssertEqual(user2Posts[10].Author.Avatar, posts[0].Author.Avatar)
		tu.AssertEqual(user2Posts[10].CreatedAt, posts[0].CreatedAt.Format(constants.TIME_LAYOUT))
		tu.AssertEqual(user2Posts[10].UpdatedAt, posts[0].UpdatedAt.Format(constants.TIME_LAYOUT))

		tu.AssertEqual(user2Posts[19].ID, posts[9].ID)
		tu.AssertEqual(user2Posts[19].Content, posts[9].Content)
		tu.AssertEqual(user2Posts[19].Image, posts[9].Image)
		tu.AssertEqual(user2Posts[19].Impressions+1, posts[9].Impressions)
		tu.AssertEqual(user2Posts[19].BookmarkCount, posts[9].BookmarkCount)
		tu.AssertEqual(user2Posts[19].RetweetCount, posts[9].RetweetCount)
		tu.AssertEqual(user2Posts[19].LikeCount, posts[9].LikeCount)
		tu.AssertEqual(user2Posts[19].Author.DisplayName, posts[9].Author.DisplayName)
		tu.AssertEqual(user2Posts[19].Author.Username, posts[9].Author.Username)
		tu.AssertEqual(user2Posts[19].Author.Avatar, posts[9].Author.Avatar)
		tu.AssertEqual(user2Posts[19].CreatedAt, posts[9].CreatedAt.Format(constants.TIME_LAYOUT))
		tu.AssertEqual(user2Posts[19].UpdatedAt, posts[9].UpdatedAt.Format(constants.TIME_LAYOUT))
	})
}
