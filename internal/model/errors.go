package model

type UserNotFoundError struct{}

func (_ UserNotFoundError) Error() string {
	return "User not found"
}

type MissingRequiredFilterData struct{}

func (_ MissingRequiredFilterData) Error() string {
	return "Missing required filter data"
}
