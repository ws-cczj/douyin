package cache

import (
	"douyin/consts"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"sync"

	"go.uber.org/zap"
)

type RelationCache struct {
}

var (
	relationCache *RelationCache
	relationOnce  sync.Once
)

func NewRelationCache() *RelationCache {
	relationOnce.Do(func() {
		relationCache = new(RelationCache)
	})
	return relationCache
}

// StringSingleSignOn 限制单点用户登录
//func (*UserCache) StringSingleSignOn() error {return nil}

// SAddResetActionUserFollowOrFollower 用户关注行为缓存重置
func (r *RelationCache) SAddResetActionUserFollowOrFollower(key string, toUserIds []int64) {
	if key != "" {
		pipe := rdbRelation.Pipeline()
		// 填充初始数据 -1
		pipe.SAdd(ctx, key, -1)
		for _, id := range toUserIds {
			pipe.SAdd(ctx, key, id)
		}
		pipe.Expire(ctx, key, consts.CacheExpired)
		if _, err := pipe.Exec(ctx); err != nil {
			zap.L().Error("cache relation SAddResetActionUserFollowOrFollower method exec fail!", zap.Error(err))
			// 如果失败就将缓存进行删除,避免脏数据
			go r.SPopNRemoveCache(key, int64(len(toUserIds)))
		}
	}
}

// SCardQueryUserFollows 查询用户关注数
func (*RelationCache) SCardQueryUserFollows(key string) (follows int64, err error) {
	if follows, err = rdbRelation.SCard(ctx, key).Result(); follows > 0 {
		return follows - 1, nil
	}
	if err != nil {
		zap.L().Error("cache relation SCardQueryUserFollows SCard method exec fail", zap.Error(err))
	}
	return -1, err
}

// SCardQueryUserFollowers 查询用户粉丝数
func (*RelationCache) SCardQueryUserFollowers(key string) (followers int64, err error) {
	if followers, err = rdbRelation.SCard(ctx, key).Result(); followers > 0 {
		return followers - 1, nil
	}
	if err != nil {
		zap.L().Error("cache relation SCardQueryUserFollowers SCard method exec fail", zap.Error(err))
	}
	return -1, err
}

// SIsMemberIsExistRelation 判断是否存在关系
func (*RelationCache) SIsMemberIsExistRelation(key string, toUserId int64) (bool, error) {
	if key == "" {
		return false, e.FailNotKnow.Err()
	}
	return rdbRelation.SIsMember(ctx, key, utils.I64toa(toUserId)).Result()
}

// TTLIsExpiredCache 判断缓存是否过期
func (r *RelationCache) TTLIsExpiredCache(key string) error {
	if key == "" {
		return e.FailNotKnow.Err()
	}
	if t := rdbRelation.TTL(ctx, key).Val(); t < 1 {
		zap.L().Warn("cache relation ttl < 0", zap.String("key", key))
		return e.FailCacheExpired.Err()
	}
	// 如果缓存没有过期就去续约
	go r.ExpireContinueCache(key)
	return nil
}

// DelCache 删除缓存
func (*RelationCache) DelCache(key string) {
	if key != "" {
		var err error
		// 启动错误重试机制，如果删除失败后果比较严重
		for i := 1; i <= consts.CacheMaxTryTimes; i++ {
			if err = rdbRelation.Del(ctx, key).Err(); err == nil {
				return
			}
			zap.L().Warn("cache relation DelCache Del method exec fail!",
				zap.Error(err),
				zap.Int("try times", i))
		}
	}
}

// SPopNRemoveCache 弹出缓存中所有数据
func (*RelationCache) SPopNRemoveCache(key string, cnt int64) {
	if key != "" {
		var err error
		// 启动错误重试机制，如果删除失败后果比较严重
		for i := 1; i <= consts.CacheMaxTryTimes; i++ {
			if err = rdbRelation.SPopN(ctx, key, cnt).Err(); err == nil {
				return
			}
			zap.L().Warn("cache relation SPopNRemoveCache SPopN method exec fail!",
				zap.Error(err),
				zap.Int("try times", i))
		}
	}
}

// ExpireContinueCache  续约缓存
func (*RelationCache) ExpireContinueCache(key string) {
	if key != "" {
		rdbRelation.Expire(ctx, key, consts.CacheExpired)
	}
}
