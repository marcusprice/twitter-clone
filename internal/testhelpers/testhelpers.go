package testhelpers

import (
	"database/sql"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
)

func QueryUser(userID int, db *sql.DB) dtypes.UserData {
	type UserData struct {
		ID          int
		Email       string
		Username    string
		FirstName   string
		LastName    string
		DisplayName string
		Password    string
		LastLogin   string
		IsActive    int
		CreatedAt   string
		UpdatedAt   string
	}

	query := `
		SELECT
			id,
			email,
			user_name,
			first_name,
			last_name,
			display_name,
			password,
			last_login,
			is_active,
			created_at,
			updated_at
		FROM User
		WHERE id = $1;
	`

	var id int
	var email string
	var user_name string
	var first_name string
	var last_name string
	var display_name string
	var password string
	var last_login_ns sql.NullString
	var is_active int
	var created_at string
	var updated_at string

	err := db.
		QueryRow(query, userID).
		Scan(
			&id, &email, &user_name, &first_name, &last_name, &display_name,
			&password, &last_login_ns, &is_active, &created_at, &updated_at)

	if err != nil {
		panic(err)
	}

	userData := dtypes.UserData{
		ID:          id,
		Email:       email,
		Username:    user_name,
		FirstName:   first_name,
		LastName:    last_name,
		DisplayName: display_name,
		Password:    password,
		LastLogin:   last_login_ns.String,
		IsActive:    is_active,
		CreatedAt:   created_at,
		UpdatedAt:   updated_at,
	}

	return userData
}

func QueryUserFollowTableCount(db *sql.DB) (count int) {
	query := `
		SELECT COUNT(*)
		FROM UserFollows;
	`

	db.QueryRow(query).Scan(&count)

	return count
}

func QueryUserFollowers(followeeID int, db *sql.DB) []dtypes.UserData {
	query := `
		SELECT 
			User.id,
			User.email,
			User.user_name,
			User.first_name,
			User.last_name,
			User.display_name,
			User.password,
			User.last_login,
			User.is_active,
			User.created_at,
			User.updated_at
		FROM 
			User
			INNER JOIN UserFollows ON UserFollows.follower_id = User.id
		WHERE
			UserFollows.followee_id = $1;
	`

	rows, err := db.Query(query, followeeID)
	if err != nil {
		panic("DB QUERY FAILED:" + err.Error())
	}
	defer rows.Close()

	var userFollowersData []dtypes.UserData
	for rows.Next() {
		var id int
		var email string
		var user_name string
		var first_name string
		var last_name string
		var display_name string
		var password string
		var last_login string
		var is_active int
		var created_at string
		var updated_at string
		err := rows.Scan(
			&id, &email, &user_name, &first_name, &last_name, &display_name,
			&password, &last_login, &is_active, &created_at, &updated_at)

		if err != nil {
			panic("DB SCAN FAILED: " + err.Error())
		}

		userData := dtypes.UserData{
			ID:          id,
			Email:       email,
			Username:    user_name,
			FirstName:   first_name,
			LastName:    last_name,
			DisplayName: display_name,
			Password:    password,
			LastLogin:   last_login,
			IsActive:    is_active,
			CreatedAt:   created_at,
			UpdatedAt:   updated_at,
		}

		userFollowersData = append(userFollowersData, userData)
	}

	return userFollowersData
}

func QueryUserPosts(userID int, db *sql.DB) []dtypes.PostData {
	query := `
		SELECT
			Post.id,
			Post.user_id,
			Post.content,
			Post.like_count,
			Post.retweet_count,
			Post.bookmark_count,
			Post.impressions,
			Post.image,
			Post.created_at,
			Post.updated_at,
			User.user_name,
			User.display_name,
			User.avatar
		FROM 
			Post
			INNER JOIN User ON User.id = Post.user_id
		WHERE Post.user_id = $1
		ORDER BY
			Post.created_at DESC;
	`

	result, err := db.Query(query, userID)
	if err != nil {
		panic(err)
	}

	var postRows []dtypes.PostData
	for result.Next() {
		var id int
		var user_id int
		var content string
		var like_count int
		var retweet_count int
		var bookmark_count int
		var impressions int
		var image string
		var created_at string
		var updated_at string
		var author_user_name string
		var author_display_name string
		var author_avatar string

		err := result.Scan(
			&id, &user_id, &content, &like_count, &retweet_count,
			&bookmark_count, &impressions, &image, &created_at, &updated_at,
			&author_user_name, &author_display_name, &author_avatar)
		if err != nil {
			panic(err)
		}

		postAuthor := dtypes.Author{
			Username:    author_user_name,
			DisplayName: author_display_name,
			Avatar:      author_avatar,
		}

		postRow := dtypes.PostData{
			Author:        postAuthor,
			ID:            id,
			UserID:        user_id,
			Content:       content,
			LikeCount:     like_count,
			RetweetCount:  retweet_count,
			BookmarkCount: bookmark_count,
			Impressions:   impressions,
			Image:         image,
			CreatedAt:     created_at,
			UpdatedAt:     updated_at,
		}

		postRows = append(postRows, postRow)
	}

	return postRows
}

