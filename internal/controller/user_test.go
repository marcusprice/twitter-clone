package controller

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/constants"
	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/model"
	"github.com/marcusprice/twitter-clone/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUserController(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		user := NewUserController(db)
		tu.AssertNotNil(user)
		tu.AssertNil(user.id)

		defer tu.ShouldPanic()
		NewUserController(nil)
	})
}

func TestUserSet(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	user := User{}
	userInput := dtypes.UserInput{
		Email:       "estecat42069@yahoo.com",
		Username:    "estecat",
		FirstName:   "Esteban",
		LastName:    "Price",
		DisplayName: "hungry boy",
	}
	user.Set(nil, userInput)
	tu.AssertEqual("estecat42069@yahoo.com", user.Email)
	tu.AssertEqual("estecat", user.Username)
	tu.AssertEqual("Esteban", user.FirstName)
	tu.AssertEqual("Price", user.LastName)
	tu.AssertEqual("hungry boy", user.DisplayName)
}

func TestUserSetFromModel(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	user := User{}
	userDbData := model.UserData{
		Email:       "estecat42069@yahoo.com",
		Username:    "estecat",
		FirstName:   "Esteban",
		LastName:    "Price",
		DisplayName: "hungry boy",
		LastLogin:   "2024-04-12 11:37:46",
		CreatedAt:   "2024-04-12 11:37:46",
		UpdatedAt:   "2024-04-12 11:37:46",
	}
	user.setFromModel(userDbData)
	tu.AssertEqual("estecat42069@yahoo.com", user.Email)
	tu.AssertEqual("estecat", user.Username)
	tu.AssertEqual("Esteban", user.FirstName)
	tu.AssertEqual("Price", user.LastName)
	tu.AssertEqual("hungry boy", user.DisplayName)
	tu.AssertEqual("2024-04-12 11:37:46", user.LastLogin.Format(constants.TIME_LAYOUT))
	tu.AssertEqual("2024-04-12 11:37:46", user.CreatedAt.Format(constants.TIME_LAYOUT))
	tu.AssertEqual("2024-04-12 11:37:46", user.UpdatedAt.Format(constants.TIME_LAYOUT))
}

func TestUserCreate(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		user := initTestUser(db)
		user.Create("password")
		storedPassword := testutil.QueryUserPassword(user.ID(), db)
		tu.AssertTrue(validPasswordHash(storedPassword, "password"))
		tu.AssertEqual(1, user.ID())

		// query the user just created to verify it was recorded in the db
		queriedUser := initTestUser(db)
		queriedUser.ByID(user.ID())
		tu.AssertEqual(user.Email, queriedUser.Email)
		tu.AssertEqual(user.Username, queriedUser.Username)
		tu.AssertEqual(user.FirstName, queriedUser.FirstName)
		tu.AssertEqual(user.LastName, queriedUser.LastName)
		tu.AssertEqual(user.DisplayName, queriedUser.DisplayName)

		duplicateUser := initTestUser(db)
		err := duplicateUser.Create("password")
		tu.AssertErrorNotNil(err)

		err = user.Create("password")
		tu.AssertErrorNotNil(err)
	})
}

func TestUserAuthenticateAndSet(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		userModel := model.NewUserModel(db)
		user := User{
			model:       userModel,
			Username:    "estecat",
			Email:       "estecat42069@yahoo.com",
			FirstName:   "Esteban",
			LastName:    "Price",
			DisplayName: "hungry cat",
		}

		user.Create("password")
		authenticatedUser := User{
			model: userModel,
			Email: "estecat42069@yahoo.com",
		}
		authenticated, err := authenticatedUser.AuthenticateAndSet("password")
		tu.AssertTrue(authenticated)
		tu.AssertErrorNil(err)
		tu.AssertEqual(1, authenticatedUser.ID())
		tu.AssertEqual("estecat", authenticatedUser.Username)
		tu.AssertEqual("estecat42069@yahoo.com", authenticatedUser.Email)
		tu.AssertEqual("Esteban", authenticatedUser.FirstName)
		tu.AssertEqual("Price", authenticatedUser.LastName)
		tu.AssertEqual("hungry cat", authenticatedUser.DisplayName)

		wrongPwdUser := User{
			model: userModel,
			Email: "estecat42069@yahoo.com",
		}
		authenticated, err = wrongPwdUser.AuthenticateAndSet("wrong_password")
		tu.AssertFalse(authenticated)
		tu.AssertErrorNil(err)

		unknownIdentifierUser := User{
			model: userModel,
			Email: "whispers_from_wallphace@gobblegobble.com",
		}
		authenticated, err = unknownIdentifierUser.AuthenticateAndSet("wrong_password")
		tu.AssertFalse(authenticated)
		tu.AssertErrorNotNil(err)
	})
}

func TestUserSetLastLogin(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		user := initTestUser(db)
		user.Create("password")
		err := user.Login()
		tu.AssertErrorNil(err)

		user.id = nil
		err = user.Login()
		tu.AssertErrorNotNil(err)
	})
}

func TestUserByID(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		user := initTestUser(db)
		user.Create("password")
		newUserID := user.ID()

		userByID := &User{model: model.NewUserModel(db)}
		err := userByID.ByID(newUserID)

		tu.AssertErrorNil(err)
		tu.AssertEqual(user.Email, userByID.Email)
		tu.AssertEqual(user.Username, userByID.Username)
		tu.AssertEqual(user.FirstName, userByID.FirstName)
		tu.AssertEqual(user.LastName, userByID.LastName)
		tu.AssertEqual(user.DisplayName, userByID.DisplayName)

		errUser := &User{model: model.NewUserModel(db)}
		err = errUser.ByID(42069)
		tu.AssertErrorNotNil(err)
		err = errUser.ByID(0)
		tu.AssertErrorNotNil(err)
		err = errUser.ByID(-20)
		tu.AssertErrorNotNil(err)
	})
}

func TestUserFollow(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		user1 := NewUserController(db)
		user2 := NewUserController(db)
		user3 := NewUserController(db)
		user1.ByID(1)
		user2.ByID(2)
		user3.ByID(3)

		err := user1.Follow(user2.Username)
		tu.AssertErrorNil(err)

		err = user2.Follow(user1.Username)
		tu.AssertErrorNil(err)

		err = user1.Follow("idontexist")
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.Is(err, model.UserNotFoundError{}))

		err = user1.Follow(user1.Username)
		var constraintError dbutils.ConstraintError
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.CHECK_ERROR, constraintError.Constraint)

		userFollowNumRows := model.QueryUserFollowTableCount(db)
		tu.AssertEqual(2, userFollowNumRows)

		err = user3.Follow(user2.Username)
		tu.AssertErrorNil(err)

		userFollowNumRows = model.QueryUserFollowTableCount(db)
		tu.AssertEqual(3, userFollowNumRows)
	})
}

func validPasswordHash(storedPassword, password string) bool {
	valid := bcrypt.CompareHashAndPassword(
		[]byte(storedPassword), []byte(password)) == nil

	return valid
}

func initTestUser(db *sql.DB) *User {
	return &User{
		model:       model.NewUserModel(db),
		Username:    "estecat",
		Email:       "estecat42069@yahoo.com",
		FirstName:   "Esteban",
		LastName:    "Price",
		DisplayName: "hungry cat",
	}
}
