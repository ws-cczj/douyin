package cache

import (
	"douyin/consts"
	"douyin/pkg/utils"
	"sync"

	"go.uber.org/zap"

	"github.com/go-redis/redis/v8"
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

// SAddRegisterActionUserFollowAndFollower 注册用户关注行为
func (*RelationCache) SAddRegisterActionUserFollowAndFollower(userId int64) (err error) {
	userFollowKey := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(userId))
	userFollowerKey := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollower, utils.I64toa(userId))

	// 启动重试机制
	for i := 0; i < consts.CacheMaxTryTimes; i++ {
		if _, err = rdbRelation.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.SAdd(ctx, userFollowKey, -1)
			pipe.SAdd(ctx, userFollowerKey, -1)
			pipe.Expire(ctx, userFollowKey, consts.CacheExpired)
			pipe.Expire(ctx, userFollowerKey, consts.CacheExpired)
			return nil
		}); err == nil {
			break
		}
		zap.L().Error("cache relation SAddRegisterActionUserFollowAndFollower method exec fail!",
			zap.Error(err),
			zap.Int("try again times", i))
	}
	return
}

// SAddActionUserFollowAndFollower 用户关注行为
func (*RelationCache) SAddActionUserFollowAndFollower(userId, toUserId int64) (err error) {
	userFollowKey := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(userId))
	userFollowerKey := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollower, utils.I64toa(toUserId))

	if _, err = rdbRelation.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SAdd(ctx, userFollowKey, toUserId)
		pipe.SAdd(ctx, userFollowerKey, userId)
		pipe.Expire(ctx, userFollowKey, consts.CacheExpired)
		pipe.Expire(ctx, userFollowerKey, consts.CacheExpired)
		return nil
	}); err != nil {
		zap.L().Error("cache relation SAddActionUserFollowAndFollower method exec fail!", zap.Error(err))
	}
	return
}

// SRemActionUserFollowAndFollower 用户取关行为
func (*RelationCache) SRemActionUserFollowAndFollower(userId, toUserId int64, userFollowKey, userFollowerKey string) (err error) {
	_, err = rdbRelation.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SRem(ctx, userFollowKey, toUserId)
		pipe.SRem(ctx, userFollowerKey, userId)
		pipe.Expire(ctx, userFollowKey, consts.CacheExpired)
		pipe.Expire(ctx, userFollowerKey, consts.CacheExpired)
		return nil
	})
	return
}

// SAddMoreActionUserFollowAndFollower 用户多次关注行为缓存重置
func (*RelationCache) SAddMoreActionUserFollowAndFollower(key string, toUserIds []int64) {
	pipe := rdbRelation.Pipeline()
	// 填充初始数据 -1
	pipe.SAdd(ctx, key, -1)
	for _, id := range toUserIds {
		pipe.SAdd(ctx, key, id)
	}
	pipe.Expire(ctx, key, consts.CacheExpired)
	if _, err := pipe.Exec(ctx); err != nil {
		zap.L().Error("cache relation SAddMoreActionUserFollow method exec fail!", zap.Error(err))
	}
}

// SCardQueryUserFollows 查询用户关注数
func (*RelationCache) SCardQueryUserFollows(userId int64) (follows int64, err error) {
	userFollowKey := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(userId))
	if follows, err = rdbRelation.SCard(ctx, userFollowKey).Result(); follows > 0 {
		go rdbRelation.Expire(ctx, userFollowKey, consts.CacheExpired)
		return follows - 1, nil
	}
	if err != nil {
		zap.L().Error("cache relation SCardQueryUserFollows SCard method exec fail", zap.Error(err))
	}
	return -1, err
}

// SCardQueryUserFollowers 查询用户粉丝数
func (*RelationCache) SCardQueryUserFollowers(userId int64) (followers int64, err error) {
	userFollowerKey := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollower, utils.I64toa(userId))
	if followers, err = rdbRelation.SCard(ctx, userFollowerKey).Result(); followers > 0 {
		go rdbRelation.Expire(ctx, userFollowerKey, consts.CacheExpired)
		return followers - 1, nil
	}
	if err != nil {
		zap.L().Error("cache relation SCardQueryUserFollowers SCard method exec fail", zap.Error(err))
	}
	return -1, err
}

// TTLIsExpiredCache 判断缓存是否过期
func (*RelationCache) TTLIsExpiredCache(keys ...string) ([]bool, error) {
	cmders, err := rdbRelation.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.TTL(ctx, key)
		}
		return nil
	})
	oks := make([]bool, len(keys))
	for _, cmder := range cmders {
		oks = append(oks, cmder.(*redis.DurationCmd).Val() > 0)
	}
	return oks, err
}
