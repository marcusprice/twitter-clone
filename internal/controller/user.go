package controller

import (
	"database/sql"
	"errors"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/model"
	"github.com/marcusprice/twitter-clone/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	UserID      int    `json:"userID"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

type User struct {
	userModel   *model.UserModel
	userID      *int
	Email       string
	Username    string
	FirstName   string
	LastName    string
	DisplayName string
}

func (u *User) Set(userID *int, userInput UserInput) {
	u.userID = userID
	u.Email = userInput.Email
	u.Username = userInput.Username
	u.FirstName = userInput.FirstName
	u.LastName = userInput.LastName
	u.DisplayName = userInput.DisplayName
}

func (u User) Create(password string) (int, error) {
	if u.userID != nil {
		if util.InDevContext() {
			panic("userID should not exist while creating a user")
		} else {
			return -1, errors.New("UserID present on User.Create, should be nil")
		}
	}

	userModel := u.userModel
	userExists, err := userModel.UsernameOrEmailExists(u.Email, u.Username)
	if err != nil {
		return -1, err
	}

	if userExists {
		return -1, dtypes.IdentifierAlreadyExistsError{}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	userID, err := userModel.New(
		u.Email, u.Username, string(hashedPassword),
		u.FirstName, u.LastName, u.DisplayName)
	if err != nil {
		return -1, err
	}

	return userID, nil
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
