package user

import (
	"database/sql"
	"douyin/cache"
	"douyin/consts"
	"douyin/models"
	"douyin/pkg/utils"
	"errors"

	"go.uber.org/zap"
)

type InfoResponse struct {
	FollowCount   int64        `json:"follow_count"`
	FollowerCount int64        `json:"follower_count"`
	WorkCount     int64        `json:"work_count"`
	FavorCount    int64        `json:"favorite_count"`
	TotalFavor    int64        `json:"total_favorited,string"`
	User          *models.User `json:"user"`
	IsFollow      bool         `json:"is_follow"`
}

func Info(userId, tkUserId int64) (*InfoResponse, error) {
	return NewUserInfoFlow(userId, tkUserId).Do()
}

func NewUserInfoFlow(userId, tkUserId int64) *InfoFlow {
	return &InfoFlow{userId: userId, tkUserId: tkUserId, user: new(models.User)}
}

type InfoFlow struct {
	userId   int64
	tkUserId int64

	followCount   int64
	followerCount int64
	workCount     int64
	favorCount    int64
	totalFavor    int64
	user          *models.User
	isFollow      bool

	data *InfoResponse
}

func (i *InfoFlow) Do() (*InfoResponse, error) {
	if err := i.checkNum(); err != nil {
		zap.L().Error("service user_info checkNum method exec fail!", zap.Error(err))
		return nil, err
	}
	if err := i.prepareData(); err != nil {
		return nil, err
	}
	if err := i.packData(); err != nil {
		return nil, err
	}
	return i.data, nil
}

func (i *InfoFlow) checkNum() (err error) {
	if i.userId == 0 || i.tkUserId == 0 {
		return errors.New("用户信息错误")
	}
	// 根据User_id查询数据库获取User信息。
	if err = models.NewUserDao().QueryUserInfoById(i.user, i.userId); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("无法查询到该用户")
		}
	}
	return
}

func (i *InfoFlow) prepareData() (err error) {
	// 1. 查看缓存中是否有该用户的Follow、Follower
	if i.followCount, err = cache.NewRelationCache().SCardQueryUserFollows(i.userId); i.followCount < 0 {
		if err != nil {
			zap.L().Error("service user_info SCardQueryUserFollows method exec fail!", zap.Error(err))
		}
		var ids []int64
		// 要么缓存过期，要么执行错误, 都会返回数据 -1
		if ids, err = models.NewRelationDao().QueryUserFollowList(i.userId); err != nil {
			zap.L().Error("handlers user_info QueryUserFollowList method exec fail!", zap.Error(err))
			return
		}
		i.followCount = int64(len(ids))
		// 使用协程重置缓存
		go func() {
			key := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(i.userId))
			cache.NewRelationCache().SAddMoreActionUserFollowAndFollower(key, ids)
		}()
	}
	if i.followerCount, err = cache.NewRelationCache().SCardQueryUserFollowers(i.userId); i.followerCount < 0 {
		if err != nil {
			zap.L().Error("service user_info SCardQueryUserFollowers method exec fail!", zap.Error(err))
		}
		var ids []int64
		if ids, err = models.NewRelationDao().QueryUserFollowerList(i.userId); err != nil {
			zap.L().Error("handlers user_info QueryUserFollowerList method exec fail!", zap.Error(err))
			return
		}
		// 判断当前用户与查询用户之间的关系
		for _, id := range ids {
			if id == i.tkUserId {
				i.isFollow = true
			}
		}
		i.followerCount = int64(len(ids))
		// 使用协程重置缓存
		go func() {
			key := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollower, utils.I64toa(i.userId))
			cache.NewRelationCache().SAddMoreActionUserFollowAndFollower(key, ids)
		}()
	}
	// 用户发表的视频列表
	var videoIds []int64
	if videoIds, err = cache.NewUserCache().SMembersQueryUserVideoList(i.userId); err != nil {
		zap.L().Error("service user_info SCardQueryUserFollows method exec fail!", zap.Error(err))
		if videoIds, err = models.NewVideoDao().QueryUserVideoList(i.userId); err != nil {
			zap.L().Error("service user_info QueryUserVideoList method exec fail!", zap.Error(err))
			return
		}
		go cache.NewUserCache().SAddMoreUserVideoList(i.userId, videoIds)
	}
	if i.workCount = int64(len(videoIds)); i.workCount > 0 {
		if i.totalFavor, err = cache.NewVideoCache().StringQueryVideosFavors(videoIds); i.totalFavor == -1 {
			if err != nil {
				zap.L().Error("service user_info StringQueryVideosFavors method exec fail!", zap.Error(err))
			}
			var favors []int64
			if favors, err = models.NewFavorDao().QueryUserVideosFavors(videoIds); err != nil {
				zap.L().Error("service user_info QueryUserVideosFavors method exec fail!", zap.Error(err))
				return
			}
			go cache.NewVideoCache().StringReSetVideosFavors(videoIds, favors)
		}
	}
	if i.favorCount, err = cache.NewFavorCache().SCardQueryUserFavorVideos(i.userId); i.favorCount < 0 {
		if err != nil {
			zap.L().Error("service user_info SCardQueryUserFollows method exec fail!", zap.Error(err))
		}
		var favorVideos []int64
		if favorVideos, err = models.NewFavorDao().QueryUserFavorVideoList(i.userId); err != nil {
			zap.L().Error("service user_info QueryUserFavorVideos method exec fail!", zap.Error(err))
			return
		}
		i.favorCount = int64(len(favorVideos))
		go cache.NewFavorCache().SAddReSetUserFavorVideo(i.userId, favorVideos)
	}
	return
}

func (i *InfoFlow) packData() error {
	i.data = &InfoResponse{
		FollowCount:   i.followerCount,
		FollowerCount: i.followerCount,
		FavorCount:    i.favorCount,
		WorkCount:     i.workCount,
		User:          i.user,
		IsFollow:      i.isFollow,
		TotalFavor:    i.totalFavor,
	}
	return nil
}
