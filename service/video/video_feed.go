package video

import "douyin/models"

type FeedResponse struct {
	NextTime int64           `json:"next_time"`
	Videos   []*models.Video `json:"videos"`
}

type FeedFlow struct {
	LastTime int64
	NextTime int64

	Token string

	Videos []*models.Video
	data   *FeedResponse
}
