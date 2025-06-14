package model

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testutil"
	"github.com/marcusprice/twitter-clone/internal/util"
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

func TestUserModelNew(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		userModel := UserModel{db: db}
		esteban := dtypes.UserInput{
			Email:       "estecat42069@yahoo.com",
			Username:    "estecat",
			FirstName:   "Esteban",
			LastName:    "Price",
			DisplayName: "Hungry Boy",
			Password:    "password",
		}

		marcus := dtypes.UserInput{
			Email:       "whispers_from_wallphace@hotmail.com",
			Username:    "suhdude",
			FirstName:   "Marcus",
			LastName:    "Price",
			DisplayName: "gobble gobble",
			Password:    "password",
		}

		beforeCreate := time.Now().UTC().Add(-1 * time.Minute)
		estebanUserData, err := userModel.New(esteban)
		afterCreate := time.Now().UTC().Add(time.Minute)
		tu.AssertErrorNil(err)
		tu.AssertEqual(1, estebanUserData.ID)
		tu.AssertEqual("estecat42069@yahoo.com", estebanUserData.Email)
		tu.AssertEqual("estecat", estebanUserData.Username)
		tu.AssertEqual("Esteban", estebanUserData.FirstName)
		tu.AssertEqual("Price", estebanUserData.LastName)
		tu.AssertEqual("", estebanUserData.Password)

		queriedUser := queryUser(estebanUserData.ID, db)
		tu.AssertEqual(1, queriedUser.ID)
		tu.AssertEqual("estecat42069@yahoo.com", queriedUser.Email)
		tu.AssertEqual("estecat", queriedUser.Username)
		tu.AssertEqual("Esteban", queriedUser.FirstName)
		tu.AssertEqual("Price", queriedUser.LastName)
		tu.AssertEqual("password", queriedUser.Password)
		tu.AssertEqual("", queriedUser.LastLogin)
		tu.AssertTrue(util.ParseTime(queriedUser.CreatedAt).After(beforeCreate))
		tu.AssertTrue(util.ParseTime(queriedUser.CreatedAt).Before(afterCreate))
		tu.AssertTrue(util.ParseTime(queriedUser.UpdatedAt).After(beforeCreate))
		tu.AssertTrue(util.ParseTime(queriedUser.UpdatedAt).Before(afterCreate))

		marcusUserData, err := userModel.New(marcus)
		marcusQueried := queryUser(marcusUserData.ID, db)
		tu.AssertErrorNil(err)
		tu.AssertEqual(2, marcusUserData.ID)
		tu.AssertEqual(2, marcusQueried.ID)

		dup := dtypes.UserInput{
			Email:       "whispers_from_wallphace@hotmail.com",
			Username:    "suhdude",
			FirstName:   "Marcus",
			LastName:    "Price",
			DisplayName: "gobble gobble",
			Password:    "password",
		}

		_, err = userModel.New(dup)
		var constraintError dbutils.ConstraintError
		errors.As(err, &constraintError)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(dbutils.IsConstraintError(err))
		tu.AssertEqual(dbutils.UNIQUE_ERROR, constraintError.Constraint)

		fieldMissing := dtypes.UserInput{
			Email:       "suhdude42069@yahoo.com",
			FirstName:   "Marcus",
			LastName:    "Price",
			DisplayName: "gobble gobble",
			Password:    "password",
		}

		_, err = userModel.New(fieldMissing)
		errors.As(err, &constraintError)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(dbutils.IsConstraintError(err))
		// go will always populate it with an empty string because of the
		// UserInput default values, so instead of NOT_NULL_ERROR it will fail
		// the check to prevent empty values((length(trim(email)) > 0))
		tu.AssertEqual(dbutils.CHECK_ERROR, constraintError.Constraint)
	})
}

