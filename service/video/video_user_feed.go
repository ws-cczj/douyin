package video

import (
	"douyin/consts"
	"douyin/models"
	"errors"
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
		return nil, err
	}
	if err := u.packData(); err != nil {
		return nil, err
	}
	return u.data, nil
}

func (u *UserFeedFlow) checkNum() error {
	if u.userId == 0 {
		return errors.New("非法用户")
	}
	if u.lastTime > time.Now().Unix() {
		zap.L().Debug("service video_user_feed uLastAt", zap.Int64("time", time.Now().Unix()))
		return errors.New("未知错误")
	}
	return nil
}

func (u *UserFeedFlow) prepareData() (err error) {
	// 1. 根据时间查询数据库中视频条数
	if err = models.NewVideoDao().QueryVideoListByTime(u.videos, u.lastTime); err != nil {
		zap.L().Error("service video_user_feed QueryVideoListByTime method exec fail!", zap.Error(err))
		return err
	}
	if u.videos[0] == nil {
		if err = models.NewVideoDao().QueryVideoList(u.videos); err != nil {
			zap.L().Error("service video_user_feed QueryVideoList method exec fail!", zap.Error(err))
			return err
		}
	}
	// 2. 根据每个视频id去查询用户信息
	for i, video := range u.videos {
		if video == nil {
			u.videos = u.videos[:i]
			break
		}
		// 通过id查询用户信息
		video.Author = new(models.User)
		if err = models.NewUserDao().QueryUserInfoById(video.Author, video.UserId); err != nil {
			zap.L().Error("service video_user_feed QueryUserInfoById method exec fail!", zap.Error(err))
			continue
		}
		// 判断用户关系
		if u.userId != video.UserId {
			var isFollow int
			if isFollow, err = models.NewRelationDao().IsExistRelation(u.userId, video.UserId); err != nil {
				zap.L().Error("service video_user_feed NewRelationDao method exec fail!", zap.Error(err))
				continue
			}
			if isFollow == 1 {
				video.Author.IsFollow = true
			}
		}
		var isFavor int
		if isFavor, err = models.NewFavorDao().IsExistFavor(u.userId, video.VideoId); err != nil {
			zap.L().Error("service video_user_feed IsExistFavor method exec fail!", zap.Error(err))
			continue
		}
		if isFavor == 1 {
			video.IsFavor = true
		}
		u.nextTime = video.CreateAt.Unix()
	}
	return
}

func (u *UserFeedFlow) packData() error {
	u.data = &FeedResponse{
		NextTime: u.nextTime,
		Videos:   u.videos,
	}
	return nil
}
