package replyqueue

import (
	"fmt"
	"sync"

	"github.com/marcusprice/twitter-clone/internal/api"
	"github.com/marcusprice/twitter-clone/internal/client"
	"github.com/marcusprice/twitter-clone/internal/constants"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/logger"
)

// TODO: write better tests for ReplyQueue

type ReplyQueue struct {
	jobs         []dtypes.ReplyGuyRequest
	lock         sync.Mutex
	cond         *sync.Cond
	coreClient   *client.CoreClient
	ollamaClient *client.OllamaClient
}

func (rq *ReplyQueue) Enqueue(request dtypes.ReplyGuyRequest) {
	rq.lock.Lock()
	rq.jobs = append(rq.jobs, request)
	rq.cond.Signal()
	rq.lock.Unlock()
}

func (rq *ReplyQueue) StartWorker() {
	go func() {
		for {
			rq.lock.Lock()
			for len(rq.jobs) == 0 {
				rq.cond.Wait()
			}

			job := rq.jobs[0]
			rq.jobs = rq.jobs[1:]
			rq.lock.Unlock()

			err := rq.process(job)
			if err != nil {
				// TODO: determine what to do on failed job beyond log
			}
		}
	}()
}

func (rq *ReplyQueue) process(job dtypes.ReplyGuyRequest) error {
	logger.LogInfo(
		fmt.Sprintf(
			"ReplyQueue.process() new process request for commentID: %d",
			job.Comment.ID))

	modelResponse, err := rq.ollamaClient.Prompt(job)
	if err != nil {
		return err
	}

	resp, err := rq.coreClient.PostComment(
		job.ParentPost.ID, job.ParentComment.ID, modelResponse.Response)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		logger.LogWarn("ReplyQueue.process() no error, but non-200 status code")
	}

	return nil
}

func NewReplyQueue() *ReplyQueue {
	// TODO: need to modularize JWT + handle expiration date
	dalecooperJWT, err := api.GenerateJWT(constants.DALE_COOPER_USER_ID)
	if err != nil {
		panic(err)
	}

	coreClient := client.NewCoreClient(dalecooperJWT)
	ollamaClient := client.NewOllamaClient()
	jobs := []dtypes.ReplyGuyRequest{}

	replyQueue := &ReplyQueue{
		jobs:         jobs,
		coreClient:   coreClient,
		ollamaClient: ollamaClient,
	}
	replyQueue.cond = sync.NewCond(&replyQueue.lock)

	return replyQueue
}
