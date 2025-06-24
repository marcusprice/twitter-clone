package client

import (
	"fmt"
	"net/http"
	"os"

	"github.com/marcusprice/twitter-clone/internal/logger"
	"github.com/marcusprice/twitter-clone/internal/util"
)

const COMMENT_API_ENDPOINT = "/api/v1/comment/create"

type CoreClient struct {
	host      string
	port      string
	client    *http.Client
	authToken string
}

func (cc *CoreClient) PostComment(postID, parentCommentID int, content string) (*http.Response, error) {
	fields := make(map[string]string)
	fields["content"] = content
	fields["postID"] = fmt.Sprintf("%d", postID)
	fields["parentCommentID"] = fmt.Sprintf("%d", parentCommentID)

	requestBody, contentType, err := util.GenerateMultipartForm(fields)
	if err != nil {
		logger.LogError("CoreClient.PostComment() error generating multipart form: " + err.Error())
		return &http.Response{}, err
	}
	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s:%s%s", cc.host, cc.port, COMMENT_API_ENDPOINT),
		requestBody)

	if err != nil {
		logger.LogError("CoreClient.PostComment() error creating new request: " + err.Error())
		return &http.Response{}, err
	}

	request.Header.Set("Authorization", "Bearer "+cc.authToken)
	request.Header.Set("Content-Type", contentType)

	apiResponse, err := cc.client.Do(request)
	if err != nil {
		return &http.Response{}, err
	}

	return apiResponse, nil
}

func NewCoreClient(jwtToken string) *CoreClient {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	client := &http.Client{}
	cc := &CoreClient{
		host:      host,
		port:      port,
		client:    client,
		authToken: jwtToken,
	}

	return cc
}
