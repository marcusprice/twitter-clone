package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/marcusprice/twitter-clone/internal/api"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/logger"
	"github.com/marcusprice/twitter-clone/internal/replyqueue"
	"github.com/marcusprice/twitter-clone/internal/util"
)

func ReplyGuyHandler(replyQueue *replyqueue.ReplyQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody dtypes.ReplyGuyRequest
		json.NewDecoder(r.Body).Decode(&requestBody)
		replyQueue.Enqueue(requestBody)
		w.WriteHeader(http.StatusAccepted)
	}
}

func main() {
	util.LoadEnvVariables()

	replyQueue := replyqueue.NewReplyQueue()
	replyQueue.StartWorker()

	mux := http.NewServeMux()
	mux.Handle(
		"/api/v1/@dalecooper/request-reply",
		api.VerifyPostMethod(
			ReplyGuyHandler(replyQueue),
		),
	)

	host := os.Getenv("REPLY_GUY_HOST")
	port := os.Getenv("REPLY_GUY_PORT")
	logger.LogInfo(fmt.Sprintf("REPLY GUY LISTENING AT %s:%s", host, port))
	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), mux)
}
