package api

import "time"

type UserPayload struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
}

type PostPayload struct {
	ID            int           `json:"postID"`
	Content       string        `json:"content"`
	LikeCount     int           `json:"likeCount"`
	RetweetCount  int           `json:"retweetCount"`
	BookmarkCount int           `json:"bookmarkCount"`
	Impressions   int           `json:"impressions"`
	Image         string        `json:"image"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	Author        AuthorPayload `json:"author"`
}

type AuthorPayload struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
}
