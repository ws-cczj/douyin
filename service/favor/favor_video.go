package favor

import (
	models2 "douyin/database/models"
	"errors"

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
		return err
	}
	return nil
}

func (f *FavorVideoFlow) checkNum() (err error) {
	if f.userId == 0 || f.videoId == 0 {
		return errors.New("用户或视频不存在")
	}
	// 1. 检查视频是否存在
	isExist, err := models2.NewVideoDao().IsExistVideoById(f.videoId)
	if err != nil {
		zap.L().Error("service favor_video IsExistVideoById method exec fail!", zap.Error(err))
		return
	}
	if !isExist {
		zap.L().Error("service favor_video videoId not exist!", zap.Int64("videoId", f.videoId))
		return errors.New("视频不存在")
	}
	// 2. 检查数据是否合法
	if f.isFavor, err = models2.NewFavorDao().IsExistFavor(f.userId, f.videoId); err != nil {
		zap.L().Error("service favor_video IsExistFavor method exec fail!")
		return
	}
	if f.action == "1" && f.isFavor == 1 || f.action == "2" && f.isFavor < 1 {
		zap.L().Warn("service favor_video action illegal")
		return errors.New("无效操作")
	}
	return
}

func (f *FavorVideoFlow) updateData() (err error) {
	switch f.action {
	case "1":
		if err = models2.NewFavorDao().AddUserFavorVideoInfoById(f.userId, f.videoId, f.isFavor); err != nil {
			zap.L().Error("service favor_video AddUserFavorVideoInfoById method exec fail!")
			return
		}
	case "2":
		if err = models2.NewFavorDao().SubUserFavorsInfoById(f.userId, f.videoId, f.isFavor); err != nil {
			zap.L().Error("service favor_video SubUserFavorsInfoById method exec fail!")
			return
		}
	default:
		return errors.New("操作非法")
	}
	return
}
