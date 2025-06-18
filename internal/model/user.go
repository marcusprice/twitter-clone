package model

import (
	"database/sql"
	_ "embed"
	"errors"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/util"
)

type UserModel struct {
	db *sql.DB
}

//go:embed queries/create-user.sql
var createUserQuery string

func (um *UserModel) New(userInput dtypes.UserInput) (UserData, error) {
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
			return UserData{}, dbutils.WrapConstraintError(err)
		}

		return UserData{}, err
	}

	out := UserData{
		ID:          userID,
		Email:       userInput.Email,
		Username:    userInput.Username,
		FirstName:   userInput.FirstName,
		LastName:    userInput.LastName,
		DisplayName: userInput.DisplayName,
		LastLogin:   "", // last login null in the db
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

func (um *UserModel) GetByID(userID int) (UserData, error) {
	if userID == 0 {
		return UserData{}, errors.New("userID required")
	}

	query := selectUserBaseQuery + "WHERE id = $1;"
	row := um.db.QueryRow(query, userID)

	return parseUserQueryRow(row)
}

func (um *UserModel) GetByIdentifier(email, username string) (UserData, error) {
	if email == "" && username == "" {
		return UserData{}, MissingRequiredFilterData{}
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

func parseUserQueryRow(row *sql.Row) (UserData, error) {
	var id int
	var email string
	var userName string
	var password string
	var firstName string
	var lastName string
	var displayName string
	var lastLogin sql.NullString
	var isActive int
	var createdAt string
	var updatedAt string

	err := row.Scan(
		&id, &email, &userName, &password, &firstName, &lastName, &displayName,
		&lastLogin, &isActive, &createdAt, &updatedAt)

	if err != nil {
		return UserData{}, UserNotFoundError{}
	}

	lastLoginString := ""
	if lastLogin.Valid {
		lastLoginString = lastLogin.String
	}

	return UserData{
		id,
		email,
		userName,
		firstName,
		lastName,
		displayName,
		password,
		lastLoginString,
		isActive,
		createdAt,
		updatedAt,
	}, nil
}
