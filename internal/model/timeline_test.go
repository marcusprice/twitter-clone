package model

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestTimelineOffsetCount(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		PostModel := &PostModel{db}
		user1 := queryUser(1, db)
		insertUserFollow(user1.ID, 2, db)
		insertUserFollow(user1.ID, 3, db)
		insertUserFollow(user1.ID, 4, db)

		count, err := PostModel.TimelineRemainingPostsCount(user1.ID, 0, 0)
		tu.AssertErrorNotNil(err)
		tu.AssertEqual("Positive limit value required", err.Error())
		tu.AssertEqual(-1, count)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, -42069, 0)
		tu.AssertErrorNotNil(err)
		tu.AssertEqual("Positive limit value required", err.Error())
		tu.AssertEqual(-1, count)

		// verify num of posts in case test db seed data changes
		numOfPosts := getNumOfPosts(db, 2, 3, 4)
		tu.AssertEqual(56, numOfPosts)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, 10, 0)
		tu.AssertErrorNil(err)
		tu.AssertEqual(46, count)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, 10, 10)
		tu.AssertErrorNil(err)
		tu.AssertEqual(36, count)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, 10, 20)
		tu.AssertErrorNil(err)
		tu.AssertEqual(26, count)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, 10, 30)
		tu.AssertErrorNil(err)
		tu.AssertEqual(16, count)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, 10, 40)
		tu.AssertErrorNil(err)
		tu.AssertEqual(6, count)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, 10, 50)
		tu.AssertErrorNil(err)
		tu.AssertEqual(0, count)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, 15, 0)
		tu.AssertErrorNil(err)
		tu.AssertEqual(41, count)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, 15, 15)
		tu.AssertErrorNil(err)
		tu.AssertEqual(26, count)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, 15, 30)
		tu.AssertErrorNil(err)
		tu.AssertEqual(11, count)

		count, err = PostModel.TimelineRemainingPostsCount(user1.ID, 15, 45)
		tu.AssertErrorNil(err)
		tu.AssertEqual(0, count)
	})
}

func getNumOfPosts(db *sql.DB, userIDs ...int) int {
	query := `
		SELECT
			COUNT(*)
		FROM Post
		WHERE
	`

	for index, id := range userIDs {
		if index == 0 {
			query += fmt.Sprintf(" user_id = %d", id)
		}

		if len(userIDs) > 1 && index != 0 {
			query += fmt.Sprintf(" OR user_id = %d", id)
		}

		if index == len(userIDs)-1 {
			query += ";"
		}
	}

	var count int

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		panic(err)
	}

	return count
}
