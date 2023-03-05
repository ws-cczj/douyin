package cache

import (
	"douyin/consts"
	"douyin/pkg/e"
	"douyin/pkg/utils"
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

// StringIncrVideoFavor 增加视频点赞数
func (*VideoCache) StringIncrVideoFavor(videoId int64) {
	videoFavorKey := utils.AddCacheKey(consts.CacheVideo, consts.CacheStringVideoFavor, utils.I64toa(videoId))
	// 使用缓存重试机制
	for i := 0; i < consts.CacheMaxTryTimes; i++ {
		if err := rdbVideo.Incr(ctx, videoFavorKey).Err(); err != nil {
			zap.L().Error("cache video StringIncrVideoFavor method exec fail!",
				zap.Error(err),
				zap.Int("try again times", i))
			continue
		}
		go rdbVideo.Expire(ctx, videoFavorKey, consts.CacheExpired)
		break
	}
}

// StringDecrVideoFavor 减少视频点赞数
func (*VideoCache) StringDecrVideoFavor(videoId int64) {
	videoFavorKey := utils.AddCacheKey(consts.CacheVideo, consts.CacheStringVideoFavor, utils.I64toa(videoId))
	// 使用缓存重试机制
	for i := 0; i < consts.CacheMaxTryTimes; i++ {
		if err := rdbVideo.Decr(ctx, videoFavorKey).Err(); err != nil {
			zap.L().Error("cache video StringDecrVideoFavor method exec fail!",
				zap.Error(err),
				zap.Int("try again times", i))
			continue
		}
		go rdbVideo.Expire(ctx, videoFavorKey, consts.CacheExpired)
		break
	}
}

// StringQueryVideoFavors 查询该视频点赞数
func (*VideoCache) StringQueryVideoFavors(videoId int64) (int64, error) {
	videoFavorKey := utils.AddCacheKey(consts.CacheVideo, consts.CacheStringVideoFavor, utils.I64toa(videoId))
	return rdbVideo.GetEx(ctx, videoFavorKey, consts.CacheExpired).Int64()
}

// StringQueryVideosFavors 查询指定视频的总点赞数
func (*VideoCache) StringQueryVideosFavors(videoIds []int64) (res int64, err error) {
	pipe := rdbVideo.Pipeline()
	for _, videoId := range videoIds {
		videoFavorKey := utils.AddCacheKey(consts.CacheVideo, consts.CacheStringVideoFavor, utils.I64toa(videoId))
		pipe.GetEx(ctx, videoFavorKey, consts.CacheExpired)
	}
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		zap.L().Error("cache video StringQueryVideosFavors method exec fail", zap.Error(err))
		return -1, err
	}
	for _, cmder := range cmders {
		res += cmder.(*redis.IntCmd).Val()
	}
	return
}

// StringReSetVideosFavors 重新设置每个视频的点赞数量
func (*VideoCache) StringReSetVideosFavors(videoIds []int64, favors []int64) {
	pipe := rdbVideo.Pipeline()
	for i, id := range videoIds {
		videoFavorKey := utils.AddCacheKey(consts.CacheVideo, consts.CacheStringVideoFavor, utils.I64toa(id))
		pipe.SetEX(ctx, videoFavorKey, favors[i], consts.CacheExpired)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		zap.L().Error("cache video StringReSetVideosFavors method exec fail!", zap.Error(err))
	}
}

// SetEXResetVideoComments 重设视频评论列表
func (v *VideoCache) SetEXResetVideoComments(key string, comments int64) {
	if key != "" {
		if err := rdbVideo.SetEX(ctx, key, comments, consts.CacheExpired).Err(); err != nil {
			zap.L().Error("cache video SetEXResetVideoComments method exec fail!", zap.Error(err))
			go v.DelCache(key)
		}
	}
}

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
			err = nil
		}
		zap.L().Error("cache video GetEx method exec fail!", zap.Error(err))
		return
	}
	return cmders[0].(*redis.IntCmd).Result()
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
