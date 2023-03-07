package models

import (
	"database/sql"
	"douyin/pkg/e"
	"sync"

	"go.uber.org/zap"
)

type Favor struct {
	Id      int64 `json:"id"`
	UserId  int64 `json:"user_id"`
	VideoId int64 `json:"video_id"`
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

// AddUserFavorVideoInfoById 添加用户点赞视频操作
func (*FavorDao) AddUserFavorVideoInfoById(userId, videoId int64, isFavor int) (err error) {
	var tx *sql.Tx
	if tx, err = db.Begin(); err == nil {
		if tx == nil {
			zap.L().Error("models favor begin tx transition fail!", zap.Error(err))
			return e.FailServerBusy.Err()
		}
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			defer wg.Done()
			// 将用户点赞进行数+1
			uStr := `update users set favor_count = favor_count + 1 where user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, userId); err != nil {
				zap.L().Error("models favor AddFavorCount exec fail!", zap.Error(err))
			}
		}()
		go func() {
			// 将用户被点赞数进行+1
			defer wg.Done()
			uStr := `update users
				set total_favor_count = total_favor_count + 1
				where user_id = (
					select user_id 
					from videos 
					where video_id = ?)`
			if _, err = tx.ExecContext(ctx, uStr, videoId); err != nil {
				zap.L().Error("models favor AddTotalFavors exec fail!", zap.Error(err))
			}
		}()
		go func() {
			// 将视频被点赞数进行+1
			defer wg.Done()
			uStr := `update videos set favored_count = favored_count + 1 where video_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, videoId); err != nil {
				zap.L().Error("models favor AddVideoFavored To Table fail!", zap.Error(err))
			}
		}()
		switch isFavor {
		case -1:
			iStr := `insert into user_favor_videos(user_id, video_id) values(?,?)`
			if _, err = tx.ExecContext(ctx, iStr, userId, videoId); err != nil {
				zap.L().Error("models favor AddFavorData To table fail!", zap.Error(err))
			}
		case 0:
			uStr := `update user_favor_videos set is_favor = ? where user_id = ? AND video_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, 1, userId, videoId); err != nil {
				zap.L().Error("models favor UpdateFavorData To table fail!", zap.Error(err))
			}
		default:
			err = e.FailNotKnow.Err()
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
		zap.L().Error("models favor tx Commit exec fail!", zap.Error(err))
		tx.Rollback()
	}
	return
}

// SubUserFavorsInfoById 通过用户id对用户取消点赞进行操作
func (*FavorDao) SubUserFavorsInfoById(userId, videoId int64, isFavor int) (err error) {
	var tx *sql.Tx
	if tx, err = db.Begin(); err == nil {
		if tx == nil {
			zap.L().Error("models favor begin tx transition fail!", zap.Error(err))
			return e.FailServerBusy.Err()
		}
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			// 将用户点赞进行-1
			defer wg.Done()
			uStr := `update users set favor_count = favor_count - 1 where user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, userId); err != nil {
				zap.L().Error("models favor SubFavorCount exec fail!", zap.Error(err))
			}
		}()
		go func() {
			// 将用户被点赞进行-1
			defer wg.Done()
			uStr := `update users
				set total_favor_count = total_favor_count - 1
				where user_id = (
					select user_id 
					from videos 
					where video_id = ?)`
			if _, err = tx.ExecContext(ctx, uStr, videoId); err != nil {
				zap.L().Error("models favor SubTotalFavorCount exec fail!", zap.Error(err))
			}
		}()
		go func() {
			// 将视频被点赞数进行-1
			defer wg.Done()
			uStr := `update videos set favored_count = favored_count - 1 where video_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, videoId); err != nil {
				zap.L().Error("models favor AddVideoFavored To Table fail!", zap.Error(err))
			}
		}()
		if isFavor != 1 {
			err = e.FailNotKnow.Err()
		}
		uStr := `update user_favor_videos set is_favor = ? where user_id = ? AND video_id = ?`
		if _, err = tx.ExecContext(ctx, uStr, 0, userId, videoId); err != nil {
			zap.L().Error("models favor UpdateFavorData To table fail!", zap.Error(err))
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
		zap.L().Error("models favor tx Commit exec fail!", zap.Error(err))
		// 防止占用资源，未提交成功一定要回滚
		tx.Rollback()
	}
	return
}

// QueryUserFavorVideoList 查询用户点赞的视频列表
func (*FavorDao) QueryUserFavorVideoList(userId int64) (favorsVideos []int64, err error) {
	qStr := `select video_id 
					from user_favor_videos 
					where user_id = ? AND is_favor = ?`
	favorsVideos = []int64{}
	if err = db.SelectContext(ctx, &favorsVideos, qStr, userId, 1); err != nil {
		zap.L().Error("models favor SelectContext method exec fail!", zap.Error(err))
	}
	return
}

// IsExistFavor 是否存在点赞
func (*FavorDao) IsExistFavor(userId, videoId int64) (isFavor int, err error) {
	qStr := `select is_favor from user_favor_videos where user_id = ? AND video_id = ?`
	if err = db.GetContext(ctx, &isFavor, qStr, userId, videoId); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("models favor IsExistFavor result is null!")
			err = nil
		}
		return -1, err
	}
	return isFavor, nil
}
