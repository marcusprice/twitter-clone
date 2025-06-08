package dtypes

type UserInput struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

type IdentifierAlreadyExistsError struct{}

func (_ IdentifierAlreadyExistsError) Error() string {
	return "Username or email already exists"
}
