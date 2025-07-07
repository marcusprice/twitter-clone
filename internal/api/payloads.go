package api

import (
	"time"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/util"
)

type TimelinePayload struct {
	Posts          []TimelinePostPayload `json:"posts"`
	HasMore        bool                  `json:"hasMore"`
	PostsRemaining int                   `json:"postsRemaining"`
}

type UserPayload struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
}

type AuthorPayload struct {
	Username        string           `json:"username"`
	DisplayName     string           `json:"displayName"`
	Avatar          string           `json:"avatar"`
	Bio             string           `json:"bio"`
	FollowerCount   int              `json:"followerCount"`
	FollowingCount  int              `json:"followingCount"`
	ViewerFollowing bool             `json:"viewerFollowing"`
	MutalFollowers  []*AuthorPayload `json:"mutualFollowers"`
}

func generateAuthorPayload(author dtypes.Author) AuthorPayload {
	authorPayload := AuthorPayload{
		Username:        author.Username,
		DisplayName:     author.DisplayName,
		Avatar:          author.Avatar,
		Bio:             author.Bio,
		FollowerCount:   author.FollowerCount,
		FollowingCount:  author.FollowingCount,
		ViewerFollowing: author.ViewerFollowing,
	}

	if authorPayload.Avatar != "" {
		authorPayload.Avatar = getUploadPath(authorPayload.Avatar)
	}

	return authorPayload
}

type RetweeterPayload struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type PostPayload struct {
	ID                   int           `json:"postID"`
	Content              string        `json:"content"`
	CommentCount         int           `json:"commentCount"`
	LikeCount            int           `json:"likeCount"`
	RetweetCount         int           `json:"retweetCount"`
	BookmarkCount        int           `json:"bookmarkCount"`
	Impressions          int           `json:"impressions"`
	Image                string        `json:"image"`
	CreatedAt            time.Time     `json:"createdAt"`
	UpdatedAt            time.Time     `json:"updatedAt"`
	Author               AuthorPayload `json:"author"`
	IsRetweet            bool          `json:"isRetweet"`
	RetweeterUsername    string        `json:"retweeterUsername"`
	RetweeterDisplayName string        `json:"retweeterDisplayName"`
	Liked                bool          `json:"liked"`
	Retweeted            bool          `json:"retweeted"`
	Bookmarked           bool          `json:"bookmarked"`
}

func generatePostPayload(post *controller.Post) PostPayload {
	if post.Author.Avatar != "" {
		post.Author.Avatar = getUploadPath(post.Author.Avatar)
	}

	author := AuthorPayload{
		Username:    post.Author.Username,
		DisplayName: post.Author.DisplayName,
		Avatar:      post.Author.Avatar,
	}

	if post.Image != "" {
		post.Image = getUploadPath(post.Image)
	}

	return PostPayload{
		ID:                   post.ID,
		Content:              post.Content,
		CommentCount:         post.CommentCount,
		LikeCount:            post.LikeCount,
		RetweetCount:         post.RetweetCount,
		BookmarkCount:        post.BookmarkCount,
		Impressions:          post.Impressions,
		Image:                post.Image,
		CreatedAt:            post.CreatedAt,
		UpdatedAt:            post.UpdatedAt,
		Author:               author,
		IsRetweet:            post.Retweeter.Username != "",
		RetweeterUsername:    post.Retweeter.Username,
		RetweeterDisplayName: post.Retweeter.DisplayName,
		Liked:                post.Liked,
		Retweeted:            post.Retweeted,
		Bookmarked:           post.Bookmarked,
	}
}

type TimelinePostPayload struct {
	Type                        string    `json:"type"`
	IsRetweet                   bool      `json:"isRetweet"`
	ID                          int       `json:"id"`
	Content                     string    `json:"content"`
	CommentCount                int       `json:"commentCount"`
	LikeCount                   int       `json:"likeCount"`
	RetweetCount                int       `json:"retweetCount"`
	BookmarkCount               int       `json:"bookmarkCount"`
	Impressions                 int       `json:"impressions"`
	Image                       string    `json:"image"`
	CreatedAt                   time.Time `json:"createdAt"`
	UpdatedAt                   time.Time `json:"updatedAt"`
	ViewerLiked                 int       `json:"viewerLiked"`
	ViewerRetweeted             int       `json:"viewerRetweeted"`
	ViewerBookmarked            int       `json:"viewerBookmarked"`
	ParentPostID                int       `json:"parentPostID"`
	ParentPostAuthorUsername    string    `json:"parentPostAuthorUsername"`
	ParentCommentID             int       `json:"parentCommentID"`
	ParentCommentAuthorUsername string    `json:"parentCommentAuthorUsername"`

	Author    AuthorPayload    `json:"author"`
	Retweeter RetweeterPayload `json:"retweeter"`
}

