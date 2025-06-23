package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/logger"
)

type OllamaClient struct {
	host   string
	port   string
	client *http.Client
}

func (oc OllamaClient) Prompt(job dtypes.ReplyGuyRequest) (dtypes.ModelResponse, error) {
	payload, err := json.Marshal(job)
	if err != nil {
		logger.LogError("OllamaClient.Prompt() error marshalling payload: " + err.Error())
		return dtypes.ModelResponse{}, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s:%s", oc.host, oc.port),
		"application/json",
		bytes.NewReader(payload))

	if err != nil {
		return dtypes.ModelResponse{}, err
	}

	defer resp.Body.Close()

	var modelResponse dtypes.ModelResponse
	json.NewDecoder(resp.Body).Decode(&modelResponse)

	return modelResponse, nil
}

func NewOllamaClient() *OllamaClient {
	ollamaHost := os.Getenv("OLLAMA_HOST")
	ollamaPort := os.Getenv("OLLAMA_PORT")
	client := &http.Client{}

	oc := &OllamaClient{
		host:   ollamaHost,
		port:   ollamaPort,
		client: client,
	}

	return oc
}
