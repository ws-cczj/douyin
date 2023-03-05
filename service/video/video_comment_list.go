package video

import (
	"douyin/cache"
	"douyin/consts"
	models "douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"sync"

	"go.uber.org/zap"
)

func CommentList(videoId, userId int64) ([]*models.Comment, error) {
	return NewVideoCommentList(videoId, userId).Do()
}

func NewVideoCommentList(videoId, userId int64) *CommentListFlow {
	return &CommentListFlow{videoId: videoId, userId: userId}
}

type CommentListFlow struct {
	videoId, userId int64

	data []*models.Comment
}

func (c *CommentListFlow) Do() ([]*models.Comment, error) {
	if err := c.checkNum(); err != nil {
		return nil, err
	}
	if err := c.prepareData(); err != nil {
		zap.L().Error("service video_comment_list prepareData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	if err := c.packData(); err != nil {
		zap.L().Error("service video_comment_list prepareData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	return c.data, nil
}

func (c *CommentListFlow) checkNum() (err error) {
	if c.videoId == 0 {
		return e.FailNotKnow.Err()
	}
	// 判断视频是否存在
	var exist bool
	if exist, err = models.NewVideoDao().IsExistVideoById(c.videoId); err != nil {
		zap.L().Error("service video_comment IsExistVideoById method exec fail!", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	if !exist {
		return e.FailVideoNotExist.Err()
	}
	return nil
}

func (c *CommentListFlow) prepareData() (err error) {
	// 找到该视频下的评论数量
	videoKey := utils.AddCacheKey(consts.CacheVideo, consts.CacheStringVideoFavor, utils.I64toa(c.videoId))
	var comments int64
	if comments, err = cache.NewVideoCache().GetEXVideoComments(videoKey); err != nil {
		zap.L().Error("service video_comment_list GetEXVideoComments method exec fail!", zap.Error(err))
		if comments, err = models.NewVideoDao().QueryVideoCommentsById(c.videoId); err != nil {
			zap.L().Error("service video_comment_list QueryVideoCommentsById method exec fail!", zap.Error(err))
			return
		}
	}
	c.data = make([]*models.Comment, comments)
	return
}

func (c *CommentListFlow) packData() (err error) {
	if len(c.data) < 1 {
		return nil
	}
	// 根据视频id查到该视频下所有评论信息
	if err = models.NewCommentDao().QueryVideoCommentListById(c.data, c.videoId); err != nil {
		zap.L().Error("service video_comment_list QueryVideoCommentListById method exec fail!", zap.Error(err))
	}
	if c.userId == 0 {
		return
	}
	// 判断当前用户与评论用户之间的关系
	relationDao := models.NewRelationDao()
	relationCache := cache.NewRelationCache()
	followKey := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(c.userId))
	var wg sync.WaitGroup
	wg.Add(len(c.data))
	for _, comment := range c.data {
		go func(ct *models.Comment) {
			defer wg.Done()
			if ct.IsFollow, err = relationCache.SIsMemberIsExistRelation(followKey, ct.UserId); err != nil {
				zap.L().Error("service video_comment SIsMemberIsExistRelation method exec fail!", zap.Error(err))
				// 如果缓存无效就去数据库查找
				var isFollow int
				if isFollow, err = relationDao.IsExistRelation(c.userId, ct.UserId); err != nil {
					zap.L().Error("service video_comment NewRelationDao method exec fail!", zap.Error(err))
				}
				if isFollow == 1 {
					ct.IsFollow = true
				}
			}
			ct.CreateAt = utils.FormatTime(ct.CreateTime)
		}(comment)
	}
	wg.Wait()
	return nil
}