func generateTimelinePostPayload(timelinePostData dtypes.TimelinePostData) TimelinePostPayload {
	authorPayload := AuthorPayload{
		Username:       timelinePostData.Author.Username,
		DisplayName:    timelinePostData.Author.DisplayName,
		Avatar:         timelinePostData.Author.Avatar,
		Bio:            timelinePostData.Author.Bio,
		FollowerCount:  timelinePostData.Author.FollowerCount,
		FollowingCount: timelinePostData.Author.FollowingCount,
	}
	retweeterPayload := RetweeterPayload(timelinePostData.Retweeter)

	if authorPayload.Avatar != "" {
		authorPayload.Avatar = getUploadPath(authorPayload.Avatar)
	}

	payload := TimelinePostPayload{
		Type:                        timelinePostData.Type,
		ID:                          timelinePostData.ID,
		Content:                     timelinePostData.Content,
		CommentCount:                timelinePostData.CommentCount,
		LikeCount:                   timelinePostData.LikeCount,
		RetweetCount:                timelinePostData.RetweetCount,
		BookmarkCount:               timelinePostData.BookmarkCount,
		Impressions:                 timelinePostData.Impressions,
		Image:                       timelinePostData.Image,
		CreatedAt:                   util.ParseTime(timelinePostData.CreatedAt),
		UpdatedAt:                   util.ParseTime(timelinePostData.UpdatedAt),
		ViewerLiked:                 timelinePostData.ViewerLiked,
		ViewerRetweeted:             timelinePostData.ViewerRetweeted,
		ViewerBookmarked:            timelinePostData.ViewerBookmarked,
		ParentPostID:                timelinePostData.ParentPostID,
		ParentPostAuthorUsername:    timelinePostData.ParentPostAuthorUsername,
		ParentCommentID:             timelinePostData.ParentCommentID,
		ParentCommentAuthorUsername: timelinePostData.ParentCommentAuthorUsername,
		Author:                      authorPayload,
		Retweeter:                   retweeterPayload,
		IsRetweet:                   retweeterPayload.Username != "",
	}

	if payload.Image != "" {
		payload.Image = getUploadPath(payload.Image)
	}

	return payload
}

type BookmarkPayload struct {
	BookmarkCreatedAt string        `json:"bookmarkCreatedAt"`
	ID                int           `json:"id"`
	Content           string        `json:"content"`
	Image             string        `json:"image"`
	LikeCount         int           `json:"likeCount"`
	RetweetCount      int           `json:"retweetCount"`
	BookmarkCount     int           `json:"bookmarkCount"`
	Impressions       int           `json:"impressions"`
	CreatedAt         string        `json:"createdAt"`
	UpdatedAt         string        `json:"updatedAt"`
	Author            AuthorPayload `json:"author"`
	Type              string        `json:"type"`
}

type BookmarkResponsePayload struct {
	Bookmarks          []BookmarkPayload `json:"bookmarks"`
	HasMore            bool              `json:"hasMore"`
	BookmarksRemaining int               `json:"bookmarksRemaining"`
}

