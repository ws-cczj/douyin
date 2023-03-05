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

	followKey string
	action    bool

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
	if err := f.packData(); err != nil {
		zap.L().Error("service relation_user_list packData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	return f.data, nil
}

func (u *UserRelationListFlow) checkNum() (err error) {
	if u.userId == 0 || u.tkUserId == 0 {
		return e.FailNotKnow.Err()
	}
	// 预热缓存
	relationCache := cache.NewRelationCache()
	u.followKey = utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(u.tkUserId))
	if err = relationCache.TTLIsExpiredCache(u.followKey); err != nil {
		zap.L().Warn("service user_publish_video_list relationCache.TTLIsExpiredCache method exec fail!", zap.Error(err))
		var ids []int64
		if ids, err = models.NewRelationDao().QueryUserFollowIds(u.tkUserId); err != nil {
			zap.L().Error("service user_publish_video_list QueryUserFollowIds method exec fail!", zap.Error(err))
		}
		relationCache.SAddResetActionUserFollowOrFollower(u.followKey, ids)
	}
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
		u.data = make([]*models.User, follows)
		// 获取关注列表
		if err = models.NewRelationDao().QueryUserFollowList(u.data, u.userId); err != nil {
			zap.L().Error("service relation_user_list QueryUserFollowList method exec fail!", zap.Error(err))
			return
		}
	} else {
		// TODO 这里没有对粉丝数进行redis存储，后续再进行
		// 获取粉丝数
		var followers int64
		if followers, err = models.NewUserDao().QueryUserFollowers(u.userId); err != nil {
			zap.L().Error("service relation_user_list QueryUserFollowers method exec fail!", zap.Error(err))
		}
		u.data = make([]*models.User, followers)
		// 获取粉丝列表
		if err = models.NewRelationDao().QueryUserFollowerList(u.data, u.userId); err != nil {
			zap.L().Error("service relation_user_list QueryUserFollowerList method exec fail!", zap.Error(err))
			return
		}
	}
	return
}

func (u *UserRelationListFlow) packData() (err error) {
	relationCache := cache.NewRelationCache()
	relationDao := models.NewRelationDao()
	var wg sync.WaitGroup
	wg.Add(len(u.data))
	for _, data := range u.data {
		go func(user *models.User) {
			defer wg.Done()
			if user.IsFollow, err = relationCache.SIsMemberIsExistRelation(u.followKey, user.UserId); err != nil {
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
	return
}
