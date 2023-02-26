package user

import (
	"douyin/models"
	"errors"
	"sync"

	"go.uber.org/zap"
)

func FavorVideoList(userId, tkUserId int64) ([]*models.Video, error) {
	return NewFavorVideoListFlow(userId, tkUserId).Do()
}

func NewFavorVideoListFlow(userId, tkUserId int64) *FavorVideoListFlow {
	return &FavorVideoListFlow{userId: userId, tkUserId: tkUserId}
}

type FavorVideoListFlow struct {
	userId, tkUserId int64
	data             []*models.Video
}

func (f *FavorVideoListFlow) Do() ([]*models.Video, error) {
	if err := f.checkNum(); err != nil {
		return nil, err
	}
	if err := f.prepareData(); err != nil {
		return nil, err
	}
	return f.data, nil
}

func (f *FavorVideoListFlow) checkNum() error {
	if f.userId == 0 {
		return errors.New("非法用户")
	}
	return nil
}

func (f *FavorVideoListFlow) prepareData() (err error) {
	// 本用户要么是游客 tk = 0 要么是用户 tk = ?
	// 目标用户要么是本人 tk = user 要么是目标对象 tk != user
	// 需要判断当前用户与这些视频的作者是否关注以及视频是否进行了点赞

	// 1. 首先需要拿到目标用户的点赞视频数,然后拿到点赞视频列表
	favors, err := models.NewUserDao().QueryUserFavorVideos(f.userId)
	if err != nil {
		zap.L().Error("service user_favor_video_list QueryUserFavorVideos method exec fail!", zap.Error(err))
		return
	}
	if favors == 0 {
		return
	}
	f.data = make([]*models.Video, favors)
	if err = models.NewVideoDao().QueryVideoListWithFavors(f.data, f.userId); err != nil {
		zap.L().Error("service user_favor_video_list QueryVideoListWithFavors method exec fail!", zap.Error(err))
	}
	// 2. 根据视频进行遍历搜索
	for _, video := range f.data {
		// 通过id查询用户信息
		video.Author = new(models.User)
		if err = models.NewUserDao().QueryUserInfoById(video.Author, video.UserId); err != nil {
			zap.L().Error("service user_favor_video_list QueryUserInfoById method exec fail!", zap.Error(err))
		}
	}
	// 判断是不是游客
	if f.tkUserId == 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(f.data))
	for _, video := range f.data {
		vdo := video
		go func() {
			// 判断用户关系
			if f.tkUserId != vdo.UserId {
				var isFollow int
				if isFollow, err = models.NewRelationDao().IsExistRelation(f.tkUserId, vdo.UserId); err != nil {
					zap.L().Error("service user_favor_video_list NewRelationDao method exec fail!", zap.Error(err))
				}
				if isFollow == 1 {
					vdo.Author.IsFollow = true
				}
			}
			var isFavor int
			if isFavor, err = models.NewFavorDao().IsExistFavor(f.tkUserId, vdo.VideoId); err != nil {
				zap.L().Error("service user_favor_video_list IsExistFavor method exec fail!", zap.Error(err))
			}
			if isFavor == 1 {
				vdo.IsFavor = true
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return
}
