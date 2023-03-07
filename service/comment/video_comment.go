package comment

import (
	"douyin/cache"
	"douyin/consts"
	models "douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"go.uber.org/zap"
)

func VideoComment(userId, videoId, commentId int64, action, content string) (*models.Comment, error) {
	return NewVideoCommentFlow(userId, videoId, commentId, action, content).Do()
}

func NewVideoCommentFlow(userId, videoId, commentId int64, action, content string) *VideoCommentFlow {
	return &VideoCommentFlow{userId: userId, videoId: videoId, commentId: commentId, action: action, content: content}
}

type VideoCommentFlow struct {
	userId, videoId, commentId int64
	action, content            string

	data *models.Comment
}

func (v *VideoCommentFlow) Do() (*models.Comment, error) {
	if err := v.checkNum(); err != nil {
		return nil, err
	}
	if err := v.updateData(); err != nil {
		zap.L().Error("service video_comment updateData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	// 如果是评论操作需要进行数据打包
	if v.action == "1" {
		if err := v.packData(); err != nil {
			zap.L().Error("service video_comment packData method exec fail!", zap.Error(err))
			return nil, e.FailServerBusy.Err()
		}
	}
	return v.data, nil
}

func (v *VideoCommentFlow) checkNum() (err error) {
	// 判断参数是否正确
	if v.userId == 0 || v.videoId == 0 {
		return e.FailNotKnow.Err()
	}
	// 判断操作是否有误
	if v.action != "1" && v.action != "2" {
		return e.FailNotKnow.Err()
	}
	// 判断内容是否合格
	if v.action == "1" {
		if v.content == "" || len(v.content) > consts.MaxCommentLenLimit {
			return e.FailCommentLenLimit.Err()
		}
		// 过滤关键词
		v.content = utils.Replace(v.content)
	}
	// 判断视频是否存在
	var exist bool
	if exist, err = models.NewVideoDao().IsExistVideoById(v.videoId); err != nil {
		zap.L().Error("service video_comment IsExistVideoById method exec fail!", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	if !exist {
		return e.FailVideoNotExist.Err()
	}
	// 判断评论是否存在
	if v.action == "2" {
		if v.commentId == 0 {
			return e.FailNotKnow.Err()
		}
		if exist, err = models.NewCommentDao().IsExistComment(v.commentId); err != nil {
			zap.L().Error("service video_comment IsExistComment method exec fail!", zap.Error(err))
			return e.FailServerBusy.Err()
		}
		if !exist {
			return e.FailCommentNotExist.Err()
		}
	}
	return
}

func (v *VideoCommentFlow) updateData() (err error) {
	if v.action == "1" {
		if v.commentId, err = models.NewCommentDao().PublishVideoComment(v.userId, v.videoId, v.content); err != nil {
			zap.L().Error("service video_comment PublishVideoComment method exec fail!", zap.Error(err))
			return
		}
	} else {
		if err = models.NewCommentDao().DeleteVideoComment(v.videoId, v.commentId); err != nil {
			zap.L().Error("service video_comment DeleteVideoComment method exec fail!", zap.Error(err))
			return
		}
	}
	// 保证缓存一致性 先删除后更新数据
	go func() {
		videoKey := utils.AddCacheKey(consts.CacheVideo, consts.CacheStringVideoComment, utils.I64toa(v.videoId))
		cache.NewVideoCache().DelCache(videoKey)

		var commentIds int64
		if commentIds, err = models.NewCommentDao().QueryVideoCommentsById(v.videoId); err != nil {
			zap.L().Error("service video_comment QueryVideoCommentsById method exec fail!", zap.Error(err))
		}
		cache.NewVideoCache().SetEXResetVideoComments(videoKey, commentIds)
	}()
	return nil
}

func (v *VideoCommentFlow) packData() (err error) {
	v.data = new(models.Comment)
	if err = models.NewCommentDao().QueryUserCommentById(v.data, v.commentId); err != nil {
		zap.L().Error("service video_comment QueryUserCommentById method exec fail!", zap.Error(err))
	}
	v.data.CreateAt = utils.FormatTime(v.data.CreateTime)
	return
}
