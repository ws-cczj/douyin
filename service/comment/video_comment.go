package comment

import (
	"douyin/consts"
	models2 "douyin/database/models"
	"douyin/pkg/utils"
	"errors"

	"go.uber.org/zap"
)

func VideoComment(userId, videoId, commentId int64, action, content string) (*models2.Comment, error) {
	return NewVideoCommentFlow(userId, videoId, commentId, action, content).Do()
}

func NewVideoCommentFlow(userId, videoId, commentId int64, action, content string) *VideoCommentFlow {
	return &VideoCommentFlow{userId: userId, videoId: videoId, commentId: commentId, action: action, content: content}
}

type VideoCommentFlow struct {
	userId, videoId, commentId int64
	action, content            string

	data *models2.Comment
}

func (v *VideoCommentFlow) Do() (*models2.Comment, error) {
	if err := v.checkNum(); err != nil {
		return nil, err
	}
	if err := v.updateData(); err != nil {
		return nil, err
	}
	if v.action == "1" {
		if err := v.packData(); err != nil {
			return nil, err
		}
	}
	return v.data, nil
}

func (v *VideoCommentFlow) checkNum() error {
	if v.userId == 0 {
		return errors.New("非法用户")
	}
	if v.videoId == 0 {
		return errors.New("视频不存在")
	}
	exist, err := models2.NewVideoDao().IsExistVideoById(v.videoId)
	if err != nil {
		zap.L().Error("service video_comment IsExistVideoById method exec fail!", zap.Error(err))
		return err
	}
	if !exist {
		return errors.New("视频不存在")
	}
	if v.action == "2" {
		if exist, err = models2.NewCommentDao().IsExistComment(v.commentId); err != nil {
			zap.L().Error("service video_comment IsExistComment method exec fail!", zap.Error(err))
			return err
		}
		if !exist {
			return errors.New("评论不存在")
		}
	}
	return nil
}

func (v *VideoCommentFlow) updateData() (err error) {
	switch v.action {
	case "1":
		if v.content == "" || len(v.content) > consts.MaxCommentLenLimit {
			return errors.New("内容字数不符合要求")
		}
		if v.commentId, err = models2.NewCommentDao().PublishVideoComment(v.userId, v.videoId, v.content); err != nil {
			zap.L().Error("service video_comment PublishVideoComment method exec fail!", zap.Error(err))
			return
		}
	case "2":
		if v.commentId == 0 {
			return errors.New("服务繁忙")
		}
		if err = models2.NewCommentDao().DeleteVideoComment(v.videoId, v.commentId); err != nil {
			zap.L().Error("service video_comment DeleteVideoComment method exec fail!", zap.Error(err))
			return
		}
	default:
		return errors.New("非法操作")
	}
	return
}

func (v *VideoCommentFlow) packData() (err error) {
	v.data = new(models2.Comment)
	if err = models2.NewCommentDao().QueryUserCommentById(v.data, v.commentId); err != nil {
		zap.L().Error("service video_comment QueryUserCommentById method exec fail!", zap.Error(err))
	}
	v.data.CreateAt = utils.FormatTime(v.data.CreateTime)
	return
}
