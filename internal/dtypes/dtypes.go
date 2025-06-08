package dtypes

type IdentifierAlreadyExistsError struct{}

func (_ IdentifierAlreadyExistsError) Error() string {
	return "Username or email already exists"
}
