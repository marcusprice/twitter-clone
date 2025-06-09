package controller

import (
	"database/sql"
	"errors"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/model"
	"github.com/marcusprice/twitter-clone/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userModel   *model.UserModel
	UserID      *int
	Email       string
	Username    string
	FirstName   string
	LastName    string
	DisplayName string
	LastLogin   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) Set(userID *int, userInput dtypes.UserInput) {
	u.UserID = userID
	u.Email = userInput.Email
	u.Username = userInput.Username
	u.FirstName = userInput.FirstName
	u.LastName = userInput.LastName
	u.DisplayName = userInput.DisplayName
}

func (u *User) setFromModel(userData model.UserData) {
	u.UserID = &userData.ID
	u.Email = userData.Email
	u.Username = userData.Username
	u.FirstName = userData.FirstName
	u.LastName = userData.LastName
	u.DisplayName = userData.DisplayName

	lastLogin, err := time.Parse(TIME_LAYOUT, userData.LastLogin)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			u.LastLogin = time.Time{}
		}
	}

	createdAt, err := time.Parse(TIME_LAYOUT, userData.CreatedAt)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			u.CreatedAt = time.Time{}
		}
	}

	updatedAt, err := time.Parse(TIME_LAYOUT, userData.UpdatedAt)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			u.UpdatedAt = time.Time{}
		}
	}

	u.LastLogin = lastLogin
	u.CreatedAt = createdAt
	u.UpdatedAt = updatedAt
}

func (u User) Create(password string) (int, error) {
	if u.UserID != nil {
		if util.InDevContext() {
			panic("userID should not exist while creating a user")
		} else {
			return -1, errors.New("UserID present on User.Create, should be nil")
		}
	}

	userModel := u.userModel
	userExists, err := userModel.UsernameOrEmailExists(u.Email, u.Username)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			return -1, err
		}
	}

	if userExists {
		return -1, dtypes.IdentifierAlreadyExistsError{}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			return -1, err
		}
	}

	userID, err := userModel.New(
		u.Email, u.Username, string(hashedPassword),
		u.FirstName, u.LastName, u.DisplayName)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (u *User) AuthenticateAndSet(pwd string) (authenticated bool, err error) {
	userModel := u.userModel
	userData, err := userModel.GetByIdentifier(u.Email, u.Username)
	if err != nil {
		return false, err
	}

	valid := bcrypt.CompareHashAndPassword(
		[]byte(userData.Password), []byte(pwd)) == nil

	if !valid {
		return false, nil
	}

	u.setFromModel(userData)

	return true, nil
}

func (u *User) SetLastLogin() error {
	if u.UserID == nil {
		err := errors.New("trying to update a user login without ID")
		if util.InDevContext() {
			panic(err)
		} else {
			return (err)
		}
	}

	return u.userModel.SetLastLogin(*u.UserID)
}

func (u *User) ByID(userID int) error {
	userData, err := u.userModel.GetByID(userID)
	if err != nil {
		return err
	}

	u.setFromModel(userData)

	return nil
}

func NewUserController(dbConn *sql.DB) *User {
	if dbConn == nil {
		panic("db conn cannot be nil")
	}

	userModel := model.NewUserModel(dbConn)
	return &User{
		userModel: userModel,
	}
}
