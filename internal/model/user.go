package model

import (
	"database/sql"
	_ "embed"
	"errors"

	"github.com/marcusprice/twitter-clone/internal/db"
	"github.com/marcusprice/twitter-clone/internal/util"
)

type UserModel struct {
	db *sql.DB
}

//go:embed queries/create-user.sql
var createUserQuery string

func (um *UserModel) New(email, userName, password, firstName, lastName, displayName string) (UserData, error) {
	var userID int
	var lastLogin sql.NullString
	var createdAt string
	var updatedAt string

	err := um.db.QueryRow(
		createUserQuery,
		email,
		userName,
		password,
		firstName,
		lastName,
		displayName,
	).Scan(&userID, &lastLogin, &createdAt, &updatedAt)

	if err != nil {
		if db.ConstraintFailed(err) {
			return UserData{}, db.WrapConstraintError(err)
		}

		return UserData{}, err
	}

	out := UserData{
		ID:          userID,
		Email:       email,
		Username:    userName,
		FirstName:   firstName,
		LastName:    lastName,
		DisplayName: displayName,
		LastLogin:   "", // last login null in the db
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return out, nil
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

//go:embed queries/update-last-login.sql
var updateLastLoginQuery string

func (um *UserModel) SetLastLogin(userID int) (lastLoginTime string, err error) {
	if userID == 0 {
		err := errors.New("Missing user ID")
		if util.InDevContext() {
			panic(err)
		} else {
			return "", err
		}
	}

	err = um.db.QueryRow(updateLastLoginQuery, userID).Scan(&lastLoginTime)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return "", UserNotFoundError{}
		} else {
			return "", err
		}
	}

	return lastLoginTime, err
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

func (um UserModel) EmailExists(email string) bool {
	return false
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
	var createdAt string
	var updatedAt string

	err := row.Scan(
		&id, &email, &userName, &password, &firstName, &lastName, &displayName,
		&lastLogin, &createdAt, &updatedAt)

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
		createdAt,
		updatedAt,
	}, nil
}
