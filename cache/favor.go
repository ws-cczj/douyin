package cache

import (
	"douyin/consts"
	"douyin/pkg/utils"
	"sync"

	"github.com/go-redis/redis/v8"
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

// SAddUserFavorVideo 用户点赞视频
func (*FavorCache) SAddUserFavorVideo(userId, videoId int64) {
	userFavorKey := utils.AddCacheKey(consts.CacheUser, consts.CacheSetUserFavor, utils.I64toa(userId))
	// 使用缓存重试机制
	for i := 0; i < consts.CacheMaxTryTimes; i++ {
		pipe := rdbFavor.TxPipeline()
		pipe.SAdd(ctx, userFavorKey, videoId)
		pipe.Expire(ctx, userFavorKey, consts.CacheExpired)
		if _, err := pipe.Exec(ctx); err != nil {
			zap.L().Error("cache favor SAddUserFavorVideo method exec fail!",
				zap.Error(err),
				zap.Int("try again times", i))
			continue
		}
		break
	}
}

// SCardQueryUserFavorVideos 获取用户喜欢的视频数
func (*FavorCache) SCardQueryUserFavorVideos(userId int64) (favors int64, err error) {
	userFavorKey := utils.AddCacheKey(consts.CacheUser, consts.CacheSetUserFavor, utils.I64toa(userId))
	if favors, err = rdbFavor.SCard(ctx, userFavorKey).Result(); favors > 0 {
		go rdbFavor.Expire(ctx, userFavorKey, consts.CacheExpired)
		return favors - 1, nil
	}
	if err != nil {
		zap.L().Error("cache favor SCardQueryUserFavorVideos method exec fail!", zap.Error(err))
	}
	return -1, err
}

// SAddReSetUserFavorVideo 重设用户点赞视频缓存
func (*FavorCache) SAddReSetUserFavorVideo(userId int64, videoId []int64) {
	userFavorKey := utils.AddCacheKey(consts.CacheUser, consts.CacheSetUserFavor, utils.I64toa(userId))
	pipe := rdbFavor.Pipeline()
	pipe.SAdd(ctx, userFavorKey, -1)
	for _, id := range videoId {
		pipe.SAdd(ctx, userFavorKey, id)
	}
	pipe.Expire(ctx, userFavorKey, consts.CacheExpired)
	if _, err := pipe.Exec(ctx); err != nil {
		zap.L().Error("cache favor SAddReSetUserFavorVideo method exec fail!", zap.Error(err))
	}
}

// SMembersQueryUserFavorVideoList 查询用户喜欢视频列表
func (*FavorCache) SMembersQueryUserFavorVideoList(userId int64) ([]int64, error) {
	userFavorKey := utils.AddCacheKey(consts.CacheUser, consts.CacheSetUserFavor, utils.I64toa(userId))
	pipe := rdbFavor.Pipeline()
	pipe.SMembers(ctx, userFavorKey)
	pipe.Expire(ctx, userFavorKey, consts.CacheExpired)
	cmders, err := pipe.Exec(ctx)
	if err != nil {
		zap.L().Error("cache favor SMembersQueryUserFavorVideoList method exec fail!", zap.Error(err))
		return nil, err
	}
	videoList := make([]int64, len(cmders))
	for _, cmder := range cmders[0].(*redis.SliceCmd).Val() {
		videoList = append(videoList, cmder.(*redis.IntCmd).Val())
	}
	return videoList, nil
}
