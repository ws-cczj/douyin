package cache

import (
	"douyin/consts"
	"douyin/pkg/e"
	"sync"

	"go.uber.org/zap"
)

type FavorCache struct {
}

var (
	favorCache *FavorCache
	favorOnce  sync.Once
)

func NewFavorCache() *FavorCache {
	favorOnce.Do(func() {
		favorCache = new(FavorCache)
	})
	return favorCache
}

// SCardQueryUserFavorVideos 获取用户喜欢的视频数
func (*FavorCache) SCardQueryUserFavorVideos(key string) (favors int64, err error) {
	if key == "" {
		return -1, e.FailNotKnow.Err()
	}
	if favors, err = rdbFavor.SCard(ctx, key).Result(); favors > 0 {
		// 剪掉初始数据 -1
		return favors - 1, nil
	}
	if err != nil {
		zap.L().Error("cache favor SCardQueryUserFavorVideos method exec fail!", zap.Error(err))
	}
	return -1, err
}

// SAddReSetUserFavorVideo 重设用户点赞视频缓存
func (f *FavorCache) SAddReSetUserFavorVideo(key string, videoIds []int64) {
	if key != "" {
		pipe := rdbFavor.Pipeline()
		pipe.SAdd(ctx, key, -1)
		for _, id := range videoIds {
			pipe.SAdd(ctx, key, id)
		}
		pipe.Expire(ctx, key, consts.CacheExpired)
		if _, err := pipe.Exec(ctx); err != nil {
			zap.L().Error("cache favor SAddReSetUserFavorVideo method exec fail!", zap.Error(err))
			go f.SPopNRemoveCache(key, int64(len(videoIds)))
		}
	}
}

// SIsMemberIsExistFavor 是否存在点赞
func (*FavorCache) SIsMemberIsExistFavor(key string, videoId int64) (bool, error) {
	if key == "" {
		return false, e.FailNotKnow.Err()
	}
	return rdbFavor.SIsMember(ctx, key, videoId).Result()
}

// DelCache 删除缓存
func (*FavorCache) DelCache(key string) {
	if key != "" {
		var err error
		// 启动错误重试机制，如果删除失败后果比较严重
		for i := 1; i <= consts.CacheMaxTryTimes; i++ {
			if err = rdbFavor.Del(ctx, key).Err(); err == nil {
				return
			}
			zap.L().Warn("cache favor DelCache Del method exec fail!",
				zap.Error(err),
				zap.Int("try times", i))
		}
	}
}

// SPopNRemoveCache 弹出缓存中所有数据
func (*FavorCache) SPopNRemoveCache(key string, cnt int64) {
	if key != "" {
		var err error
		// 启动错误重试机制，如果删除失败后果比较严重
		for i := 1; i <= consts.CacheMaxTryTimes; i++ {
			if err = rdbFavor.SPopN(ctx, key, cnt).Err(); err == nil {
				return
			}
			zap.L().Warn("cache favor SPopNRemoveCache SPopN method exec fail!",
				zap.Error(err),
				zap.Int("try times", i))
		}
	}
}

// TTLIsExpiredCache 判断缓存是否过期
func (*FavorCache) TTLIsExpiredCache(key string) error {
	if key == "" {
		return e.FailNotKnow.Err()
	}
	if t := rdbFavor.TTL(ctx, key).Val(); t < 1 {
		zap.L().Warn("cache favor ttl < 0", zap.String("key", key))
		return e.FailCacheExpired.Err()
	}
	// 如果缓存没有过期就去续约
	go rdbFavor.Expire(ctx, key, consts.CacheExpired)
	return nil
}
