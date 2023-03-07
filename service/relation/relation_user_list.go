package relation

import (
	"douyin/cache"
	"douyin/consts"
	models "douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"sync"

	"go.uber.org/zap"
)

func UserFollowList(userId, tkUserId int64, action bool) ([]*models.User, error) {
	return NewUserRelationListFlow(userId, tkUserId, action).Do()
}

func NewUserRelationListFlow(userId, tkUserId int64, action bool) *UserRelationListFlow {
	return &UserRelationListFlow{userId: userId, tkUserId: tkUserId, action: action}
}

type UserRelationListFlow struct {
	userId, tkUserId int64

	followKey, followerKey string
	action                 bool

	data []*models.User
}

func (f *UserRelationListFlow) Do() ([]*models.User, error) {
	if err := f.checkNum(); err != nil {
		return nil, err
	}
	if err := f.prepareData(); err != nil {
		zap.L().Error("service relation_user_list prepareData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	// 如果不是游客访问就去查找关系
	if f.tkUserId != 0 && len(f.data) > 0 {
		if err := f.packData(); err != nil {
			zap.L().Error("service relation_user_list packData method exec fail!", zap.Error(err))
			return nil, e.FailServerBusy.Err()
		}
	}
	return f.data, nil
}

func (u *UserRelationListFlow) checkNum() (err error) {
	if u.userId == 0 {
		return e.FailNotKnow.Err()
	}
	var wg sync.WaitGroup
	wg.Add(1)
	relationCache := cache.NewRelationCache()
	go func() {
		defer wg.Done()
		// 预热用户关注缓存 这里预热的是目标用户
		u.followKey = utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(u.userId))
		if err = relationCache.TTLIsExpiredCache(u.followKey); err != nil {
			zap.L().Warn("service user_publish_video_list relationCache.TTLIsExpiredCache method exec fail!", zap.Error(err))
			var ids []int64
			if ids, err = models.NewRelationDao().QueryUserFollowIds(u.userId); err != nil {
				zap.L().Error("service user_publish_video_list QueryUserFollowIds method exec fail!", zap.Error(err))
			}
			relationCache.SAddResetActionUserFollowOrFollower(u.followKey, ids)
		}
	}()
	// 预热用户粉丝缓存
	u.followerKey = utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollower, utils.I64toa(u.userId))
	if err = relationCache.TTLIsExpiredCache(u.followerKey); err != nil {
		zap.L().Warn("service user_publish_video_list relationCache.TTLIsExpiredCache method exec fail!", zap.Error(err))
		var ids []int64
		if ids, err = models.NewRelationDao().QueryUserFollowerIds(u.userId); err != nil {
			zap.L().Error("service user_publish_video_list QueryUserFollowerIds method exec fail!", zap.Error(err))
		}
		relationCache.SAddResetActionUserFollowOrFollower(u.followerKey, ids)
	}
	wg.Wait()
	return nil
}

func (u *UserRelationListFlow) prepareData() (err error) {
	if u.action {
		var follows int64
		if follows, err = cache.NewRelationCache().SCardQueryUserFollows(u.followKey); follows < 0 {
			zap.L().Warn("service relation_user_list SCardQueryUserFollows method exec fail!", zap.Error(err))
			// 如果缓存无效就去数据库查找
			if follows, err = models.NewUserDao().QueryUserFollows(u.userId); err != nil {
				zap.L().Error("service relation_user_list QueryUserFollows method exec fail!", zap.Error(err))
			}
		}
		// 这里先查询关注数量的原因有两个，1：如果没有关注直接可以返回，2：提前申请容量，避免触发append扩容机制节省资源
		if follows > 0 || err != nil {
			u.data = make([]*models.User, follows)
			// 获取关注列表
			if err = models.NewRelationDao().QueryUserFollowList(u.data, u.userId); err != nil {
				zap.L().Error("service relation_user_list QueryUserFollowList method exec fail!", zap.Error(err))
			}
		}
	} else {
		// 获取粉丝数
		var followers int64
		if followers, err = cache.NewRelationCache().SCardQueryUserFollowers(u.followerKey); followers < 0 {
			zap.L().Warn("service relation_user_list SCardQueryUserFollowers method exec fail!", zap.Error(err))
			// 如果缓存无效就去数据库查找
			if followers, err = models.NewUserDao().QueryUserFollowers(u.userId); err != nil {
				zap.L().Error("service relation_user_list QueryUserFollowers method exec fail!", zap.Error(err))
			}
		}
		// 如果错误不为空说明一定是数据库查找粉丝数错误，这里错误是可以容忍的，直接查询数据库数据即可!
		if followers > 0 || err != nil {
			u.data = make([]*models.User, followers)
			// 获取粉丝列表
			if err = models.NewRelationDao().QueryUserFollowerList(u.data, u.userId); err != nil {
				zap.L().Error("service relation_user_list QueryUserFollowerList method exec fail!", zap.Error(err))
			}
		}
	}
	return
}

func (u *UserRelationListFlow) packData() (err error) {
	tkUserKey := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(u.tkUserId))
	relationCache := cache.NewRelationCache()
	relationDao := models.NewRelationDao()
	var wg sync.WaitGroup
	wg.Add(len(u.data))
	for i, data := range u.data {
		if data == nil {
			u.data = u.data[:i]
			wg.Add(i - len(u.data))
			zap.L().Warn("service relation_user_list user is null data!", zap.Any("data", u.data))
			return
		}
		go func(user *models.User) {
			defer wg.Done()
			// 这里查找的关系是当前访问的用户
			if user.IsFollow, err = relationCache.SIsMemberIsExistRelation(tkUserKey, user.UserId); err != nil {
				zap.L().Error("service video_comment SIsMemberIsExistRelation method exec fail!", zap.Error(err))
				// 如果缓存无效就去数据库查找
				var isFollow int
				if isFollow, err = relationDao.IsExistRelation(u.tkUserId, user.UserId); err != nil {
					zap.L().Error("service relation_user_list IsExistRelation method exec fail!", zap.Error(err))
				}
				if isFollow == 1 {
					user.IsFollow = true
				}
			}
		}(data)
	}
	wg.Wait()
	return
}
