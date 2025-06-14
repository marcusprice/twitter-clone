package model

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
