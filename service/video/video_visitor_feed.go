package video

import (
	"douyin/consts"
	models "douyin/database/models"
	"douyin/pkg/e"
	"sync"
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
		zap.L().Error("service video_visitor_feed prepareData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	if err := u.packData(); err != nil {
		zap.L().Error("service video_visitor_feed packData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	return u.data, nil
}

func (u *VisitorFeedFlow) checkNum() error {
	if u.lastTime > time.Now().Unix() {
		return e.FailNotKnow.Err()
	}
	return nil
}

func (u *VisitorFeedFlow) prepareData() (err error) {
	videoDao := models.NewVideoDao()
	// 根据时间查询数据库中视频条数
	formatT := time.UnixMilli(u.lastTime).Format("2006-01-02 15:04:05")
	if err = videoDao.QueryVideoListByTime(u.videos, formatT); err != nil {
		zap.L().Error("service video_visitor_feed QueryVideoListByTime method exec fail!", zap.Error(err))
		return
	}
	if u.videos[0] == nil {
		if err = videoDao.QueryVideoList(u.videos); err != nil {
			zap.L().Error("service video_visitor QueryVideoList method exec fail!", zap.Error(err))
			return
		}
	}
	// 填充数据
	userDao := models.NewUserDao()
	var wg sync.WaitGroup
	for i, video := range u.videos {
		if video == nil {
			u.videos = u.videos[:i]
			break
		}
		wg.Add(1)
		go func(vdo *models.Video) {
			defer wg.Done()
			vdo.Author = new(models.User)
			if err = userDao.QueryUserInfoById(vdo.Author, vdo.UserId); err != nil {
				zap.L().Error("service video_visitor_feed QueryUserInfoById method exec fail!", zap.Error(err))
			}
		}(video)
	}
	wg.Wait()
	return nil
}

func (u *VisitorFeedFlow) packData() error {
	u.data = &FeedResponse{
		NextTime: u.videos[0].CreateAt.UnixMilli(),
		Videos:   u.videos,
	}
	return nil
}
