package dtypes

import (
	"time"

	"github.com/marcusprice/twitter-clone/internal/permissions"
)

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

type CommentInput struct {
	UserID          int
	PostID          int
	ParentCommentID int
	Content         string
	Image           string
}

type UserData struct {
	ID          int
	Email       string
	Username    string
	FirstName   string
	LastName    string
	DisplayName string
	Avatar      string
	Password    string
	LastLogin   string
	IsActive    int
	Role        permissions.Role
	CreatedAt   string
	UpdatedAt   string
}

type Author struct {
	Username    string
	DisplayName string
	Avatar      string
}

type PostData struct {
	Author        Author
	Retweeter     Retweeter
	ID            int
	UserID        int
	Content       string
	CommentCount  int
	LikeCount     int
	RetweetCount  int
	BookmarkCount int
	Impressions   int
	Image         string
	CreatedAt     string
	UpdatedAt     string
	Liked         int
	Retweeted     int
	Bookmarked    int
}

type CommentData struct {
	Author          Author
	Retweeter       Retweeter
	ID              int
	PostID          int
	UserID          int
	Depth           int
	ParentCommentID int
	Content         string
	LikeCount       int
	RetweetCount    int
	BookmarkCount   int
	Impressions     int
	Image           string
	CreatedAt       string
	UpdatedAt       string
}

type Retweeter struct {
	Username    string
	DisplayName string
}

type BookmarkData struct {
	BookmarkCreatedAt string
	ID                int
	Content           string
	Image             string
	LikeCount         int
	RetweetCount      int
	BookmarkCount     int
	Impressions       int
	CreatedAt         string
	UpdatedAt         string
	Author            Author
	Type              string
}

type IdentifierAlreadyExistsError struct{}

func (_ IdentifierAlreadyExistsError) Error() string {
	return "Username or email already exists"
}

type ModelResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	DoneReason         string    `json:"done_reason"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int       `json:"eval_duration"`
}
