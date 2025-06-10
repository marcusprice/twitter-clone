package controller

import (
	"database/sql"
	"testing"

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

func TestSetFromModel(t *testing.T) {
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
	tu.AssertEqual("2024-04-12 11:37:46", user.LastLogin.Format(TIME_LAYOUT))
	tu.AssertEqual("2024-04-12 11:37:46", user.CreatedAt.Format(TIME_LAYOUT))
	tu.AssertEqual("2024-04-12 11:37:46", user.UpdatedAt.Format(TIME_LAYOUT))
}

func TestCreate(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		userModel := model.NewUserModel(db)
		user := User{
			model:       userModel,
			Email:       "estecat42069@yahoo.com",
			Username:    "estecat",
			FirstName:   "Esteban",
			LastName:    "Price",
			DisplayName: "hungry boy",
		}
		user.Create("password")
		storedPassword := queryUserPassword(user.ID(), db)
		tu.AssertTrue(validPasswordHash(storedPassword, "password"))
		tu.AssertEqual(1, user.ID())

		// query the user just created to verify it was recorded in the db
		queriedUser := User{model: userModel}
		queriedUser.ByID(user.ID())
		tu.AssertEqual(user.Email, queriedUser.Email)
		tu.AssertEqual(user.Username, queriedUser.Username)
		tu.AssertEqual(user.FirstName, queriedUser.FirstName)
		tu.AssertEqual(user.LastName, queriedUser.LastName)
		tu.AssertEqual(user.DisplayName, queriedUser.DisplayName)

		duplicateUser := User{
			model:    userModel,
			Email:    user.Email,
			Username: user.Username,
		}

		_, err := duplicateUser.Create("password")
		tu.AssertErrorNotNil(err)
		// verify that panic is called when user.id set (returns error in prod)
		defer tu.ShouldPanic()
		user.Create("password")
	})
}

func TestPanicOnNewUserNilDB(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.ShouldPanic()
	NewUserController(nil)
}

func queryUserPassword(userID int, db *sql.DB) string {
	query := `SELECT password FROM User WHERE id = $1;`

	var password string
	err := db.QueryRow(query, userID).Scan(&password)
	if err != nil {
		panic(err)
	}

	return password
}

func validPasswordHash(storedPassword, password string) bool {
	valid := bcrypt.CompareHashAndPassword(
		[]byte(storedPassword), []byte(password)) == nil

	return valid
}
