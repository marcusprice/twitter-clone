package testhelpers

import (
	"database/sql"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
)

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
