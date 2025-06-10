package model

import (
	"database/sql"
	"testing"

	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestNewUserModel(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		userModel := NewUserModel(db)
		tu.AssertNotNil(userModel)

		defer tu.ShouldPanic()
		NewUserModel(nil)
	})
}

// func TestUserModelNew(t *testing.T) {
// 	testutil.WithTestDB(t, func(db *sql.DB) {
// 		userModel := UserModel{db: db}

// 	})
// }
