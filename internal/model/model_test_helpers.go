package model

import "database/sql"

func QueryUserFollowTableCount(db *sql.DB) (count int) {
	query := `
		SELECT COUNT(*)
		FROM UserFollows;
	`

	db.QueryRow(query).Scan(&count)

	return count
}
