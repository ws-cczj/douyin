package models

import (
	"sync"

	"go.uber.org/zap"
)

type Video struct {
	Id           int64  `json:"id"`
	VideoId      int64  `json:"video_id"`
	UserId       int64  `json:"user_id"`
	FavorCount   int64  `json:"favorite_count"`
	CommentCount int64  `json:"comment_count"`
	Title        string `json:"title"`
	PlayUrl      string `json:"play_url"`
	CoverUrl     string `json:"cover_url"`
	IsFavor      bool   `json:"is_favor"`
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

// PublishVideo TODO 发布一条视频
func (*VideoDao) PublishVideo() {}

// QueryUserVideoList 查询用户发布的视频列表
func (*VideoDao) QueryUserVideoList(userId int64) (videoIds []int64, err error) {
	qStr := `select video_id from videos where user_id = ? AND is_delete = ?`
	videoIds = make([]int64, 4)
	if err = db.GetContext(ctx, &videoIds, qStr, userId, 0); err != nil {
		zap.L().Error("models video QueryUserVideoList method exec fail!", zap.Error(err))
	}
	return
}
