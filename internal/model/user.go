package model

import (
	"database/sql"
	_ "embed"
	"errors"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/logger"
	"github.com/marcusprice/twitter-clone/internal/permissions"
	"github.com/marcusprice/twitter-clone/internal/util"
)

type UserModel struct {
	db *sql.DB
}

//go:embed queries/create-user.sql
var createUserQuery string

func (um *UserModel) New(userInput dtypes.UserInput) (dtypes.UserData, error) {
	var userID int
	var lastLogin sql.NullString
	var createdAt string
	var updatedAt string

	err := um.db.QueryRow(
		createUserQuery,
		userInput.Email,
		userInput.Username,
		userInput.Password,
		userInput.FirstName,
		userInput.LastName,
		userInput.DisplayName,
	).Scan(&userID, &lastLogin, &createdAt, &updatedAt)

	if err != nil {
		if dbutils.ConstraintFailed(err) {
			return dtypes.UserData{}, dbutils.WrapConstraintError(err)
		}

		logger.LogError("error creating user: " + err.Error())
		return dtypes.UserData{}, err
	}

	out := dtypes.UserData{
		ID:          userID,
		Email:       userInput.Email,
		Username:    userInput.Username,
		FirstName:   userInput.FirstName,
		LastName:    userInput.LastName,
		DisplayName: userInput.DisplayName,
		LastLogin:   "", // last login null in the db
		Role:        permissions.USER_ROLE,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return out, nil
}

//go:embed queries/create-user-follows.sql
var createUserFollowsQuery string

func (um *UserModel) Follow(followerID, followeeID int) error {
	result, err := um.db.Exec(createUserFollowsQuery, followerID, followeeID)
	if err != nil {
		if dbutils.IsUniqueConstraintError(err) {
			// user already likes this post, likely a duplicate request
			return nil
		}

		if dbutils.ConstraintFailed(err) {
			return dbutils.WrapConstraintError(err)
		}

		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil

}

//go:embed queries/delete-user-follows.sql
var deleteUserFollowsQuery string

func (um *UserModel) UnFollow(followerID, followeeID int) error {
	result, err := um.db.Exec(deleteUserFollowsQuery, followerID, followeeID)

	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

//go:embed queries/select-user-base-query.sql
var selectUserBaseQuery string

func (um *UserModel) GetByID(userID int) (dtypes.UserData, error) {
	if userID == 0 {
		return dtypes.UserData{}, errors.New("userID required")
	}

	query := selectUserBaseQuery + "WHERE id = $1;"
	row := um.db.QueryRow(query, userID)

	return parseUserQueryRow(row)
}

//go:embed queries/select-user-bookmark-count.sql
var selectUserBookmarkCountQuery string

func (um *UserModel) GetBookmarkCount(userID int) (int, error) {
	var count int
	err := um.db.QueryRow(selectUserBookmarkCountQuery, userID).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

//go:embed queries/select-user-bookmarks.sql
var selectUserBookmarksQuery string

func (um *UserModel) GetBookmarks(userID, limit, offset int) ([]dtypes.BookmarkData, error) {
	result, err := um.db.Query(selectUserBookmarksQuery, userID, limit, offset)
	if err != nil {
		return []dtypes.BookmarkData{}, err
	}

	var bookmarks []dtypes.BookmarkData
	for result.Next() {
		var bookmark_created_at string
		var id int
		var content string
		var image string
		var like_count int
		var retweet_count int
		var bookmark_count int
		var impressions int
		var created_at string
		var updated_at string
		var author_user_name string
		var author_display_name string
		var author_avatar string
		var content_type string

		err := result.Scan(
			&bookmark_created_at, &id, &content, &image, &like_count,
			&retweet_count, &bookmark_count, &impressions, &created_at,
			&updated_at, &author_user_name, &author_display_name,
			&author_avatar, &content_type,
		)

		author := dtypes.Author{
			Username:    author_user_name,
			DisplayName: author_display_name,
			Avatar:      author_avatar,
		}

		bookmarkData := dtypes.BookmarkData{
			BookmarkCreatedAt: bookmark_created_at,
			ID:                id,
			Content:           content,
			Image:             image,
			LikeCount:         like_count,
			RetweetCount:      retweet_count,
			BookmarkCount:     bookmark_count,
			Impressions:       impressions,
			CreatedAt:         created_at,
			UpdatedAt:         updated_at,
			Author:            author,
			Type:              content_type,
		}
		bookmarks = append(bookmarks, bookmarkData)
		if err != nil {
			logger.LogError("UserModel.GetBookmarks() errors scanning row: " + err.Error())
			return []dtypes.BookmarkData{}, err
		}
	}

	return bookmarks, nil
}

func (um *UserModel) GetByIdentifier(email, username string) (dtypes.UserData, error) {
	if email == "" && username == "" {
		return dtypes.UserData{}, MissingRequiredFilterData{}
	}

	filterValue := ""
	query := selectUserBaseQuery
	if email != "" {
		query += "WHERE email = $1;"
		filterValue = email
	} else {
		query += "WHERE user_name = $1;"
		filterValue = username
	}

	row := um.db.QueryRow(query, filterValue)
	return parseUserQueryRow(row)
}

//go:embed queries/user-login.sql
var userLoginQuery string

func (um *UserModel) Login(userID int) (lastLoginTime string, isActive int, err error) {
	if userID == 0 {
		err := errors.New("Missing user ID")
		if util.InDevContext() {
			panic(err)
		} else {
			return "", 0, err
		}
	}

	err = um.db.QueryRow(userLoginQuery, userID).Scan(&lastLoginTime, &isActive)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return "", 0, UserNotFoundError{}
		} else {
			return "", 0, err
		}
	}

	return lastLoginTime, isActive, err
}

//go:embed queries/check-unique-user.sql
var checkUniqueUserQuery string

func (um *UserModel) UsernameOrEmailExists(email, username string) (bool, error) {
	var count int
	err := um.db.QueryRow(checkUniqueUserQuery, email, username).Scan(&count)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			return false, err
		}
	}

	return count > 0, nil
}

//go:embed queries/check-user-exists.sql
var checkUserExistsQuery string

func (um *UserModel) Exists(userID int) (bool, error) {
	var count int
	err := um.db.QueryRow(checkUserExistsQuery, userID).Scan(&count)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			return false, err
		}
	}

	return count > 0, nil
}

func NewUserModel(dbConn *sql.DB) *UserModel {
	if dbConn == nil {
		panic("db conn cannot be nil")
	}

	return &UserModel{db: dbConn}
}

func parseUserQueryRow(row *sql.Row) (dtypes.UserData, error) {
	var id int
	var email string
	var userName string
	var password string
	var firstName string
	var lastName string
	var displayName string
	var avatar string
	var lastLogin sql.NullString
	var isActive int
	var role int
	var createdAt string
	var updatedAt string

	err := row.Scan(
		&id, &email, &userName, &password, &firstName, &lastName, &displayName,
		&avatar, &lastLogin, &isActive, &role, &createdAt, &updatedAt)

	if err != nil {
		return dtypes.UserData{}, UserNotFoundError{}
	}

	lastLoginString := ""
	if lastLogin.Valid {
		lastLoginString = lastLogin.String
	}

	return dtypes.UserData{
		ID:          id,
		Email:       email,
		Username:    userName,
		FirstName:   firstName,
		LastName:    lastName,
		DisplayName: displayName,
		Avatar:      avatar,
		Password:    password,
		LastLogin:   lastLoginString,
		IsActive:    isActive,
		Role:        permissions.Role(role),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
