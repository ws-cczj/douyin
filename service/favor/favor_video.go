package favor

import (
	"douyin/cache"
	"douyin/consts"
	models "douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"go.uber.org/zap"
)

func FavorVideo(userId, videoId int64, action string) error {
	return NewFavorVideoFlow(userId, videoId, action).Do()
}

func NewFavorVideoFlow(userId, videoId int64, action string) *FavorVideoFlow {
	return &FavorVideoFlow{userId: userId, videoId: videoId, action: action}
}

type FavorVideoFlow struct {
	action          string
	userId, videoId int64

	isFavor int // -1 不存在数据, 0 存在但未点赞, 1存在并且点赞了
}

func (f *FavorVideoFlow) Do() error {
	if err := f.checkNum(); err != nil {
		return err
	}
	if err := f.updateData(); err != nil {
		zap.L().Error("service favor_video updateData method exec fail!", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	return nil
}

func (f *FavorVideoFlow) checkNum() (err error) {
	if f.userId == 0 || f.videoId == 0 {
		return e.FailServerBusy.Err()
	}
	if f.action != "1" && f.action != "2" {
		return e.FailNotKnow.Err()
	}
	// 1. 检查视频是否存在
	var isExist bool
	if isExist, err = models.NewVideoDao().IsExistVideoById(f.videoId); err != nil {
		zap.L().Error("service favor_video IsExistVideoById method exec fail!", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	if !isExist {
		zap.L().Error("service favor_video videoId not exist!", zap.Int64("videoId", f.videoId))
		return e.FailVideoNotExist.Err()
	}
	// 2. 检查数据是否合法
	if f.isFavor, err = models.NewFavorDao().IsExistFavor(f.userId, f.videoId); err != nil {
		zap.L().Error("service favor_video IsExistFavor method exec fail!")
		return e.FailServerBusy.Err()
	}
	if f.action == "1" && f.isFavor == 1 || f.action == "2" && f.isFavor < 1 {
		zap.L().Warn("service favor_video action illegal")
		return e.FailRepeatAction.Err()
	}
	return
}

func (f *FavorVideoFlow) updateData() (err error) {
	if f.action == "1" {
		if err = models.NewFavorDao().AddUserFavorVideoInfoById(f.userId, f.videoId, f.isFavor); err != nil {
			zap.L().Error("service favor_video AddUserFavorVideoInfoById method exec fail!")
			return
		}
	} else {
		if err = models.NewFavorDao().SubUserFavorsInfoById(f.userId, f.videoId, f.isFavor); err != nil {
			zap.L().Error("service favor_video SubUserFavorsInfoById method exec fail!")
			return
		}
	}
	// 保证缓存一致性，先删除后更新缓存
	go func() {
		favorKey := utils.StrI64(consts.CacheSetUserFavor, f.userId)
		cache.NewFavorCache().DelCache(favorKey)

		var videos []int64
		if videos, err = models.NewFavorDao().QueryUserFavorVideoList(f.userId); err != nil {
			zap.L().Error("service favor_video QueryUserFavorVideoList method exec fail!", zap.Error(err))
		}
		cache.NewFavorCache().SAddReSetUserFavorVideo(favorKey, videos)
	}()
	return nil
}
