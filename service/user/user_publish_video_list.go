package user

import (
	"douyin/cache"
	"douyin/consts"
	"douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"sync"

	"go.uber.org/zap"
)

func PublishVideoList(userId, tkUserId int64) ([]*models.Video, error) {
	return NewPublishVideoListFlow(userId, tkUserId).Do()
}

func NewPublishVideoListFlow(userId, tkUserId int64) *PublishVideoListFlow {
	return &PublishVideoListFlow{userId: userId, tkUserId: tkUserId}
}

type PublishVideoListFlow struct {
	userId, tkUserId int64

	followKey, favorKey string

	data []*models.Video
}

func (p *PublishVideoListFlow) Do() ([]*models.Video, error) {
	if err := p.checkNum(); err != nil {
		return nil, err
	}
	if err := p.prepareData(); err != nil {
		return nil, e.FailServerBusy.Err()
	}
	// 如果是游客访问不用去判断关注和点赞
	if p.tkUserId != 0 {
		if err := p.packData(); err != nil {
			return nil, e.FailServerBusy.Err()
		}
	}
	return p.data, nil
}

func (p *PublishVideoListFlow) checkNum() error {
	if p.userId == 0 {
		return e.FailServerBusy.Err()
	}
	return nil
}

func (p *PublishVideoListFlow) prepareData() (err error) {
	var wg sync.WaitGroup
	if p.tkUserId != 0 {
		// 只有不为0的情况下才add，否则不进行add
		wg.Add(1)
		// 查询用户关注缓存是否过期, 如果过期则对缓存进行重置操作
		go func() {
			defer wg.Done()
			relationCache := cache.NewRelationCache()
			p.followKey = utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(p.tkUserId))
			if err = relationCache.TTLIsExpiredCache(p.followKey); err != nil {
				zap.L().Warn("service user_publish_video_list relationCache.TTLIsExpiredCache method exec fail!", zap.Error(err))
				var ids []int64
				if ids, err = models.NewRelationDao().QueryUserFollowIds(p.tkUserId); err != nil {
					zap.L().Error("service user_publish_video_list QueryUserFollowIds method exec fail!", zap.Error(err))
				}
				relationCache.SAddResetActionUserFollowOrFollower(p.followKey, ids)
			}
		}()
	}
	user := new(models.User)
	// 查询目标用户信息
	if err = models.NewUserDao().QueryUserInfoById(user, p.userId); err != nil {
		zap.L().Error("service user_video_list QueryUserInfoById method exec fail!", zap.Error(err))
		return
	}
	wg.Wait()
	// 更新用户点赞缓存缓存
	wg.Add(1)
	go func() {
		defer wg.Done()
		favorCache := cache.NewFavorCache()
		p.favorKey = utils.AddCacheKey(consts.CacheFavor, consts.CacheSetUserVideo, utils.I64toa(p.tkUserId))
		if err = favorCache.TTLIsExpiredCache(p.favorKey); err != nil {
			zap.L().Warn("service user_publish_video_list favorCache.TTLIsExpiredCache method exec fail!", zap.Error(err))
			var ids []int64
			if ids, err = models.NewFavorDao().QueryUserFavorVideoList(p.tkUserId); err != nil {
				zap.L().Error("service user_publish_video_list QueryUserFavorVideoList method exec fail!", zap.Error(err))
			}
			favorCache.SAddReSetUserFavorVideo(p.favorKey, ids)
		}
	}()
	// 判断当前查询用户是否关注了目标用户
	if p.tkUserId != user.UserId {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 通过缓存查找关系
			if user.IsFollow, err = cache.NewRelationCache().SIsMemberIsExistRelation(p.followKey, p.userId); err != nil {
				zap.L().Error("service user_publish_video_list SIsMemberIsExistRelation method exec fail!", zap.Error(err))
				// 如果缓存无效就去数据库查找
				var isFollow int
				if isFollow, err = models.NewRelationDao().IsExistRelation(p.userId, user.UserId); err != nil {
					zap.L().Error("service user_video_list NewRelationDao method exec fail!", zap.Error(err))
				}
				if isFollow == 1 {
					user.IsFollow = true
				}
			}
		}()
	}
	p.data = make([]*models.Video, user.WorkCount)
	// 根据用户id查询用户发布的视频列表信息
	if err = models.NewVideoDao().QueryUserVideoListById(p.data, p.userId); err != nil {
		zap.L().Error("service user_video_list QueryUserVideoListById method exec fail!", zap.Error(err))
		return
	}
	wg.Wait()
	// 填充数据
	for _, video := range p.data {
		video.Author = user
	}
	return
}

func (p *PublishVideoListFlow) packData() (err error) {
	// 使用协程简化循环TTL时间
	var wg sync.WaitGroup
	wg.Add(len(p.data))
	for _, video := range p.data {
		go func(vdo *models.Video) {
			defer wg.Done()
			// 通过缓存查找点赞
			if vdo.IsFavor, err = cache.NewFavorCache().SIsMemberIsExistFavor(p.favorKey, vdo.VideoId); err != nil {
				zap.L().Error("service user_publish_video_list SIsMemberIsExistFavor method exec fail!", zap.Error(err))
				// 如果缓存无效就去数据库中找
				var isFavor int
				if isFavor, err = models.NewFavorDao().IsExistFavor(p.userId, vdo.VideoId); err != nil {
					zap.L().Error("service user_publish_video_list IsExistFavor method exec fail!", zap.Error(err))
				}
				if isFavor == 1 {
					vdo.IsFavor = true
				}
			}
		}(video)
	}
	wg.Wait()
	return
}
