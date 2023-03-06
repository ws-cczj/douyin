package cache

import (
	"douyin/consts"
	"douyin/pkg/e"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type VideoCache struct {
}

var (
	videoCache *VideoCache
	videoOnce  sync.Once
)

func NewVideoCache() *VideoCache {
	videoOnce.Do(func() {
		videoCache = new(VideoCache)
	})
	return videoCache
}

// SetEXResetVideoComments 重设视频评论数量
func (v *VideoCache) SetEXResetVideoComments(key string, comments int64) {
	if key != "" {
		if err := rdbVideo.SetEX(ctx, key, comments, consts.CacheExpired).Err(); err != nil {
			zap.L().Error("cache video SetEXResetVideoComments method exec fail!", zap.Error(err))
			go v.DelCache(key)
		}
	}
}

// GetEXVideoComments 获取视频评论数量并且设置缓存时间
func (v *VideoCache) GetEXVideoComments(key string) (comments int64, err error) {
	if key == "" {
		return -1, e.FailNotKnow.Err()
	}
	pipe := rdbVideo.TxPipeline()
	pipe.Get(ctx, key)
	pipe.Expire(ctx, key, consts.CacheExpired)
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		if err == redis.Nil {
			zap.L().Warn("cache video key is not exist!", zap.String("key", key))
		} else {
			zap.L().Error("cache video GetEx method exec fail!", zap.Error(err))
		}
		return
	}
	return cmders[0].(*redis.StringCmd).Int64()
}

// DelCache 删除缓存
func (*VideoCache) DelCache(key string) {
	if key != "" {
		var err error
		// 启动错误重试机制，如果删除失败后果比较严重
		for i := 1; i <= consts.CacheMaxTryTimes; i++ {
			if err = rdbVideo.Del(ctx, key).Err(); err == nil {
				return
			}
			zap.L().Warn("cache video DelCache Del method exec fail!",
				zap.Error(err),
				zap.Int("try times", i))
		}
	}
}

// TTLIsExpiredCache 判断缓存是否过期
func (*VideoCache) TTLIsExpiredCache(key string) error {
	if key == "" {
		return e.FailNotKnow.Err()
	}
	if t := rdbVideo.TTL(ctx, key).Val(); t < 1 {
		zap.L().Error("cache video ttl < 0", zap.String("key", key))
		return e.FailCacheExpired.Err()
	}
	// 如果缓存没有过期就去续约
	go rdbVideo.Expire(ctx, key, consts.CacheExpired)
	return nil
}