func TestUserGetById(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		userModel := UserModel{db: db}
		user := dtypes.UserInput{
			Email:       "estecat42069@yahoo.com",
			Username:    "estecat",
			FirstName:   "Esteban",
			LastName:    "Price",
			DisplayName: "Hungry Boy",
			Password:    "password",
		}
		id := insertUser(user, db)
		userData, err := userModel.GetByID(id)
		tu.AssertErrorNil(err)
		tu.AssertEqual("estecat42069@yahoo.com", userData.Email)
		tu.AssertEqual("estecat", userData.Username)
		tu.AssertEqual("Esteban", userData.FirstName)
		tu.AssertEqual("Price", userData.LastName)
		tu.AssertEqual("Hungry Boy", userData.DisplayName)
		tu.AssertEqual("password", userData.Password)
		tu.AssertEqual("", userData.LastLogin)

		var unsetInt int
		_, err = userModel.GetByID(unsetInt)
		tu.AssertErrorNotNil(err)
	})
}

func TestUserGetByIdentifier(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		userModel := UserModel{db: db}
		user := dtypes.UserInput{
			Email:       "estecat42069@yahoo.com",
			Username:    "estecat",
			FirstName:   "Esteban",
			LastName:    "Price",
			DisplayName: "Hungry Boy",
			Password:    "password",
		}
		insertUser(user, db)
		userByEmail, err := userModel.GetByIdentifier(user.Email, "")
		tu.AssertErrorNil(err)
		tu.AssertEqual("estecat42069@yahoo.com", userByEmail.Email)
		tu.AssertEqual("estecat", userByEmail.Username)
		tu.AssertEqual("Esteban", userByEmail.FirstName)
		tu.AssertEqual("Price", userByEmail.LastName)
		tu.AssertEqual("Hungry Boy", userByEmail.DisplayName)
		tu.AssertEqual("password", userByEmail.Password)

		userByUsername, err := userModel.GetByIdentifier("", user.Username)
		tu.AssertErrorNil(err)
		tu.AssertEqual("estecat42069@yahoo.com", userByUsername.Email)
		tu.AssertEqual("estecat", userByUsername.Username)
		tu.AssertEqual("Esteban", userByUsername.FirstName)
		tu.AssertEqual("Price", userByUsername.LastName)
		tu.AssertEqual("Hungry Boy", userByUsername.DisplayName)
		tu.AssertEqual("password", userByUsername.Password)
		tu.AssertEqual("", userByUsername.LastLogin)

		_, err = userModel.GetByIdentifier("not-a-person@yahoo.com", "")
		var notFoundError UserNotFoundError
		tu.AssertTrue(errors.As(err, &notFoundError))
		tu.AssertErrorNotNil(err)

		beforeUpdate := time.Now().UTC().Add(-1 * time.Minute)
		updateLastLogin(userByEmail.ID, db)
		afterUpdate := time.Now().UTC().Add(time.Minute)

		userByUsername, _ = userModel.GetByIdentifier("", user.Username)
		lastLogin := util.ParseTime(userByUsername.LastLogin)
		tu.AssertTrue(lastLogin.After(beforeUpdate))
		tu.AssertTrue(lastLogin.Before(afterUpdate))

		_, err = userModel.GetByIdentifier("", "")
		var missingRequiredFilterData MissingRequiredFilterData
		tu.AssertTrue(errors.As(err, &missingRequiredFilterData))
		tu.AssertErrorNotNil(err)
	})
}

