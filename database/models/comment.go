package models

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Comment struct {
	Id         int64 `json:"id" db:"id"`
	VideoId    int64 `json:"video_id" db:"video_id"`
	*User      `json:"user"`
	Content    string    `json:"content" db:"content"`
	CreateAt   string    `json:"create_date"`
	CreateTime time.Time `db:"create_at"`
}

type CommentDao struct {
}

var (
	commentDao  *CommentDao
	commentOnce sync.Once
)

// NewCommentDao 使用饿汉式单例模式初始化FavorDao对象
func NewCommentDao() *CommentDao {
	commentOnce.Do(func() {
		commentDao = new(CommentDao)
	})
	return commentDao
}

func (c *CommentDao) PublishVideoComment(userId, videoId int64, content string) (commentId int64, err error) {
	var tx *sql.Tx
	if tx, err = db.Begin(); err == nil {
		if tx == nil {
			zap.L().Error("models comment begin tx transition fail!", zap.Error(err))
			return 0, errors.New("服务繁忙")
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 增加视频的评论数量
			uStr := `update videos set comment_count = comment_count + 1 where video_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, videoId); err != nil {
				zap.L().Error("models comment incr comment fail!", zap.Error(err))
			}
		}()
		// 增加一条评论
		iStr := `insert into video_comments(user_id,video_id,content) values (?,?,?)`
		var res sql.Result
		if res, err = tx.ExecContext(ctx, iStr, userId, videoId, content); err != nil {
			zap.L().Error("models comment insert comment exec fail!", zap.Error(err))
		}
		commentId, _ = res.LastInsertId()
		wg.Wait()
	}
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return
	}
	if err = tx.Commit(); err != nil {
		zap.L().Error("models comment tx Commit exec fail!", zap.Error(err))
		tx.Rollback()
	}
	return
}

func (c *CommentDao) DeleteVideoComment(videoId, commentId int64) (err error) {
	var tx *sql.Tx
	if tx, err = db.Begin(); err == nil {
		if tx == nil {
			zap.L().Error("models comment begin tx transition fail!", zap.Error(err))
			return errors.New("服务繁忙")
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 删除一条评论
			dStr := `update video_comments set is_delete = ? where id = ?`
			if _, err = tx.ExecContext(ctx, dStr, 1, commentId); err != nil {
				zap.L().Error("models comment delete comment fail!", zap.Error(err))
			}
		}()
		// 减少视频的评论数量
		uStr := `update videos set comment_count = comment_count - 1 where video_id = ?`
		if _, err = tx.ExecContext(ctx, uStr, videoId); err != nil {
			zap.L().Error("models comment decr comment fail!", zap.Error(err))
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
		zap.L().Error("models comment tx Commit exec fail!", zap.Error(err))
		tx.Rollback()
	}
	return
}

// QueryUserCommentById 查询用户评论通过Id
func (*CommentDao) QueryUserCommentById(comment *Comment, commentId int64) (err error) {
	qStr := `select u.user_id,username,follow_count,follower_count,avatar,background_image,
				signature,total_favor_count,work_count,favor_count,t.id,t.content,t.create_at
				from users u, (select id,user_id,content,create_at
               					from video_comments
               					where id = ?) t
				where u.user_id = t.user_id`
	if err = db.GetContext(ctx, comment, qStr, commentId); err != nil {
		zap.L().Error("models comment GetContext method exec fail!", zap.Error(err))
	}
	return
}

// QueryVideoCommentsById 根据id来查询该视频下所有的评论列表ids
func (*CommentDao) QueryVideoCommentsById(videoId int64) (ids int64, err error) {
	qStr := `select COUNT(*) from video_comments where video_id = ? AND is_delete = ?`
	if err = db.GetContext(ctx, &ids, qStr, videoId, 0); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("models comment Get comments data is null!")
			return 0, nil
		}
		zap.L().Error("models comment GetContext method exec fail!", zap.Error(err))
	}
	return
}

// QueryVideoCommentListById 根据id来查询该视频下所有的评论列表
func (*CommentDao) QueryVideoCommentListById(comments []*Comment, videoId int64) (err error) {
	qStr := `select u.user_id,username,follow_count,follower_count,avatar,background_image,
       			signature,total_favor_count,work_count,favor_count,t.id,t.content,t.create_at
				from users u, (select id,user_id,content,create_at
               					from video_comments
               					where video_id = ? AND is_delete = ?) t
				where u.user_id = t.user_id
				order by t.create_at DESC`
	if err = db.SelectContext(ctx, &comments, qStr, videoId, 0); err != nil {
		zap.L().Error("models comment SelectContext method exec fail!", zap.Error(err))
	}
	return
}

// IsExistComment 判断评论是否存在
func (*CommentDao) IsExistComment(commentId int64) (bool, error) {
	qStr := `select is_delete from video_comments where id = ?`
	var isDelete int
	if err := db.GetContext(ctx, &isDelete, qStr, commentId); err != nil {
		zap.L().Error("models comment ExecContext method exec fail!", zap.Error(err))
		return false, err
	}
	return isDelete == 0, nil
}
