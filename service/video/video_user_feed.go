package video

import (
	"douyin/consts"
	models "douyin/database/models"
	"douyin/pkg/e"
	"sync"
	"time"

	"go.uber.org/zap"
)

type FeedResponse struct {
	NextTime int64           `json:"next_time"`
	Videos   []*models.Video `json:"video_list"`
}

func UserFeed(lastTime, userId int64) (*FeedResponse, error) {
	return NewVideoUserFeedFlow(lastTime, userId).Do()
}

func NewVideoUserFeedFlow(lastTime, userId int64) *UserFeedFlow {
	return &UserFeedFlow{lastTime: lastTime, userId: userId, videos: make([]*models.Video, consts.MaxFeedVideos)}
}

type UserFeedFlow struct {
	lastTime int64
	userId   int64

	nextTime int64
	videos   []*models.Video

	data *FeedResponse
}

func (u *UserFeedFlow) Do() (*FeedResponse, error) {
	if err := u.checkNum(); err != nil {
		return nil, err
	}
	if err := u.prepareData(); err != nil {
		zap.L().Error("service video_user_feed prepareData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	if err := u.packData(); err != nil {
		zap.L().Error("service video_user_feed packData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	return u.data, nil
}

func (u *UserFeedFlow) checkNum() error {
	if u.userId == 0 {
		return e.FailNotKnow.Err()
	}
	if u.lastTime > time.Now().Unix() {
		zap.L().Debug("service video_user_feed uLastAt", zap.Int64("time", u.lastTime))
		return e.FailNotKnow.Err()
	}
	return nil
}

func (u *UserFeedFlow) prepareData() (err error) {
	// 1. 根据时间查询数据库中视频条数
	if err = models.NewVideoDao().QueryVideoListByTime(u.videos, u.lastTime); err != nil {
		zap.L().Error("service video_user_feed QueryVideoListByTime method exec fail!", zap.Error(err))
		return
	}
	if u.videos[0] == nil {
		if err = models.NewVideoDao().QueryVideoList(u.videos); err != nil {
			zap.L().Error("service video_user_feed QueryVideoList method exec fail!", zap.Error(err))
			return
		}
	}
	// 2. 根据每个视频id去查询用户信息
	var wg sync.WaitGroup
	for i, video := range u.videos {
		if video == nil {
			u.videos = u.videos[:i]
			break
		}
		wg.Add(1)
		go func(vdo *models.Video) {
			defer wg.Done()
			// 通过id查询用户信息
			vdo.Author = new(models.User)
			if err = models.NewUserDao().QueryUserInfoById(vdo.Author, vdo.UserId); err != nil {
				zap.L().Error("service video_user_feed QueryUserInfoById method exec fail!", zap.Error(err))
			}
			// 判断用户关系
			if u.userId != vdo.UserId {
				var isFollow int
				if isFollow, err = models.NewRelationDao().IsExistRelation(u.userId, vdo.UserId); err != nil {
					zap.L().Error("service video_user_feed NewRelationDao method exec fail!", zap.Error(err))
				}
				if isFollow == 1 {
					vdo.Author.IsFollow = true
				}
			}
			var isFavor int
			if isFavor, err = models.NewFavorDao().IsExistFavor(u.userId, vdo.VideoId); err != nil {
				zap.L().Error("service video_user_feed IsExistFavor method exec fail!", zap.Error(err))
			}
			if isFavor == 1 {
				vdo.IsFavor = true
			}
		}(video)
	}
	wg.Wait()
	return nil
}

func (u *UserFeedFlow) packData() error {
	u.data = &FeedResponse{
		NextTime: u.videos[len(u.videos)-1].CreateAt.Unix(),
		Videos:   u.videos,
	}
	return nil
}