func generateBookmarkPayload(bookmarkData []dtypes.BookmarkData, bookmarksRemaining int) BookmarkResponsePayload {
	var bookmarks []BookmarkPayload
	for _, bookmark := range bookmarkData {
		authorPayload := AuthorPayload{
			Username:       bookmark.Author.Username,
			DisplayName:    bookmark.Author.DisplayName,
			Avatar:         bookmark.Author.Avatar,
			Bio:            bookmark.Author.Bio,
			FollowerCount:  bookmark.Author.FollowerCount,
			FollowingCount: bookmark.Author.FollowingCount,
		}

		bp := BookmarkPayload{
			BookmarkCreatedAt: bookmark.BookmarkCreatedAt,
			ID:                bookmark.ID,
			Content:           bookmark.Content,
			Image:             bookmark.Image,
			LikeCount:         bookmark.LikeCount,
			RetweetCount:      bookmark.RetweetCount,
			BookmarkCount:     bookmark.BookmarkCount,
			Impressions:       bookmark.Impressions,
			CreatedAt:         bookmark.CreatedAt,
			UpdatedAt:         bookmark.UpdatedAt,
			Author:            authorPayload,
			Type:              bookmark.Type,
		}
		bookmarks = append(bookmarks, bp)
	}

	return BookmarkResponsePayload{
		Bookmarks:          bookmarks,
		HasMore:            bookmarksRemaining > 0,
		BookmarksRemaining: bookmarksRemaining,
	}
}

type CommentPayload struct {
	ID                   int           `json:"commentID"`
	PostID               int           `json:"postID"`
	ParentCommentID      int           `json:"parentCommentID"`
	Content              string        `json:"content"`
	LikeCount            int           `json:"likeCount"`
	RetweetCount         int           `json:"retweetCount"`
	BookmarkCount        int           `json:"bookmarkCount"`
	Impressions          int           `json:"impressions"`
	Image                string        `json:"image"`
	CreatedAt            time.Time     `json:"createdAt"`
	UpdatedAt            time.Time     `json:"updatedAt"`
	Author               AuthorPayload `json:"author"`
	IsRetweet            bool          `json:"isRetweet"`
	RetweeterUsername    string        `json:"retweeterUsername"`
	RetweeterDisplayName string        `json:"retweeterDisplayName"`
}

func generateCommentPayload(comment *controller.Comment) *CommentPayload {
	if comment.Image != "" {
		comment.Image = getUploadPath(comment.Image)
	}

	if comment.Author.Avatar != "" {
		comment.Author.Avatar = getUploadPath(comment.Author.Avatar)
	}
	author := AuthorPayload{
		Username:    comment.Author.Username,
		DisplayName: comment.Author.DisplayName,
		Avatar:      comment.Author.Avatar,
	}

	return &CommentPayload{
		ID:                   comment.ID,
		PostID:               comment.PostID,
		ParentCommentID:      comment.ParentCommentID,
		Content:              comment.Content,
		LikeCount:            comment.LikeCount,
		RetweetCount:         comment.RetweetCount,
		BookmarkCount:        comment.BookmarkCount,
		Impressions:          comment.Impressions,
		Image:                comment.Image,
		CreatedAt:            comment.CreatedAt,
		UpdatedAt:            comment.UpdatedAt,
		IsRetweet:            comment.IsRetweet,
		RetweeterUsername:    comment.RetweeterUsername,
		RetweeterDisplayName: comment.RetweeterDisplayName,
		Author:               author,
	}
}

type CommentFromPostPayload struct {
	ID              int                       `json:"commentID"`
	PostID          int                       `json:"postID"`
	ParentCommentID int                       `json:"parentCommentID"`
	Content         string                    `json:"content"`
	LikeCount       int                       `json:"likeCount"`
	RetweetCount    int                       `json:"retweetCount"`
	BookmarkCount   int                       `json:"bookmarkCount"`
	Impressions     int                       `json:"impressions"`
	Image           string                    `json:"image"`
	CreatedAt       time.Time                 `json:"createdAt"`
	UpdatedAt       time.Time                 `json:"updatedAt"`
	Author          AuthorPayload             `json:"author"`
	Replies         []*CommentFromPostPayload `json:"replies"`
}

