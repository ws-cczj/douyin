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
		zap.L().Error("models relation follow data query fail!", zap.Error(err))
	}
	return
}

// QueryUserFollowIds 查询用户的关注列表ids
func (*RelationDao) QueryUserFollowIds(userId int64) (ids []int64, err error) {
	ids = []int64{}
	qStr := `select to_user_id from user_relations where user_id = ? AND is_follow = ?`
	if err = db.SelectContext(ctx, &ids, qStr, userId, 1); err != nil {
		zap.L().Error("models relation follow ids query fail!", zap.Error(err))
	}
	return
}

// QueryUserFollowerList 查询用户的粉丝列表
func (*RelationDao) QueryUserFollowerList(user []*User, toUserId int64) (err error) {
	qStr := `select user_id,username,avatar,background_image,signature,
       follow_count,follower_count,work_count,favor_count,total_favor_count
       from users where user_id in (select user_id 
                                    from user_relations
                                    where to_user_id = ? AND is_follow = ?)`
	if err = db.SelectContext(ctx, &user, qStr, toUserId, 1); err != nil {
		zap.L().Error("models relation follower data query fail!", zap.Error(err))
	}
	return
}

// QueryUserFollowerIds 查询用户的粉丝列表ids
func (*RelationDao) QueryUserFollowerIds(userId int64) (ids []int64, err error) {
	ids = []int64{}
	qStr := `select user_id from user_relations where to_user_id = ? AND is_follow = ?`
	if err = db.SelectContext(ctx, &ids, qStr, userId, 1); err != nil {
		zap.L().Error("models relation follower ids query fail!", zap.Error(err))
	}
	return
}

// QueryUserFriendsList 查询用户朋友列表
func (*RelationDao) QueryUserFriendsList(user []*User, userId int64) (err error) {
	qStr := `select user_id,username,avatar,background_image,signature,
       follow_count,follower_count,work_count,favor_count,total_favor_count
       from users where user_id in (select to_user_id 
                                    from user_relations
                                    where user_id = ? AND is_friend = 1)`
	if err = db.SelectContext(ctx, &user, qStr, userId); err != nil {
		zap.L().Error("models relation friend data query fail!", zap.Error(err))
	}
	return
}

// Action1UserRelation 动态处理用户关系
func (*RelationDao) Action1UserRelation(userId, toUserId int64, isFollow, isFollower int) (err error) {
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
				defer wg.Done()
				iStr := `insert into user_relations(user_id,to_user_id,is_friend) values (?,?,?)`
				if _, err = tx.ExecContext(ctx, iStr, userId, toUserId, isFollower); err != nil {
					zap.L().Error("models relation AddUserRelation exec fail!", zap.Error(err))
				}
			}()
		} else {
			go func() {
				defer wg.Done()
				uStr := `update user_relations set is_follow = 1, is_friend = ? where user_id = ? AND to_user_id = ?`
				if _, err = tx.ExecContext(ctx, uStr, isFollower, userId, toUserId); err != nil {
					zap.L().Error("models relation UpdateUserRelation exec fail!", zap.Error(err))
				}
			}()
		}
		// 更新用户表信息
		go func() {
			defer wg.Done()
			uStr := `update users set follow_count = follow_count + 1 where user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, userId); err != nil {
				zap.L().Error("models relation UpdateUserFollow exec fail!", zap.Error(err))
			}
		}()
		go func() {
			defer wg.Done()
			uStr := `update users set follower_count = follower_count + 1 where user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, toUserId); err != nil {
				zap.L().Error("models relation UpdateUserFollower exec fail!", zap.Error(err))
			}
		}()
		// 判断是否需要更新朋友关系
		if isFollower == 1 {
			uStr := `update user_relations 
						set is_friend = ?
                      	where user_id = ? AND to_user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, isFollower, toUserId, userId); err != nil {
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
func (*RelationDao) Action2UserRelation(userId, toUserId int64, isFollow, isFollower int) (err error) {
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
			defer wg.Done()
			uStr := `update user_relations set is_follow = 0, is_friend = 0 where user_id = ? AND to_user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, userId, toUserId); err != nil {
				zap.L().Error("models relation UpdateRelation exec fail!", zap.Error(err))
			}
		}()
		go func() {
			defer wg.Done()
			uStr := `update users set follow_count = follow_count - 1 where user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, userId); err != nil {
				zap.L().Error("models relation UpdateUserFollow exec fail!", zap.Error(err))
			}
		}()
		go func() {
			defer wg.Done()
			uStr := `update users set follower_count = follower_count - 1 where user_id = ?`
			if _, err = tx.ExecContext(ctx, uStr, toUserId); err != nil {
				zap.L().Error("models relation UpdateUserFollower exec fail!", zap.Error(err))
			}
		}()
		if isFollower == 1 {
			uStr := `update user_relations 
						set is_friend = ?
                      	where user_id = ? AND to_user_id = ?`
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

// IsExistFriend 是否存在朋友关系
func (*RelationDao) IsExistFriend(userId, toUserId int64) (bool, error) {
	qStr := `select is_friend from user_relations where user_id = ? AND to_user_id = ?`
	var isFriend int
	if err := db.GetContext(ctx, &isFriend, qStr, userId, toUserId); err != nil {
		zap.L().Error("models relation query IsFriend fail!", zap.Error(err))
	}
	return isFriend == 1, nil
}

// QueryUserFriendsById 根据id查询朋友数目
func (*RelationDao) QueryUserFriendsById(userId int64) (friends int64, err error) {
	qStr := `select Count(*) from user_relations where user_id = ? AND is_friend = 1`
	if err = db.GetContext(ctx, &friends, qStr, userId); err != nil {
		zap.L().Error("models relation query userFriends fail!", zap.Error(err))
	}
	return
}
