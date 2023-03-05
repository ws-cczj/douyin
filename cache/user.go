package cache

import (
	"douyin/consts"
	"douyin/pkg/utils"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type UserCache struct {
}

var (
	userCache *UserCache
	userOnce  sync.Once
)

func NewUserCache() *UserCache {
	userOnce.Do(func() {
		userCache = new(UserCache)
	})
	return userCache
}

// SAddPublishVideo 发布视频
func (*UserCache) SAddPublishVideo(userId, videoId int64) (err error) {
	userVideoKey := utils.AddCacheKey(consts.CacheUser, consts.CacheSetUserVideo, utils.I64toa(userId))

	_, err = rdbUser.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SAdd(ctx, userVideoKey, videoId)
		pipe.Expire(ctx, userVideoKey, consts.CacheExpired)
		return nil
	})
	return
}

// SMembersQueryUserVideoList 查询用户发布的视频列表
func (*UserCache) SMembersQueryUserVideoList(userId int64) ([]int64, error) {
	userVideoKey := utils.AddCacheKey(consts.CacheUser, consts.CacheSetUserVideo, utils.I64toa(userId))
	cmders, err := rdbUser.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SMembers(ctx, userVideoKey)
		pipe.Expire(ctx, userVideoKey, consts.CacheExpired)
		return nil
	})
	if err != nil {
		zap.L().Error("cache user SMembersUserVideoList method exec fail!", zap.Error(err))
		return nil, err
	}
	vals := cmders[0].(*redis.StringSliceCmd).Val()
	videoList := make([]int64, len(vals))
	for _, val := range vals {
		videoList = append(videoList, utils.AtoI64(val))
	}
	return videoList, nil
}

// SAddReSetUserVideoList 用户发布视频缓存重置
func (u *UserCache) SAddReSetUserVideoList(key string, videoIds []int64) {
	if key != "" {
		pipe := rdbRelation.Pipeline()
		// 填充初始数据 -1
		pipe.SAdd(ctx, key, -1)
		for _, id := range videoIds {
			pipe.SAdd(ctx, key, id)
		}
		pipe.Expire(ctx, key, consts.CacheExpired)
		if _, err := pipe.Exec(ctx); err != nil {
			zap.L().Error("cache user SAddMoreUserVideoList method exec fail!", zap.Error(err))
			u.SPopNRemoveCache(key, int64(len(videoIds)))
		}
	}
}

// DelCache 删除缓存
func (*UserCache) DelCache(key string) {
	if key != "" {
		var err error
		// 启动错误重试机制，如果删除失败后果比较严重
		for i := 1; i <= consts.CacheMaxTryTimes; i++ {
			if err = rdbUser.Del(ctx, key).Err(); err == nil {
				return
			}
			zap.L().Error("cache user DelCache Del method exec fail!",
				zap.Error(err),
				zap.Int("try times", i))
		}
	}
}

// SPopNRemoveCache 弹出缓存中所有数据
func (*UserCache) SPopNRemoveCache(key string, cnt int64) {
	if key != "" {
		var err error
		// 启动错误重试机制，如果删除失败后果比较严重
		for i := 1; i <= consts.CacheMaxTryTimes; i++ {
			if err = rdbUser.SPopN(ctx, key, cnt).Err(); err == nil {
				return
			}
			zap.L().Error("cache user SPopNRemoveCache SPopN method exec fail!",
				zap.Error(err),
				zap.Int("try times", i))
		}
	}
}

// ExpireContinueCache  续约缓存
func (*UserCache) ExpireContinueCache(key string) {
	if key != "" {
		rdbUser.Expire(ctx, key, consts.CacheExpired)
	}
}
