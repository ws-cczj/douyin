package cache

import (
	"context"
	"douyin/conf"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var (
	rdbUser     *redis.Client
	rdbRelation *redis.Client
	rdbVideo    *redis.Client
	rdbFavor    *redis.Client
	rdbComment  *redis.Client
)

// InitRedis 初始化所有Redis连接。
func InitRedis() {
	var err error
	defer func() {
		if err != nil {
			panic(fmt.Errorf("initRedis fail!, err: %v", err))
		}
	}()
	rdbUser = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.RDB.Addr,
		Password: conf.Conf.RDB.Password,
		PoolSize: conf.Conf.RDB.PoolSize,
		DB:       conf.Conf.RDB.UserDB, // 用户发布的视频缓存
	})
	err = rdbUser.Ping(ctx).Err()
	rdbRelation = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.RDB.Addr,
		Password: conf.Conf.RDB.Password,
		PoolSize: conf.Conf.RDB.PoolSize,
		DB:       conf.Conf.RDB.RelationDB, // 用户与用户之间的关系缓存
	})
	err = rdbRelation.Ping(ctx).Err()
	rdbVideo = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.RDB.Addr,
		Password: conf.Conf.RDB.Password,
		PoolSize: conf.Conf.RDB.PoolSize,
		DB:       conf.Conf.RDB.VideoDB, // 视频点赞数，视频评论数等缓存
	})
	err = rdbVideo.Ping(ctx).Err()
	rdbFavor = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.RDB.Addr,
		Password: conf.Conf.RDB.Password,
		PoolSize: conf.Conf.RDB.PoolSize,
		DB:       conf.Conf.RDB.FavorDB, // 视频评论缓存
	})
	err = rdbFavor.Ping(ctx).Err()
	rdbComment = redis.NewClient(&redis.Options{
		Addr:     conf.Conf.RDB.Addr,
		Password: conf.Conf.RDB.Password,
		PoolSize: conf.Conf.RDB.PoolSize,
		DB:       conf.Conf.RDB.CommentDB, // 视频评论缓存
	})
	err = rdbComment.Ping(ctx).Err()
}

// Close 统一关闭redis连接
func Close() {
	_ = rdbUser.Close()
	_ = rdbRelation.Close()
	_ = rdbVideo.Close()
	_ = rdbComment.Close()
}
