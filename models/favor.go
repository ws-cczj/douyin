package models

import (
	"database/sql"
	"database/sql/driver"
	"douyin/pkg/utils"
	"sync"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

type Favor struct {
	IsFavor int   `json:"is_favor"`
	Id      int64 `json:"id"`
	UserId  int64 `json:"user_id"`
	VideoId int64 `json:"video_id"`
}

func (f Favor) Value() (driver.Value, error) {
	return []interface{}{f.VideoId, f.IsFavor}, nil
}

type FavorDao struct {
}

var (
	favorDao  *FavorDao
	favorOnce sync.Once
)

// NewFavorDao 使用饿汉式单例模式初始化FavorDao对象
func NewFavorDao() *FavorDao {
	favorOnce.Do(func() {
		favorDao = new(FavorDao)
	})
	return favorDao
}

// QueryUserFavorVideos 查询用户点赞的视频总数
func (*FavorDao) QueryUserFavorVideos(userId int64) (favors int64, err error) {
	qStr := `select Count(*) from user_favor_videos where user_id = ? AND is_favor = ?`
	if err = db.GetContext(ctx, &favors, qStr, userId, 1); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("models favor QueryUserFavorVideos data is null", zap.Error(err))
			err = nil
		}
	}
	return
}

// QueryUserVideosFavors 查询用户每个视频的获赞数
func (*FavorDao) QueryUserVideosFavors(videoIds []int64) (favors []int64, err error) {
	qStr := `select Count(*)
        		from user_favor_videos
        		where video_id in(?) AND is_favor = ?
        		order by FIND_IN_SET(video_id, ?);`
	favors = make([]int64, len(videoIds))
	qry, args, _ := sqlx.In(qStr, videoIds, utils.ISlice64toa(videoIds))
	query := db.Rebind(qry)
	if err = db.SelectContext(ctx, &favors, query, args...); err != nil {
		zap.L().Error("models favor QueryUserVideosFavors method exec fail!", zap.Error(err))
	}
	return
}

// QueryUserFavorVideoList 查询用户点赞的视频列表
func (*FavorDao) QueryUserFavorVideoList(userId int64) (favorsVideos []int64, err error) {
	qStr := `select video_id from user_favor_videos where user_id = ? AND is_favor = ?`
	favorsVideos = make([]int64, 0)
	if err = db.SelectContext(ctx, &favorsVideos, qStr, userId, 1); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("models favor QueryUserFavorVideoList data is null", zap.Error(err))
			err = nil
		}
	}
	return
}