type PostAndCommentsPayload struct {
	ID            int                       `json:"postID"`
	Content       string                    `json:"content"`
	CommentCount  int                       `json:"commentCount"`
	LikeCount     int                       `json:"likeCount"`
	RetweetCount  int                       `json:"retweetCount"`
	BookmarkCount int                       `json:"bookmarkCount"`
	Impressions   int                       `json:"impressions"`
	Image         string                    `json:"image"`
	CreatedAt     time.Time                 `json:"createdAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
	Author        AuthorPayload             `json:"author"`
	Liked         bool                      `json:"liked"`
	Comments      []*CommentFromPostPayload `json:"comments"`
}

func generatePostAndCommentsPayload(post *controller.Post) PostAndCommentsPayload {
	postAndCommentsPayload := PostAndCommentsPayload{}
	postAndCommentsPayload.Comments = []*CommentFromPostPayload{}
	for _, comment := range post.Comments {
		commentPayload := &CommentFromPostPayload{}
		repliesPayload := []*CommentFromPostPayload{}

		if comment.Image != "" {
			comment.Image = getUploadPath(comment.Image)
		}

		if comment.Author.Avatar != "" {
			comment.Author.Avatar = getUploadPath(comment.Author.Avatar)
		}

		for _, reply := range comment.Replies {
			authorPayload := AuthorPayload{
				Username:       reply.Author.Username,
				DisplayName:    reply.Author.DisplayName,
				Avatar:         reply.Author.Avatar,
				Bio:            reply.Author.Bio,
				FollowerCount:  reply.Author.FollowerCount,
				FollowingCount: reply.Author.FollowingCount,
			}
			replyPayload := &CommentFromPostPayload{}
			replyPayload.ID = reply.ID
			replyPayload.PostID = reply.PostID
			replyPayload.ParentCommentID = reply.ParentCommentID
			replyPayload.Content = reply.Content
			replyPayload.LikeCount = reply.LikeCount
			replyPayload.RetweetCount = reply.RetweetCount
			replyPayload.BookmarkCount = reply.BookmarkCount
			replyPayload.Impressions = reply.Impressions
			replyPayload.Image = reply.Image
			replyPayload.CreatedAt = reply.CreatedAt
			replyPayload.UpdatedAt = reply.UpdatedAt
			replyPayload.Author = authorPayload
			repliesPayload = append(repliesPayload, replyPayload)
		}

		authorPayload := AuthorPayload{
			Username:       comment.Author.Username,
			DisplayName:    comment.Author.DisplayName,
			Avatar:         comment.Author.Avatar,
			Bio:            comment.Author.Bio,
			FollowerCount:  comment.Author.FollowerCount,
			FollowingCount: comment.Author.FollowingCount,
		}

		commentPayload.ID = comment.ID
		commentPayload.PostID = comment.PostID
		commentPayload.ParentCommentID = comment.ParentCommentID
		commentPayload.Content = comment.Content
		commentPayload.LikeCount = comment.LikeCount
		commentPayload.RetweetCount = comment.RetweetCount
		commentPayload.BookmarkCount = comment.BookmarkCount
		commentPayload.Impressions = comment.Impressions
		commentPayload.Image = comment.Image
		commentPayload.CreatedAt = comment.CreatedAt
		commentPayload.UpdatedAt = comment.UpdatedAt
		commentPayload.Author = authorPayload
		commentPayload.Replies = repliesPayload

		postAndCommentsPayload.Comments = append(
			postAndCommentsPayload.Comments,
			commentPayload,
		)
	}

	if post.Author.Avatar != "" {
		post.Author.Avatar = getUploadPath(post.Author.Avatar)
	}

	if post.Image != "" {
		post.Image = getUploadPath(post.Image)
	}

	authorPayload := AuthorPayload{
		Username:       post.Author.Username,
		DisplayName:    post.Author.DisplayName,
		Avatar:         post.Author.Avatar,
		Bio:            post.Author.Bio,
		FollowerCount:  post.Author.FollowerCount,
		FollowingCount: post.Author.FollowingCount,
	}

	postAndCommentsPayload.ID = post.ID
	postAndCommentsPayload.Content = post.Content
	postAndCommentsPayload.CommentCount = post.CommentCount
	postAndCommentsPayload.LikeCount = post.LikeCount
	postAndCommentsPayload.RetweetCount = post.RetweetCount
	postAndCommentsPayload.BookmarkCount = post.BookmarkCount
	postAndCommentsPayload.Impressions = post.Impressions
	postAndCommentsPayload.Image = post.Image
	postAndCommentsPayload.CreatedAt = post.CreatedAt
	postAndCommentsPayload.UpdatedAt = post.UpdatedAt
	postAndCommentsPayload.Author = authorPayload
	postAndCommentsPayload.Liked = post.Liked

	return postAndCommentsPayload
}
