package video

import "douyin/models"

func VisitorFeed(lastTime int64) (*FeedResponse, error) {
	return NewVisitorFeedFlow(lastTime).Do()
}

func NewVisitorFeedFlow(lastTime int64) *VisitorFeedFlow {
	return &VisitorFeedFlow{LastTime: lastTime}
}

type VisitorFeedFlow struct {
	LastTime int64

	NextTime int64
	Videos   []*models.Video

	data *FeedResponse
}

func (u *VisitorFeedFlow) Do() (*FeedResponse, error) {
	return nil, nil
}
