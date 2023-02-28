package user

import (
	models2 "douyin/database/models"
	"errors"
	"sync"

	"go.uber.org/zap"
)

func PublishVideoList(userId, tkUserId int64) ([]*models2.Video, error) {
	return NewPublishVideoListFlow(userId, tkUserId).Do()
}

func NewPublishVideoListFlow(userId, tkUserId int64) *PublishVideoListFlow {
	return &PublishVideoListFlow{userId: userId, tkUserId: tkUserId}
}

type PublishVideoListFlow struct {
	userId, tkUserId int64
	data             []*models2.Video
}

func (p *PublishVideoListFlow) Do() ([]*models2.Video, error) {
	if err := p.checkNum(); err != nil {
		return nil, err
	}
	if err := p.prepareData(); err != nil {
		return nil, err
	}
	return p.data, nil
}

func (p *PublishVideoListFlow) checkNum() error {
	if p.userId == 0 {
		return errors.New("非法用户")
	}
	return nil
}
func (p *PublishVideoListFlow) prepareData() (err error) {
	// 查询目标用户信息
	user := new(models2.User)
	if err = models2.NewUserDao().QueryUserInfoById(user, p.userId); err != nil {
		zap.L().Error("service user_video_list QueryUserInfoById method exec fail!", zap.Error(err))
		return
	}
	p.data = make([]*models2.Video, user.WorkCount)
	if err = models2.NewVideoDao().QueryUserVideoListById(p.data, p.userId); err != nil {
		zap.L().Error("service user_video_list QueryUserVideoListById method exec fail!", zap.Error(err))
		return
	}
	// 如果是游客访问
	if p.tkUserId == 0 {
		for _, video := range p.data {
			video.Author = user
		}
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(p.data))
	for _, video := range p.data {
		vdo := video
		go func() {
			defer wg.Done()
			vdo.Author = user
			if p.tkUserId != p.userId {
				var isFollow int
				if isFollow, err = models2.NewRelationDao().IsExistRelation(p.userId, vdo.UserId); err != nil {
					zap.L().Error("service user_video_list NewRelationDao method exec fail!", zap.Error(err))
				}
				if isFollow == 1 {
					vdo.Author.IsFollow = true
				}
			}
			var isFavor int
			if isFavor, err = models2.NewFavorDao().IsExistFavor(p.userId, vdo.VideoId); err != nil {
				zap.L().Error("service user_video_list IsExistFavor method exec fail!", zap.Error(err))
			}
			if isFavor == 1 {
				vdo.IsFavor = true
			}
		}()
	}
	wg.Wait()
	return
}
