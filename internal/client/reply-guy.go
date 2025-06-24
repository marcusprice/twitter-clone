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

type ReplyGuyClient struct {
	host   string
	port   string
	client *http.Client
}

func (rg *ReplyGuyClient) RequestReply(request dtypes.ReplyGuyRequest) error {
	json, err := json.Marshal(request)
	if err != nil {
		logger.LogError("ReplyGuyClient.RequestReply() error marshalling json: " + err.Error())
		if util.InDevContext() {
			panic(err)
		}

		return err
	}

	resp, err := http.Post(
		rg.address()+DALE_COOPER_ENDPOINT,
		"application/json",
		bytes.NewReader(json),
	)
	if err != nil {
		logger.LogError("ReplyGuyClient.RequestReply() post request failed: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		logger.LogWarn("ReplyGuyClient.RequestReply() non-status accepted code returned from reply-guy")
	}

	return nil
}

func (rg *ReplyGuyClient) address() string {
	return fmt.Sprintf("http://%s:%s", rg.host, rg.port)
}

func (rg *ReplyGuyClient) GetReplyGuys() []string {
	return []string{"@dalecooper"}
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
