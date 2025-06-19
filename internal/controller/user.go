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
	model       *model.UserModel
	id          *int // TODO: change this to a regular in, use 0 value as null check
	Email       string
	Username    string
	FirstName   string
	LastName    string
	DisplayName string
	IsActive    bool
	LastLogin   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u User) ID() int {
	return *u.id
}

func (u *User) Set(userID *int, userInput dtypes.UserInput) {
	u.id = userID
	u.Email = userInput.Email
	u.Username = userInput.Username
	u.FirstName = userInput.FirstName
	u.LastName = userInput.LastName
	u.DisplayName = userInput.DisplayName
}

func (u *User) setFromModel(userData dtypes.UserData) {
	u.id = &userData.ID
	u.Email = userData.Email
	u.Username = userData.Username
	u.FirstName = userData.FirstName
	u.LastName = userData.LastName
	u.DisplayName = userData.DisplayName
	u.IsActive = userData.IsActive != 0
	u.LastLogin = util.ParseTime(userData.LastLogin)
	u.CreatedAt = util.ParseTime(userData.CreatedAt)
	u.UpdatedAt = util.ParseTime(userData.UpdatedAt)
}

func (u *User) Create(password string) error {
	if u.id != nil {
		if util.InDevContext() {
			panic("User.id should not exist while creating a user")
		} else {
			return errors.New("User.id present on User.Create, should be nil")
		}
	}

	model := u.model
	userExists, err := model.UsernameOrEmailExists(u.Email, u.Username)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			return err
		}
	}

	if userExists {
		return dtypes.IdentifierAlreadyExistsError{}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if util.InDevContext() {
			panic(err)
		} else {
			return err
		}
	}

	userInput := dtypes.UserInput{
		Email:       u.Email,
		Username:    u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		DisplayName: u.DisplayName,
		Password:    string(hashedPassword),
	}

	userData, err := model.New(userInput)

	if err != nil {
		return err
	}

	u.setFromModel(userData)
	return nil
}

func (u *User) Follow(followeeUsername string) error {
	followeeData, err := u.model.GetByIdentifier("", followeeUsername)
	if err != nil {
		return err
	}

	err = u.model.Follow(u.ID(), followeeData.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) UnFollow(followeeUsername string) error {
	if u.Username == followeeUsername {
		return errors.New("cannot unfollow yourself")
	}

	followeeData, err := u.model.GetByIdentifier("", followeeUsername)
	if err != nil {
		return err
	}

	err = u.model.UnFollow(u.ID(), followeeData.ID)
	if err != nil {
		return err
	}

	return nil

}

func (u *User) AuthenticateAndSet(pwd string) (authenticated bool, err error) {
	model := u.model
	userData, err := model.GetByIdentifier(u.Email, u.Username)
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

func (user *User) Login() error {
	if user.id == nil {
		err := errors.New("trying to update a user login without ID")
		if util.InDevContext() {
			panic(err)
		} else {
			return (err)
		}
	}
	lastLoginTime, isActive, err := user.model.Login(user.ID())
	user.LastLogin = util.ParseTime(lastLoginTime)
	user.IsActive = isActive != 0
	return err
}

func (u *User) ByID(userID int) error {
	userData, err := u.model.GetByID(userID)
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

	model := model.NewUserModel(dbConn)
	return &User{
		model: model,
	}
}
