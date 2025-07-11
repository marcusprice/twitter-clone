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

const GENERATE_ENDPOINT = "/api/generate"

type OllamaClient struct {
	host   string
	port   string
	client *http.Client
}

func (oc OllamaClient) Prompt(job dtypes.ReplyGuyRequest) (dtypes.ModelResponse, error) {
	ollamaRequestPayload := dtypes.OllamaRequest{
		Stream: false,
		Model:  job.Model,
		Prompt: formatPrompt(job),
	}

	payload, err := json.Marshal(ollamaRequestPayload)
	if err != nil {
		logger.LogError("OllamaClient.Prompt() error marshalling payload: " + err.Error())
		return dtypes.ModelResponse{}, err
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%s%s", oc.host, oc.port, GENERATE_ENDPOINT),
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

func formatPrompt(request dtypes.ReplyGuyRequest) string {
	prompt := "***************************************************\n\n"
	prompt += request.Comment.Content + "\n\n"
	prompt += fmt.Sprintf("posted by: @%s", request.Comment.Author.Username) + "\n\n"
	prompt += "***************************************************\n\n"
	prompt += "The user's prompt has ended, the following is additional context for the LLM: \n\n"
	prompt += fmt.Sprintf(
		"The top level post was posted by user @%s and the content read:\n\n%s",
		request.ParentPost.Author.Username,
		request.ParentPost.Content)

	if request.Comment.Author.Username == request.ParentComment.Author.Username {
		prompt += fmt.Sprintf(
			"\n\n(the top level post was posted by the same user @%s)",
			request.ParentPost.Author.Username,
		)
	}

	if request.ParentComment.ID != 0 {
		prompt += "\n\n"
		prompt += "This user is replying to another comment, the content of"
		prompt += "the top level comment in the thread is:"
		prompt += request.ParentComment.Content
	}

	prompt += "\n\n"
	prompt += "Feel free to reply/acknowledge the other users (include their "
	prompt += "username with the @ symbol) if it warrants it."

	return prompt
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