func TestUserUsernameOrEmailExists(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		userModel := UserModel{db}
		user := dtypes.UserInput{
			Email:       "estecat42069@yahoo.com",
			Username:    "estecat",
			FirstName:   "Esteban",
			LastName:    "Price",
			DisplayName: "Hungry Boy",
			Password:    "password",
		}
		insertUser(user, db)
		exists, err := userModel.UsernameOrEmailExists("estecat42069@yahoo.com", "")
		tu.AssertErrorNil(err)
		tu.AssertTrue(exists)

		exists, err = userModel.UsernameOrEmailExists("", "estecat")
		tu.AssertErrorNil(err)
		tu.AssertTrue(exists)

		exists, err = userModel.UsernameOrEmailExists("yayayaya@hotmail.com", "")
		tu.AssertErrorNil(err)
		tu.AssertFalse(exists)

		exists, err = userModel.UsernameOrEmailExists("", "whispers_from_wallphace")
		tu.AssertErrorNil(err)
		tu.AssertFalse(exists)

		exists, err = userModel.UsernameOrEmailExists("yayayaya@hotmail.com", "whispers_from_wallphace")
		tu.AssertErrorNil(err)
		tu.AssertFalse(exists)
	})
}

func TestUserLogin(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		userModel := UserModel{db}
		user := dtypes.UserInput{
			Email:       "estecat42069@yahoo.com",
			Username:    "estecat",
			FirstName:   "Esteban",
			LastName:    "Price",
			DisplayName: "Hungry Boy",
			Password:    "password",
		}

		userID := insertUser(user, db)
		beforeUpdate := time.Now().UTC().Add(-1 * time.Minute)
		methodReturnTimestamp, isActive, err := userModel.Login(userID)
		afterUpdate := time.Now().UTC().Add(time.Minute)

		userData := queryUser(userID, db)
		timestamp := util.ParseTime(methodReturnTimestamp)
		lastLogin := util.ParseTime(userData.LastLogin)
		tu.AssertTrue(timestamp.After(beforeUpdate))
		tu.AssertTrue(timestamp.Before(afterUpdate))
		tu.AssertTrue(lastLogin.After(beforeUpdate))
		tu.AssertTrue(lastLogin.Before(afterUpdate))
		tu.AssertEqual(1, isActive)

		var unsetID int
		_, _, err = userModel.Login(unsetID)
		tu.AssertErrorNotNil(err)

		_, _, err = userModel.Login(42069)
		var userNotFoundError UserNotFoundError
		tu.AssertTrue(errors.As(err, &userNotFoundError))
		tu.AssertErrorNotNil(err)
	})
}

func queryUser(userID int, db *sql.DB) UserData {
	query := `
	SELECT 
		id,
		email,
		user_name,
		first_name,
		last_name,
		display_name,
		password,
		last_login,
		created_at,
		updated_at
	FROM User
	WHERE id = $1`

	var id int
	var email string
	var username string
	var firstName string
	var lastName string
	var displayName string
	var password string
	var lastLogin sql.NullString
	var createdAt string
	var updatedAt string

	err := db.QueryRow(query, userID).
		Scan(
			&id, &email, &username, &firstName, &lastName,
			&displayName, &password, &lastLogin, &createdAt, &updatedAt,
		)

	if err != nil {
		panic(err)
	}

	lastLoginString := ""
	if lastLogin.Valid {
		lastLoginString = lastLogin.String
	} else {
		lastLoginString = ""
	}

	return UserData{
		ID:          id,
		Email:       email,
		Username:    username,
		FirstName:   firstName,
		LastName:    lastName,
		DisplayName: displayName,
		Password:    password,
		LastLogin:   lastLoginString,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func insertUser(userInput dtypes.UserInput, db *sql.DB) int {
	query := `
	INSERT INTO User 
		(email, user_name, first_name, last_name, display_name, password, is_active)
	Values
		($1, $2, $3, $4, $5, $6, 0)
	RETURNING id;
	`

	var userID int
	err := db.QueryRow(query, userInput.Email, userInput.Username,
		userInput.FirstName, userInput.LastName,
		userInput.DisplayName, userInput.Password).Scan(&userID)

	if err != nil {
		panic(err)
	}

	return userID
}

func updateLastLogin(userID int, db *sql.DB) {
	query := `
	UPDATE User 
	SET last_login = CURRENT_TIMESTAMP
	WHERE id = $1;
	`

	_, err := db.Exec(query, userID)
	if err != nil {
		panic(err)
	}
}
