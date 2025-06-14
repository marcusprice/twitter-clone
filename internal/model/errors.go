package model

type MissingRequiredFilterData struct{}

func (_ MissingRequiredFilterData) Error() string {
	return "Missing required filter data"
}

type UserNotFoundError struct{}

func (_ UserNotFoundError) Error() string {
	return "User not found"
}

type PostNotFoundError struct{}

func (_ PostNotFoundError) Error() string {
	return "Post not found"
}
