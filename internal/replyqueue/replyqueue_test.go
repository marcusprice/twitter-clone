package replyqueue

import (
	"sync"
	"testing"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestReplyQueueEnqueue(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	jobs := []dtypes.ReplyGuyRequest{}
	rq := &ReplyQueue{jobs: jobs}
	rq.cond = sync.NewCond(&rq.lock)
	comment := dtypes.ReplyGuyComment{Content: "yodel"}
	newJob := dtypes.ReplyGuyRequest{Comment: comment}
	rq.Enqueue(newJob)
	tu.AssertEqual("yodel", rq.jobs[0].Comment.Content)

	comment = dtypes.ReplyGuyComment{Content: "1"}
	newJob = dtypes.ReplyGuyRequest{Comment: comment}
	rq.Enqueue(newJob)

	comment = dtypes.ReplyGuyComment{Content: "2"}
	newJob = dtypes.ReplyGuyRequest{Comment: comment}
	rq.Enqueue(newJob)

	comment = dtypes.ReplyGuyComment{Content: "3"}
	newJob = dtypes.ReplyGuyRequest{Comment: comment}
	rq.Enqueue(newJob)

	comment = dtypes.ReplyGuyComment{Content: "4"}
	newJob = dtypes.ReplyGuyRequest{Comment: comment}
	rq.Enqueue(newJob)

	comment = dtypes.ReplyGuyComment{Content: "5"}
	newJob = dtypes.ReplyGuyRequest{Comment: comment}
	rq.Enqueue(newJob)

	tu.AssertEqual("1", rq.jobs[1].Comment.Content)
	tu.AssertEqual("2", rq.jobs[2].Comment.Content)
	tu.AssertEqual("3", rq.jobs[3].Comment.Content)
	tu.AssertEqual("4", rq.jobs[4].Comment.Content)
	tu.AssertEqual("5", rq.jobs[5].Comment.Content)
}
