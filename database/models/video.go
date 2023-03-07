package models

import (
	"database/sql"
	"douyin/consts"
	"douyin/pkg/e"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Video struct {
	VideoId      int64     `json:"id" db:"video_id"`
	UserId       int64     `json:"-" db:"user_id"`
	FavorCount   int64     `json:"favorite_count" db:"favored_count"`
	CommentCount int64     `json:"comment_count" db:"comment_count"`
	Author       *User     `json:"author"`
	Title        string    `json:"title" db:"title"`
	PlayUrl      string    `json:"play_url" db:"play_url"`
	CoverUrl     string    `json:"cover_url" db:"cover_url"`
	IsFavor      bool      `json:"is_favorite"`
	CreateAt     time.Time `json:"-" db:"create_at"`
}

type VideoDao struct {
}

var (
	videoDao  *VideoDao
	videoOnce sync.Once
)

// NewVideoDao 使用饿汉式单例模式初始化VideoDao对象
func NewVideoDao() *VideoDao {
	videoOnce.Do(func() {
		videoDao = new(VideoDao)
	})
	return videoDao
}

// PublishVideo 发布一条视频
func (*VideoDao) PublishVideo(videoId, userId int64, playUrl, coverUrl, title string) (err error) {
	var tx *sql.Tx
	if tx, err = db.Begin(); err == nil {
		if tx == nil {
			zap.L().Error("models relation begin tx transition fail!", zap.Error(err))
			return e.FailServerBusy.Err()
		}
		var wg sync.WaitGroup
		wg.Add(1)
		// 增加一条视频
		go func() {
			defer wg.Done()
			iStr := `insert into videos(video_id,user_id,title,play_url,cover_url) values (?,?,?,?,?)`
			if _, err = db.ExecContext(ctx, iStr, videoId, userId, title, playUrl, coverUrl); err != nil {
				zap.L().Error("models video ExecContext method exec fail!", zap.Error(err))
			}
		}()
		// 增加用户视频数
		uStr := `update users set work_count = work_count + 1 where user_id = ?`
		if _, err = db.ExecContext(ctx, uStr, userId); err != nil {
			zap.L().Error("models video Update User videos fail!", zap.Error(err))
		}
		wg.Wait()
	}
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return
	}
	if err = tx.Commit(); err != nil {
		zap.L().Error("models video tx Commit exec fail!", zap.Error(err))
		tx.Rollback()
	}
	return
}

// QueryUserVideosById 查询用户发布的视频列表ids
func (*VideoDao) QueryUserVideosById(userId int64) (ids []int64, err error) {
	qStr := `select video_id from videos where user_id = ? AND is_delete = ?`
	ids = []int64{}
	if err = db.SelectContext(ctx, &ids, qStr, userId, 0); err != nil {
		zap.L().Error("models video SelectContext method exec fail!", zap.Error(err))
	}
	return
}

// QueryUserVideoListById 查询用户发布的视频列表
func (*VideoDao) QueryUserVideoListById(videos []*Video, userId int64) (err error) {
	qStr := `select video_id,user_id,play_url,cover_url,favored_count,comment_count,title
				from videos
				where user_id = ? AND is_delete = ?`
	if err = db.SelectContext(ctx, &videos, qStr, userId, 0); err != nil {
		zap.L().Error("models video SelectContext method exec fail!", zap.Error(err))
	}
	return
}

// QueryVideoListWithFavors 根据视频数量查询视频列表
func (*VideoDao) QueryVideoListWithFavors(videos []*Video, userId int64) (err error) {
	qStr := `select video_id,user_id,play_url,cover_url,favored_count,comment_count,title
        		from videos
        		where video_id in (select video_id 
					from user_favor_videos 
					where user_id = ? AND is_favor = ?
					) AND is_delete = ?`
	if err = db.SelectContext(ctx, &videos, qStr, userId, 1, 0); err != nil {
		zap.L().Error("models favor SelectContext method exec fail!", zap.Error(err))
	}
	return
}

// QueryVideoListByTime 根据时间来查询视频列表
func (*VideoDao) QueryVideoListByTime(videos []*Video, lastTime string) (err error) {
	qStr := `select video_id,user_id,play_url,cover_url,favored_count,comment_count,title,create_at
				from videos
				where create_at > ? 
  					AND create_at <= Now()
  					AND is_delete = ?
				order by create_at DESC
				limit ?`
	if err = db.SelectContext(ctx, &videos, qStr, lastTime, 0, consts.MaxFeedVideos); err != nil {
		zap.L().Error("models video SelectContext method exec fail!", zap.Error(err))
	}
	return
}

// QueryVideoList 查询视频列表
func (*VideoDao) QueryVideoList(videos []*Video) (err error) {
	qStr := `select video_id,user_id,play_url,cover_url,favored_count,comment_count,title,create_at
				from videos
				where is_delete = ?
				order by favored_count DESC
				limit ?`
	if err = db.SelectContext(ctx, &videos, qStr, 0, consts.MaxFeedVideos); err != nil {
		zap.L().Error("models video QueryVideoList method exec fail!", zap.Error(err))
	}
	return
}

// QueryVideoCommentsById 根据id查询视频评论数
func (*VideoDao) QueryVideoCommentsById(videoId int64) (comments int64, err error) {
	qStr := `select comment_count from videos where video_id = ?`
	if err = db.GetContext(ctx, &comments, qStr, videoId); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("models video comments data is null!")
			return 0, nil
		}
		zap.L().Error("models video GetContext method exec fail!", zap.Error(err))
	}
	return
}

// IsExistVideoById 判断该视频是否存在
func (*VideoDao) IsExistVideoById(videoId int64) (bool, error) {
	qStr := `select is_delete from videos where video_id = ?`
	var isDelete int
	if err := db.GetContext(ctx, &isDelete, qStr, videoId); err != nil {
		zap.L().Error("models video IsExistVideoById method exec fail!", zap.Error(err))
		return false, err
	}
	return isDelete == 0, nil
}
