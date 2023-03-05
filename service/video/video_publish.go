package video

import (
	"douyin/cache"
	"douyin/consts"
	"douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"go.uber.org/zap"
)

func Publish(userId int64, playUrl, coverUrl, title string) error {
	return NewVideoPublishFlow(userId, playUrl, coverUrl, title).Do()
}

func NewVideoPublishFlow(userId int64, playUrl, coverUrl, title string) *PublishFlow {
	return &PublishFlow{userId: userId, playUrl: playUrl, coverUrl: coverUrl, title: title}
}

type PublishFlow struct {
	userId                   int64
	playUrl, coverUrl, title string
}

func (p *PublishFlow) Do() (err error) {
	if err = p.checkNum(); err != nil {
		return
	}
	if err = p.updateData(); err != nil {
		zap.L().Error("service video_publish updateData method exec fail!", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	return nil
}

func (p *PublishFlow) checkNum() error {
	if p.title == "" {
		return e.FailVideoTitleCantNull.Err()
	}
	if len(p.title) > consts.MaxVideoTileLimit {
		return e.FailVideoTitleLimit.Err()
	}
	return nil
}

func (p *PublishFlow) updateData() (err error) {
	videoId := utils.GenID()
	if err = models.NewVideoDao().PublishVideo(videoId, p.userId, p.playUrl, p.coverUrl, p.title); err != nil {
		zap.L().Error("service video_publish PublishVideo method exec fail!", zap.Error(err))
	}
	// 保证缓存一致性，先删除缓存，再更新数据库
	go func() {
		userKey := utils.AddCacheKey(consts.CacheUser, consts.CacheSetUserVideo, utils.I64toa(p.userId))
		cache.NewUserCache().DelCache(userKey)

		var videoIds []int64
		if videoIds, err = models.NewVideoDao().QueryUserVideosById(p.userId); err != nil {
			zap.L().Error("service video_publish QueryUserVideosById method exec fail!", zap.Error(err))
		}
		if len(videoIds) > 0 {
			cache.NewUserCache().SAddReSetUserVideoList(userKey, videoIds)
		}
	}()
	return nil
}
