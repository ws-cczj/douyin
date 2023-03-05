package user

import (
	"douyin/cache"
	"douyin/consts"
	"douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"go.uber.org/zap"
	"sync"
)

func FavorVideoList(userId, tkUserId int64) ([]*models.Video, error) {
	return NewFavorVideoListFlow(userId, tkUserId).Do()
}

func NewFavorVideoListFlow(userId, tkUserId int64) *FavorVideoListFlow {
	return &FavorVideoListFlow{userId: userId, tkUserId: tkUserId}
}

type FavorVideoListFlow struct {
	userId, tkUserId int64

	favorKey string

	data []*models.Video
}

func (f *FavorVideoListFlow) Do() ([]*models.Video, error) {
	if err := f.checkNum(); err != nil {
		return nil, err
	}
	if err := f.prepareData(); err != nil {
		zap.L().Error("service user_favor_video_list packData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	if len(f.data) > 0 {
		if err := f.packData(); err != nil {
			zap.L().Error("service user_favor_video_list packData method exec fail!", zap.Error(err))
			return nil, e.FailServerBusy.Err()
		}
	}
	return f.data, nil
}

func (f *FavorVideoListFlow) checkNum() error {
	// 这里有两种情况，token id可能为空，因此不做判断！
	if f.userId == 0 {
		return e.FailServerBusy.Err()
	}
	// 预热点赞缓存, 这里预热的是目标用户的点赞缓存, 并不是当前用户
	f.favorKey = utils.AddCacheKey(consts.CacheFavor, consts.CacheSetUserFavor, utils.I64toa(f.userId))
	if err := cache.NewFavorCache().TTLIsExpiredCache(f.favorKey); err != nil {
		zap.L().Warn("service user_favor_video_list TTLIsExpiredCache method exec fail!", zap.Error(err))
		var videos []int64
		if videos, err = models.NewFavorDao().QueryUserFavorVideoList(f.userId); err != nil {
			zap.L().Error("service user_favor_video_list QueryUserFavorVideoList method exec fail!", zap.Error(err))
		}
		cache.NewFavorCache().SAddReSetUserFavorVideo(f.favorKey, videos)
	}
	return nil
}

func (f *FavorVideoListFlow) prepareData() (err error) {
	// 本用户要么是游客 tk = 0 要么是用户 tk = ?
	// 目标用户要么是本人 tk = user 要么是目标对象 tk != user
	// 需要判断当前用户与这些视频的作者是否关注以及视频是否进行了点赞

	// 首先需要拿到目标用户的点赞视频数,然后拿到点赞视频列表
	var favors int64
	if favors, err = cache.NewFavorCache().SCardQueryUserFavorVideos(f.favorKey); favors < 0 {
		zap.L().Warn("service user_favor_video_list SCardQueryUserFavorVideos method exec fail!", zap.Error(err))
		if favors, err = models.NewUserDao().QueryUserFavorVideos(f.userId); err != nil {
			zap.L().Error("service user_favor_video_list QueryUserFavorVideos method exec fail!", zap.Error(err))
			return
		}
	}
	// 如果没有喜欢直接返回
	if favors == 0 {
		return
	}
	f.data = make([]*models.Video, favors)
	if err = models.NewVideoDao().QueryVideoListWithFavors(f.data, f.userId); err != nil {
		zap.L().Error("service user_favor_video_list QueryVideoListWithFavors method exec fail!", zap.Error(err))
	}
	return
}

func (f *FavorVideoListFlow) packData() (err error) {
	userDao := models.NewUserDao()
	favorDao := models.NewFavorDao()
	relationDao := models.NewRelationDao()
	favorCache := cache.NewFavorCache()
	relationCache := cache.NewRelationCache()

	// 根据视频进行遍历填充数据
	// 这里需要注意，现在查找的是当前用户对与这些视频的关系，也就是 tkUserId to videoId
	tkUStr := utils.I64toa(f.tkUserId)
	followKey := utils.AddCacheKey(consts.CacheRelation, consts.CacheSetUserFollow, tkUStr)
	favorKey := utils.AddCacheKey(consts.CacheFavor, consts.CacheSetUserFavor, tkUStr)

	var wg sync.WaitGroup
	wg.Add(len(f.data))
	for _, video := range f.data {
		go func(vdo *models.Video) {
			defer wg.Done()
			// 通过id查询用户信息
			vdo.Author = new(models.User)
			if err = userDao.QueryUserInfoById(vdo.Author, vdo.UserId); err != nil {
				zap.L().Error("service user_favor_video_list QueryUserInfoById method exec fail!", zap.Error(err))
			}
			// 如果不是游客访问
			if f.tkUserId != 0 {
				if f.tkUserId != vdo.UserId {
					// 查找当前用户与视频作者之间的关系
					if vdo.Author.IsFollow, err = relationCache.SIsMemberIsExistRelation(followKey, vdo.UserId); err != nil {
						zap.L().Error("service user_favor_video_list SIsMemberIsExistRelation method exec fail!", zap.Error(err))
						// 如果缓存无效就去数据库查找
						var isFollow int
						if isFollow, err = relationDao.IsExistRelation(f.tkUserId, vdo.UserId); err != nil {
							zap.L().Error("service user_favor_video_list NewRelationDao method exec fail!", zap.Error(err))
						}
						if isFollow == 1 {
							vdo.Author.IsFollow = true
						}
					}
				}
				// 通过缓存查找点赞
				if vdo.IsFavor, err = favorCache.SIsMemberIsExistFavor(favorKey, vdo.VideoId); err != nil {
					zap.L().Error("service user_favor_video_list SIsMemberIsExistFavor method exec fail!", zap.Error(err))
					// 如果缓存无效就去数据库中找
					var isFavor int
					if isFavor, err = favorDao.IsExistFavor(f.tkUserId, vdo.VideoId); err != nil {
						zap.L().Error("service user_favor_video_list IsExistFavor method exec fail!", zap.Error(err))
					}
					if isFavor == 1 {
						vdo.IsFavor = true
					}
				}
			}
		}(video)
	}
	wg.Wait()
	return
}
