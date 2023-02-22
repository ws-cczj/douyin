package video

import "douyin/models"

type FeedResponse struct {
	NextTime int64           `json:"next_time"`
	Videos   []*models.Video `json:"videos"`
}

func UserFeed(lastTime, userId int64) (*FeedResponse, error) {
	return NewUserFeedFlow(lastTime, userId).Do()
}

func NewUserFeedFlow(lastTime, userId int64) *UserFeedFlow {
	return &UserFeedFlow{LastTime: lastTime, userId: userId}
}

type UserFeedFlow struct {
	LastTime int64
	userId   int64

	nextTime int64
	videos   []*models.Video

	data *FeedResponse
}

func (u *UserFeedFlow) Do() (*FeedResponse, error) {
	return nil, nil
}
