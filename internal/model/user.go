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

func (um *UserModel) New(email, userName, password, firstName, lastName, displayName string) (int, error) {
	result, err := um.db.Exec(
		createUserQuery,
		email,
		userName,
		password,
		firstName,
		lastName,
		displayName,
	)

	if err != nil {
		if db.ConstraintFailed(err) {
			return -1, db.WrapConstraintError(err)
		}

		return -1, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			return -1, err
		}
	}

	return int(userID), nil
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

func (um *UserModel) OneOrNone(email, username string) (UserData, error) {
	if email == "" && username == "" {
		return UserData{}, MissingRequiredFilterData{}
	}

	var row *sql.Row
	query := selectUserBaseQuery
	if email != "" && username != "" {
		query += "WHERE email = $1 AND user_name = $2;"
		row = um.db.QueryRow(query, email, username)
	} else if email != "" {
		query += "WHERE email = $1;"
		row = um.db.QueryRow(query, email)
	} else {
		query += "WHERE user_name = $1;"
		row = um.db.QueryRow(query, username)
	}

	return parseUserQueryRow(row)
}

//go:embed queries/update-last-login.sql
var updateLastLoginQuery string

func (um *UserModel) SetLastLogin(userID int) error {
	if userID == 0 {
		err := errors.New("Missing user ID")
		if util.InDevContext() {
			panic(err)
		} else {
			return err
		}
	}

	result, err := um.db.Exec(updateLastLoginQuery, userID)
	if err != nil {
		return err
	}
	rowsUpdated, err := result.RowsAffected()
	if int(rowsUpdated) == 0 {
		return UserNotFoundError{}
	}

	return err
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
	var lastLogin string
	var createdAt string
	var updatedAt string

	err := row.Scan(
		&id, &email, &userName, &password, &firstName, &lastName, &displayName,
		&lastLogin, &createdAt, &updatedAt)

	if err != nil {
		return UserData{}, UserNotFoundError{}
	}

	return UserData{
		id,
		email,
		userName,
		firstName,
		lastName,
		displayName,
		password,
		lastLogin,
		createdAt,
		updatedAt,
	}, nil
}
