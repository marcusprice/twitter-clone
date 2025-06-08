package model

import (
	"database/sql"
	_ "embed"

	"github.com/marcusprice/twitter-clone/internal/db"
	"github.com/marcusprice/twitter-clone/internal/util"
)

//go:embed queries/create-user.sql
var createUserQuery string

//go:embed queries/check-unique-user.sql
var checkUniqueUserQuery string

type UserModel struct {
	db *sql.DB
}

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

func (um *UserModel) UsernameOrEmailExists(email, username string) (bool, error) {
	result, err := um.db.Query(checkUniqueUserQuery, email, username)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			return false, err
		}
	}
	defer result.Close()
	result.Next()
	var count int
	err = result.Scan(&count)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			return false, err
		}
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
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
