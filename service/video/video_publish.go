package video

import (
	"douyin/consts"
	"douyin/models"
	"douyin/pkg/utils"
	"errors"

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
	if err = p.check(); err != nil {
		return
	}
	if err = p.publish(); err != nil {
		return err
	}
	return nil
}

func (p *PublishFlow) check() error {
	if p.title == "" {
		return errors.New("标题为空!")
	}
	if len(p.title) > consts.MaxVideoTileLimit {
		return errors.New("视频标题超过要求长度!")
	}
	return nil
}

func (p *PublishFlow) publish() (err error) {
	videoId := utils.GenID()
	if err = models.NewVideoDao().PublishVideo(videoId, p.userId, p.playUrl, p.coverUrl, p.title); err != nil {
		zap.L().Error("service video_publish PublishVideo method exec fail!", zap.Error(err))
	}
	return
}