func QueryComment(commentID int, db *sql.DB) dtypes.CommentData {
	var id int
	var post_id int
	var user_id int
	var depth int
	var parent_comment_id sql.NullInt64
	var content string
	var image string
	var like_count int
	var retweet_count int
	var bookmark_count int
	var impressions int
	var created_at string
	var updated_at string
	var author_username string
	var author_display_name string
	var author_avatar string

	query := `
		SELECT
			Comment.id,
			Comment.post_id,
			Comment.user_id,
			Comment.depth,
			Comment.parent_comment_id,
			Comment.content,
			Comment.image,
			Comment.like_count,
			Comment.retweet_count,
			Comment.bookmark_count,
			Comment.impressions,
			Comment.created_at,
			Comment.updated_at,
			User.user_name,
			User.display_name,
			User.avatar
		FROM
			Comment
			INNER JOIN User ON User.id = Comment.user_id
		WHERE
			Comment.id = $1;
	`

	err := db.
		QueryRow(query, commentID).
		Scan(
			&id, &post_id, &user_id, &depth, &parent_comment_id, &content,
			&image, &like_count, &retweet_count, &bookmark_count, &impressions,
			&created_at, &updated_at, &author_username, &author_display_name,
			&author_avatar)

	if err != nil {
		panic(err)
	}

	author := dtypes.Author{
		Username:    author_username,
		DisplayName: author_display_name,
		Avatar:      author_avatar,
	}

	commentData := dtypes.CommentData{
		ID:              id,
		PostID:          post_id,
		UserID:          user_id,
		Depth:           depth,
		ParentCommentID: int(parent_comment_id.Int64),
		Content:         content,
		Image:           image,
		LikeCount:       like_count,
		RetweetCount:    retweet_count,
		BookmarkCount:   bookmark_count,
		Impressions:     impressions,
		CreatedAt:       created_at,
		UpdatedAt:       updated_at,
		Author:          author,
	}

	return commentData
}

func CreateUserFollows(followerID, followeeID int, db *sql.DB) (rowID int) {
	query := `
		INSERT INTO UserFollows (follower_id, followee_id)
		VALUES ($1, $2)
		RETURNING id;
	`

	err := db.QueryRow(query, followerID, followeeID).Scan(&rowID)
	if err != nil {
		panic(err)
	}

	return rowID
}

func CreatePost(postInput dtypes.PostInput, db *sql.DB) (rowID int) {
	query := `
		INSERT INTO Post (content, image, user_id) VALUES ($1, $2, $3)
		RETURNING id;
	`

	err := db.
		QueryRow(query, postInput.Content, postInput.Image, postInput.UserID).
		Scan(&rowID)

	if err != nil {
		panic(err)
	}

	return rowID
}

func CreateRetweet(postID, userID int, db *sql.DB) (rowID int) {
	query := `
		INSERT INTO PostRetweet (post_id, user_id)
		VALUES ($1, $2)
		RETURNING id;
	`

	err := db.QueryRow(query, postID, userID).Scan(&rowID)
	if err != nil {
		panic(err)
	}

	return rowID
}

func CreateComment(commentInput dtypes.CommentInput, db *sql.DB) (rowID int) {
	query := `
		INSERT INTO Comment
			(post_id, user_id, content, image, depth)
		VALUES
			($1, $2, $3, $4, 0)
		RETURNING
			id;
	`

	err := db.
		QueryRow(
			query, commentInput.PostID, commentInput.UserID,
			commentInput.Content, commentInput.Image).
		Scan(&rowID)

	if err != nil {
		panic(err)
	}

	return rowID
}
