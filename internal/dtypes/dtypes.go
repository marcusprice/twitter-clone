package dtypes

type UserInput struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

type PostInput struct {
	UserID  int
	Content string
	Image   string
}

type UserData struct {
	ID          int
	Email       string
	Username    string
	FirstName   string
	LastName    string
	DisplayName string
	Password    string
	LastLogin   string
	IsActive    int
	CreatedAt   string
	UpdatedAt   string
}

type PostData struct {
	Author        PostAuthor
	Retweeter     PostRetweeter
	ID            int
	UserID        int
	Content       string
	LikeCount     int
	RetweetCount  int
	BookmarkCount int
	Impressions   int
	Image         string
	CreatedAt     string
	UpdatedAt     string
}

type PostAuthor struct {
	Username    string
	DisplayName string
	Avatar      string
}

type PostRetweeter struct {
	Username    string
	DisplayName string
}

type IdentifierAlreadyExistsError struct{}

func (_ IdentifierAlreadyExistsError) Error() string {
	return "Username or email already exists"
}
