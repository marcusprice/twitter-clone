package testhelpers

import "github.com/marcusprice/twitter-clone/internal/dtypes"

type MockReplyGuyClient struct {
	CalledWith dtypes.ReplyGuyRequest
}

func (rg *MockReplyGuyClient) RunAsync() bool {
	return false
}

func (rg *MockReplyGuyClient) GetReplyGuys() []string {
	return []string{"@dalecooper"}
}

func (rg *MockReplyGuyClient) RequestReply(request dtypes.ReplyGuyRequest) {
	rg.CalledWith = request
}
