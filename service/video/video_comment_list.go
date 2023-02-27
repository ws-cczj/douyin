package video

import (
	"douyin/models"
	"douyin/pkg/utils"
	"errors"
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
		return nil, err
	}
	if err := c.packData(); err != nil {
		return nil, err
	}
	return c.data, nil
}

func (c *CommentListFlow) checkNum() error {
	if c.videoId == 0 {
		return errors.New("服务繁忙")
	}
	return nil
}

func (c *CommentListFlow) prepareData() (err error) {
	// 找到该视频下的评论数量
	var comments int64
	if comments, err = models.NewVideoDao().QueryVideoCommentsById(c.videoId); err != nil {
		zap.L().Error("service video_comment_list QueryVideoCommentsById method exec fail!", zap.Error(err))
		return
	}
	c.data = make([]*models.Comment, comments)
	return
}

func (c *CommentListFlow) packData() (err error) {
	// 根据视频id查到该视频下所有评论信息
	if err = models.NewCommentDao().QueryVideoCommentListById(c.data, c.videoId); err != nil {
		zap.L().Error("service video_comment_list QueryVideoCommentListById method exec fail!", zap.Error(err))
	}
	if c.userId == 0 {
		return
	}
	// 判断当前用户与评论用户之间的关系
	var wg sync.WaitGroup
	wg.Add(len(c.data))
	for _, comment := range c.data {
		ct := comment
		go func() {
			defer wg.Done()
			var isFollow int
			if isFollow, err = models.NewRelationDao().IsExistRelation(c.userId, ct.User.UserId); err != nil {
				zap.L().Error("service video_comment_list IsExistRelation method exec fail!", zap.Error(err))
			}
			if isFollow == 1 {
				ct.User.IsFollow = true
			}
			ct.CreateAt = utils.FormatTime(ct.CreateTime)
		}()
	}
	wg.Wait()
	return
}
