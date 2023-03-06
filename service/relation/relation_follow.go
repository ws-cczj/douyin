package relation

import (
	"douyin/cache"
	"douyin/consts"
	"douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"go.uber.org/zap"
)

func UserFollow(userId, toUserId int64) error {
	return NewUserFollowFlow(userId, toUserId).Do()
}

func NewUserFollowFlow(userId, toUserId int64) *UserFollowFlow {
	return &UserFollowFlow{userId: userId, toUserId: toUserId}
}

type UserFollowFlow struct {
	isFollow, isFollow2 int

	userId, toUserId int64

	follow2Key string
}

func (u *UserFollowFlow) Do() error {
	if err := u.checkNum(); err != nil {
		return err
	}
	if err := u.prepareData(); err != nil {
		return e.FailServerBusy.Err()
	}
	if err := u.updateData(); err != nil {
		return e.FailServerBusy.Err()
	}
	return nil
}

func (u *UserFollowFlow) checkNum() (err error) {
	if u.userId == 0 || u.toUserId == 0 {
		return e.FailServerBusy.Err()
	}
	// 不允许自己关注自己
	if u.userId == u.toUserId {
		return e.FailCantFollowYourself.Err()
	}
	// 查看关系情况
	if u.isFollow, err = models.NewRelationDao().IsExistRelation(u.userId, u.toUserId); err != nil {
		zap.L().Error("service relation_follow IsExistRelation method exec fail!", zap.Error(err))
		err = e.FailServerBusy.Err()
	}
	if u.isFollow == 1 {
		err = e.FailRepeatAction.Err()
	}
	return
}

func (u *UserFollowFlow) prepareData() (err error) {
	// 缓存查询目标用户与当前用户之间的关系
	u.follow2Key = utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollower, utils.I64toa(u.toUserId))
	relationCache := cache.NewRelationCache()
	if err = relationCache.TTLIsExpiredCache(u.follow2Key); err == nil {
		var isFollow bool
		if isFollow, err = relationCache.SIsMemberIsExistRelation(u.follow2Key, u.userId); err == nil {
			if isFollow {
				u.isFollow2 = 1
			}
			return
		}
		zap.L().Error("service relation_follow SIsMemberIsExistRelation method exec fail!", zap.Error(err))
	}
	// 查找目标用户与当前用户之间的关系
	if u.isFollow2, err = models.NewRelationDao().IsExistRelation(u.toUserId, u.userId); err != nil {
		zap.L().Error("service relation_follow IsExistRelation method exec fail!", zap.Error(err))
	}
	// 因为粉丝情况只有两种，要么没有或者不是粉丝，要么是粉丝需要互相关注，因此这里直接简化为两种情况。
	if u.isFollow2 != 1 {
		u.isFollow2 = 0
	}
	return
}

func (u *UserFollowFlow) updateData() (err error) {
	// 此时isFollow 有两种情况, 不存在 -1, 存在但未成立 0.
	// isFollow2有两种情况, 不存在或未成立 0, 成立 1.
	relationDao := models.NewRelationDao()
	if err = relationDao.Action1UserRelation(u.userId, u.toUserId, u.isFollow, u.isFollow2); err != nil {
		zap.L().Error("service relation_follow Action1UserRelation method exec fail!", zap.Error(err))
	}
	relationCache := cache.NewRelationCache()

	// 保证缓存一致性，先删除再重置缓存 这里重置的是当前用户的关注缓存
	go func() {
		key := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(u.userId))
		relationCache.DelCache(key)
		var ids []int64
		if ids, err = relationDao.QueryUserFollowIds(u.userId); err != nil {
			zap.L().Error("service relation_follow QueryUserFollowIds method exec fail!", zap.Error(err))
		}
		relationCache.SAddResetActionUserFollowOrFollower(key, ids)
	}()

	// 更新目标用户的粉丝缓存
	go func() {
		relationCache.DelCache(u.follow2Key)
		var ids []int64
		if ids, err = relationDao.QueryUserFollowerIds(u.toUserId); err != nil {
			zap.L().Error("service relation_follow QueryUserFollowerIds method exec fail!", zap.Error(err))
		}
		relationCache.SAddResetActionUserFollowOrFollower(u.follow2Key, ids)
	}()
	return nil
}
