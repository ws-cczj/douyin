package models

import (
	"douyin/consts"
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
	iStr := `insert into videos(video_id,user_id,title,play_url,cover_url) values (?,?,?,?,?)`
	if _, err = db.ExecContext(ctx, iStr, videoId, userId, title, playUrl, coverUrl); err != nil {
		zap.L().Error("models video ExecContext method exec fail!", zap.Error(err))
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
func (*VideoDao) QueryVideoListByTime(videos []*Video, lastTime int64) (err error) {
	qStr := `select video_id,user_id,play_url,cover_url,favored_count,comment_count,title,create_at
				from videos 
				where create_at >= ? AND is_delete = ?
				order by create_at ASC
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
