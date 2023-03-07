package video

import (
	"douyin/cache"
	"douyin/consts"
	models "douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"sync"
	"time"

	"go.uber.org/zap"
)

type FeedResponse struct {
	NextTime int64           `json:"next_time,omitempty"`
	Videos   []*models.Video `json:"video_list,omitempty"`
}

func UserFeed(lastTime, userId int64) (*FeedResponse, error) {
	return NewVideoUserFeedFlow(lastTime, userId).Do()
}

func NewVideoUserFeedFlow(lastTime, userId int64) *UserFeedFlow {
	return &UserFeedFlow{lastTime: lastTime, userId: userId, videos: make([]*models.Video, consts.MaxFeedVideos)}
}

type UserFeedFlow struct {
	lastTime int64
	userId   int64

	nextTime int64
	videos   []*models.Video

	data *FeedResponse
}

func (u *UserFeedFlow) Do() (*FeedResponse, error) {
	if err := u.checkNum(); err != nil {
		return nil, err
	}
	if err := u.prepareData(); err != nil {
		zap.L().Error("service video_user_feed prepareData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	if err := u.packData(); err != nil {
		zap.L().Error("service video_user_feed packData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	return u.data, nil
}

func (u *UserFeedFlow) checkNum() error {
	if u.userId == 0 {
		return e.FailNotKnow.Err()
	}
	if u.lastTime > time.Now().UnixMilli() {
		zap.L().Debug("service video_user_feed uLastAt", zap.Int64("time", u.lastTime))
		return e.FailNotKnow.Err()
	}
	return nil
}

func (u *UserFeedFlow) prepareData() (err error) {
	// 1. 根据时间查询数据库中视频条数
	videoDao := models.NewVideoDao()
	formatT := time.UnixMilli(u.lastTime).Format("2006-01-02 15:04:05")
	if err = videoDao.QueryVideoListByTime(u.videos, formatT); err != nil {
		zap.L().Error("service video_user_feed QueryVideoListByTime method exec fail!", zap.Error(err))
		return
	}
	if u.videos[0] == nil {
		if err = videoDao.QueryVideoList(u.videos); err != nil {
			zap.L().Error("service video_user_feed QueryVideoList method exec fail!", zap.Error(err))
			return
		}
	}
	// 2. 根据每个视频id去查询用户信息
	userDao := models.NewUserDao()
	relationDao := models.NewRelationDao()
	favorDao := models.NewFavorDao()
	favorCache := cache.NewFavorCache()
	relationCache := cache.NewRelationCache()
	followKey := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, utils.I64toa(u.userId))
	favorKey := utils.AddCacheKey(consts.CacheFavor, consts.CacheSetUserFavor, utils.I64toa(u.userId))
	var wg sync.WaitGroup
	for i, video := range u.videos {
		if video == nil {
			u.videos = u.videos[:i]
			break
		}
		wg.Add(1)
		go func(vdo *models.Video) {
			defer wg.Done()
			// 通过id查询用户信息
			vdo.Author = new(models.User)
			if err = userDao.QueryUserInfoById(vdo.Author, vdo.UserId); err != nil {
				zap.L().Error("service video_user_feed QueryUserInfoById method exec fail!", zap.Error(err))
			}
			// 判断用户关系
			if u.userId != vdo.UserId {
				if vdo.Author.IsFollow, err = relationCache.SIsMemberIsExistRelation(followKey, vdo.UserId); err != nil {
					zap.L().Error("service video_user_feed SIsMemberIsExistRelation method exec fail!", zap.Error(err))
					// 如果缓存无效就去数据库查找
					var isFollow int
					if isFollow, err = relationDao.IsExistRelation(u.userId, vdo.UserId); err != nil {
						zap.L().Error("service video_user_feed NewRelationDao method exec fail!", zap.Error(err))
					}
					if isFollow == 1 {
						vdo.Author.IsFollow = true
					}
				}
			}
			if vdo.IsFavor, err = favorCache.SIsMemberIsExistFavor(favorKey, vdo.VideoId); err != nil {
				zap.L().Error("service video_user_feed SIsMemberIsExistFavor method exec fail!", zap.Error(err))
				// 如果缓存无效就去数据库中找
				var isFavor int
				if isFavor, err = favorDao.IsExistFavor(u.userId, vdo.VideoId); err != nil {
					zap.L().Error("service video_user_feed IsExistFavor method exec fail!", zap.Error(err))
				}
				if isFavor == 1 {
					vdo.IsFavor = true
				}
			}
		}(video)
	}
	wg.Wait()
	return nil
}

func (u *UserFeedFlow) packData() error {
	u.data = &FeedResponse{
		NextTime: u.videos[0].CreateAt.UnixMilli(),
		Videos:   u.videos,
	}
	return nil
}
