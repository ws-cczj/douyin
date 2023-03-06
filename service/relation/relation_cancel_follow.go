package relation

import (
	"douyin/cache"
	"douyin/consts"
	"douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"go.uber.org/zap"
)

func UserCancelFollow(userId, toUserId int64) error {
	return NewUserCancelFollowFlow(userId, toUserId).Do()
}

func NewUserCancelFollowFlow(userId, toUserId int64) *UserCancelFollowFlow {
	return &UserCancelFollowFlow{userId: userId, toUserId: toUserId}
}

type UserCancelFollowFlow struct {
	isFollow, isFollow2 int

	userId, toUserId int64

	followKey string
}

func (u *UserCancelFollowFlow) Do() error {
	if err := u.checkNum(); err != nil {
		return err
	}
	if err := u.updateData(); err != nil {
		zap.L().Error("service relation_cancel_follow updateData method exec fail!", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	return nil
}

func (u *UserCancelFollowFlow) checkNum() (err error) {
	if u.userId == 0 || u.toUserId == 0 {
		return e.FailServerBusy.Err()
	}
	// 查找缓存
	u.followKey = utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(u.userId))
	relationCache := cache.NewRelationCache()
	if err = relationCache.TTLIsExpiredCache(u.followKey); err == nil {
		var isFollow bool
		if isFollow, err = relationCache.SIsMemberIsExistRelation(u.followKey, u.toUserId); err == nil {
			if isFollow {
				u.isFollow = 1
				return
			}
			return e.FailRepeatAction.Err()
		}
		zap.L().Error("service relation_cancel_follow SIsMemberIsExistRelation method exec fail!", zap.Error(err))
	}
	// 如果redis缓存获取失败或者缓存过期就去数据库查找
	if u.isFollow, err = models.NewRelationDao().IsExistRelation(u.userId, u.toUserId); err != nil {
		zap.L().Error("service relation_follow IsExistRelation method exec fail!", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	if u.isFollow != 1 {
		return e.FailRepeatAction.Err()
	}
	return
}

func (u *UserCancelFollowFlow) updateData() (err error) {
	relationDao := models.NewRelationDao()
	// 判断关系
	if u.isFollow2, err = relationDao.IsExistRelation(u.toUserId, u.userId); err != nil {
		zap.L().Error("service relation_follow IsExistRelation method exec fail!", zap.Error(err))
		return
	}
	// 此时isFollow 只能为 1, isFollow2有三种情况: 不存在关系 -1, 存在但未成立 0, 存在并且成立 1
	if err = relationDao.Action2UserRelation(u.userId, u.toUserId, u.isFollow, u.isFollow2); err != nil {
		zap.L().Error("service relation_follow Action2UserRelation method exec fail!", zap.Error(err))
		return
	}
	relationCache := cache.NewRelationCache()
	// 缓存一致性，先删除缓存，再更新数据 这里更新的是当前用户的关注
	go func() {
		relationCache.DelCache(u.followKey)
		var ids []int64
		if ids, err = relationDao.QueryUserFollowIds(u.userId); err != nil {
			zap.L().Error("service relation_follow QueryUserFollowIds method exec fail!", zap.Error(err))
		}
		relationCache.SAddResetActionUserFollowOrFollower(u.followKey, ids)
	}()
	// 更新目标用户的粉丝缓存
	go func() {
		key := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollower, utils.I64toa(u.toUserId))
		relationCache.DelCache(key)
		var ids []int64
		if ids, err = relationDao.QueryUserFollowerIds(u.toUserId); err != nil {
			zap.L().Error("service relation_follow QueryUserFollowerIds method exec fail!", zap.Error(err))
		}
		relationCache.SAddResetActionUserFollowOrFollower(key, ids)
	}()
	return nil
}
