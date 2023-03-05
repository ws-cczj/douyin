package relation

import (
	"douyin/cache"
	"douyin/consts"
	"douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"go.uber.org/zap"
	"sync"
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
}

func (u *UserFollowFlow) Do() error {
	if err := u.checkNum(); err != nil {
		return err
	}
	if err := u.prepareData(); err != nil {
		return e.FailServerBusy.Err()
	}
	if err := u.updateData(); err != nil {
		return err
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
	var wg sync.WaitGroup
	wg.Add(1)
	relationDao := models.NewRelationDao()
	relationCache := cache.NewRelationCache()
	// 提前预热缓存
	go func() {
		defer wg.Done()
		key := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(u.toUserId))
		if err = relationCache.TTLIsExpiredCache(key); err != nil {
			zap.L().Error("service relation_follow TTLIsExpiredCache method exec fail!", zap.Error(err))
			var ids []int64
			if ids, err = relationDao.QueryUserFollowIds(u.toUserId); err != nil {
				zap.L().Error("service relation_follow QueryUserFollowIds method exec fail!", zap.Error(err))
			}
			err = nil
			if len(ids) > 0 {
				relationCache.SAddResetActionUserFollowOrFollower(key, ids)
			}
		}
	}()
	// 查看关系情况
	if u.isFollow, err = relationDao.IsExistRelation(u.userId, u.toUserId); err != nil {
		zap.L().Error("service relation_follow IsExistRelation method exec fail!", zap.Error(err))
		err = e.FailServerBusy.Err()
	}
	if u.isFollow == 1 {
		err = e.FailRepeatAction.Err()
	}
	wg.Wait()
	return
}

func (u *UserFollowFlow) prepareData() (err error) {
	// 通过缓存查找关系
	key := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(u.toUserId))
	relationCache := cache.NewRelationCache()
	var isFollow2 bool
	if isFollow2, err = relationCache.SIsMemberIsExistRelation(key, u.userId); err == nil {
		if isFollow2 {
			u.isFollow2 = 1
		}
		return
	}
	zap.L().Error("service relation_follow SIsMemberIsExistRelation method exec fail!", zap.Error(err))
	// 如果缓存无效就去数据库查找
	if u.isFollow2, err = models.NewRelationDao().IsExistRelation(u.toUserId, u.userId); err != nil {
		zap.L().Error("service relation_follow IsExistRelation method exec fail!", zap.Error(err))
		err = e.FailServerBusy.Err()
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
	if err = models.NewRelationDao().Action1UserRelation(u.userId, u.toUserId, u.isFollow, u.isFollow2); err != nil {
		zap.L().Error("service relation_follow Action1UserRelation method exec fail!", zap.Error(err))
		err = e.FailServerBusy.Err()
	}
	// 保证缓存一致性，先删除再重置缓存
	go func() {
		key := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(u.userId))
		relationCache := cache.NewRelationCache()
		relationCache.DelCache(key)
		var ids []int64
		if ids, err = models.NewRelationDao().QueryUserFollowIds(u.userId); err != nil {
			zap.L().Error("service relation_follow QueryUserFollowIds method exec fail!", zap.Error(err))
			err = nil
		}
		if len(ids) > 0 {
			relationCache.SAddResetActionUserFollowOrFollower(key, ids)
		}
	}()
	return
}
