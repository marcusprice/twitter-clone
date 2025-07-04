package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/logger"
	"github.com/marcusprice/twitter-clone/internal/util"
)

const DALE_COOPER_ENDPOINT = "/api/v1/@dalecooper/request-reply"

type ReplyGuyRequester interface {
	RunAsync() bool
	GetReplyGuys() []string
	RequestReply(request dtypes.ReplyGuyRequest)
}

// TODO: retries

type ReplyGuyClient struct {
	host   string
	port   string
	client *http.Client
}

func (rg *ReplyGuyClient) RunAsync() bool {
	return true
}

func (rg *ReplyGuyClient) GetReplyGuys() []string {
	return []string{"@dalecooper"}
}

func (rg *ReplyGuyClient) RequestReply(request dtypes.ReplyGuyRequest) {
	json, err := json.Marshal(request)
	if err != nil {
		logger.LogError("ReplyGuyClient.RequestReply() error marshalling json: " + err.Error())
		if util.InDevContext() {
			panic(err)
		}
	}

	resp, err := http.Post(
		rg.address()+DALE_COOPER_ENDPOINT,
		"application/json",
		bytes.NewReader(json),
	)
	if err != nil {
		logger.LogError("ReplyGuyClient.RequestReply() post request failed: " + err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		logger.LogWarn("ReplyGuyClient.RequestReply() non-status accepted code returned from reply-guy")
	}
}

func (rg *ReplyGuyClient) address() string {
	return fmt.Sprintf("http://%s:%s", rg.host, rg.port)
}

func NewReplyGuyClient() *ReplyGuyClient {
	host := os.Getenv("REPLY_GUY_HOST")
	port := os.Getenv("REPLY_GUY_PORT")

	client := &http.Client{}

	replyGuyClient := &ReplyGuyClient{
		host:   host,
		port:   port,
		client: client,
	}

	return replyGuyClient
}
