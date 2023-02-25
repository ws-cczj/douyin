package video

import (
	"douyin/consts"
	"douyin/models"
	"errors"
	"time"

	"go.uber.org/zap"
)

func VisitorFeed(lastTime int64) (*FeedResponse, error) {
	return NewVideoVisitorFeedFlow(lastTime).Do()
}

func NewVideoVisitorFeedFlow(lastTime int64) *VisitorFeedFlow {
	return &VisitorFeedFlow{lastTime: lastTime, videos: make([]*models.Video, consts.MaxFeedVideos)}
}

type VisitorFeedFlow struct {
	lastTime int64

	nextTime int64
	videos   []*models.Video

	data *FeedResponse
}

func (u *VisitorFeedFlow) Do() (*FeedResponse, error) {
	if err := u.checkNum(); err != nil {
		return nil, err
	}
	if err := u.prepareData(); err != nil {
		return nil, err
	}
	if err := u.packData(); err != nil {
		return nil, err
	}
	return u.data, nil
}

func (u *VisitorFeedFlow) checkNum() error {
	if u.lastTime > time.Now().Unix() {
		return errors.New("未知错误")
	}
	return nil
}

func (u *VisitorFeedFlow) prepareData() (err error) {
	// 根据时间查询数据库中视频条数
	if err = models.NewVideoDao().QueryVideoListByTime(u.videos, u.lastTime); err != nil {
		zap.L().Error("service video_visitor_feed QueryVideoListByTime method exec fail!", zap.Error(err))
		return err
	}
	if u.videos[0] == nil {
		if err = models.NewVideoDao().QueryVideoList(u.videos); err != nil {
			zap.L().Error("service video_visitor QueryVideoList method exec fail!", zap.Error(err))
			return err
		}
	}
	for i, video := range u.videos {
		if video == nil {
			u.videos = u.videos[:i]
			break
		}
		video.Author = new(models.User)
		if err = models.NewUserDao().QueryUserInfoById(video.Author, video.UserId); err != nil {
			zap.L().Error("service video_visitor_feed QueryUserInfoById method exec fail!", zap.Error(err))
			continue
		}
		u.nextTime = video.CreateAt.Unix()
	}
	return
}

func (u *VisitorFeedFlow) packData() error {
	u.data = &FeedResponse{
		NextTime: u.nextTime,
		Videos:   u.videos,
	}
	return nil
}
