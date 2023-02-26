package models

import (
	"database/sql"
	"errors"
	"sync"

	"go.uber.org/zap"
)

type Relation struct {
	IsFollow int
	Id       int64
	UserId   int64
	ToUserId int64
}

type RelationDao struct {
}

var (
	relationDao  *RelationDao
	relationOnce sync.Once
)

// NewRelationDao 使用饿汉式单例模式初始化UserDao对象
func NewRelationDao() *RelationDao {
	relationOnce.Do(func() {
		relationDao = new(RelationDao)
	})
	return relationDao
}

// IsExistRelation 判断是否存在关系
func (*RelationDao) IsExistRelation(userId, toUserId int64) (isFollow int, err error) {
	qStr := `select is_follow from user_relations where user_id = ? AND to_user_id = ?`
	if err = db.GetContext(ctx, &isFollow, qStr, userId, toUserId); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("models relation IsExistRelation result is null!")
			err = nil
		}
		return -1, err
	}
	return isFollow, nil
}

// QueryUserFollowList 查询用户的关注列表
func (*RelationDao) QueryUserFollowList(user []*User, userId int64) (err error) {
	qStr := `select user_id,username,avatar,background_image,signature,
       follow_count,follower_count,work_count,favor_count,total_favor_count
       from users where user_id in (select to_user_id 
                                    from user_relations
                                    where user_id = ? AND is_follow = ?)`
	if err = db.SelectContext(ctx, &user, qStr, userId, 1); err != nil {
		zap.L().Error("models relation SelectContext method exec fail!", zap.Error(err))
	}
	return
}

// QueryUserFollowerList 查询用户的关注列表
func (*RelationDao) QueryUserFollowerList(user []*User, toUserId int64) (err error) {
	qStr := `select user_id,username,avatar,background_image,signature,
       follow_count,follower_count,work_count,favor_count,total_favor_count
       from users where user_id in (select user_id 
                                    from user_relations
                                    where to_user_id = ? AND is_follow = ?)`
	if err = db.SelectContext(ctx, &user, qStr, toUserId, 1); err != nil {
		zap.L().Error("models relation SelectContext method exec fail!", zap.Error(err))
	}
	return
}

// Action1UserRelation 动态处理用户关系
func (*RelationDao) Action1UserRelation(userId, toUserId int64, isFollow, isFriend int) (err error) {
	var tx *sql.Tx
	if tx, err = db.Begin(); err == nil {
		if tx == nil {
			zap.L().Error("models relation begin tx transition fail!", zap.Error(err))
			return errors.New("服务繁忙")
		}
		// 添加用户关系
		var wg sync.WaitGroup
		wg.Add(3)
		// 要么该关系不存在，要么该关系为未成立
		if isFollow == -1 {
			go func() {
				iStr := `insert into user_relations(user_id,to_user_id,is_friend) values (?,?,?)`
				if _, err = tx.ExecContext(ctx, iStr, userId, toUserId, isFriend); err != nil {
					zap.L().Error("models relation AddUserRelation exec fail!", zap.Error(err))
				}
				wg.Done()
			}()
		} else {
			go func() {
				uStr := `update user_relations set is_follow = ? where user_id = ?`
				if _, err = tx.ExecContext(ctx, uStr, 1, userId); err != nil {
					zap.L().Error("models relation UpdateUserRelation exec fail!", zap.Error(err))
				}
				wg.Done()
			}()
		}
		// 更新用户表信息
		go func() {
			uStr := `update users set follow_count = follow_count + 1 where user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, userId); err != nil {
				zap.L().Error("models relation UpdateUserFollow exec fail!", zap.Error(err))
			}
			wg.Done()
		}()
		go func() {
			uStr := `update users set follower_count = follower_count + 1 where user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, toUserId); err != nil {
				zap.L().Error("models relation UpdateUserFollower exec fail!", zap.Error(err))
			}
			wg.Done()
		}()
		// 判断是否需要更新朋友关系
		if isFriend == 1 {
			uStr := `update user_relations 
						set is_friend = ?
                      	where to_user_id = ? AND user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, isFriend, toUserId, userId); err != nil {
				zap.L().Error("models relation AddUserRelation exec fail!", zap.Error(err))
			}
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
		zap.L().Error("models relation tx Commit exec fail!", zap.Error(err))
		tx.Rollback()
	}
	return
}

// Action2UserRelation 移除用户关系
func (*RelationDao) Action2UserRelation(userId, toUserId int64, isFollow, isFriend int) (err error) {
	var tx *sql.Tx
	if tx, err = db.Begin(); err == nil {
		if tx == nil {
			zap.L().Error("models relation begin tx transition fail!", zap.Error(err))
			return errors.New("服务繁忙")
		}
		// 添加用户关系
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			uStr := `update user_relations set is_follow = 0 AND is_friend = 0 where user_id = ? AND to_user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, userId, toUserId); err != nil {
				zap.L().Error("models relation UpdateRelation exec fail!", zap.Error(err))
			}
			wg.Done()
		}()
		go func() {
			uStr := `update users set follow_count = follow_count - 1 where user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, userId); err != nil {
				zap.L().Error("models relation UpdateUserFollow exec fail!", zap.Error(err))
			}
			wg.Done()
		}()
		go func() {
			uStr := `update users set follower_count = follower_count - 1 where user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, toUserId); err != nil {
				zap.L().Error("models relation UpdateUserFollower exec fail!", zap.Error(err))
			}
			wg.Done()
		}()
		if isFriend == 1 {
			uStr := `update user_relations 
						set is_friend = ?
                      	where to_user_id = ? AND user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, 0, toUserId, userId); err != nil {
				zap.L().Error("models relation SubUserRelation exec fail!", zap.Error(err))
			}
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
		zap.L().Error("models relation tx Commit exec fail!", zap.Error(err))
		tx.Rollback()
	}
	return
}
