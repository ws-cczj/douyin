package models

import (
	"database/sql"
	"douyin/pkg/utils"
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
func (*RelationDao) IsExistRelation(userId, toUserId int64) (bool, error) {
	qStr := `select is_follow from user_relations where user_id = ? AND to_user_id = ?`
	var isFollow int
	if err := db.GetContext(ctx, &isFollow, qStr, userId, toUserId); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("models relation IsExistRelation result is null!")
			err = nil
		}
		return false, err
	}
	return isFollow == 1, nil
}

// QueryUserFollowList 查询用户的关注列表
func (*RelationDao) QueryUserFollowList(userId int64) (toUserIds []int64, err error) {
	qStr := `select to_user_id from user_relations where user_id = ? AND is_follow = ?`
	toUserIds = make([]int64, 0)
	if err = db.GetContext(ctx, &toUserIds, qStr, userId, 1); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("models relation QueryUserFollowList result is null!")
			err = nil
		}
	}
	zero := utils.SearchZero(toUserIds)
	return toUserIds[:zero], err
}

// QueryUserFollowerList 查询用户的粉丝列表
func (*RelationDao) QueryUserFollowerList(userId int64) (toUserIds []int64, err error) {
	qStr := `select user_id from user_relations where to_user_id = ? AND is_follow = ?`
	toUserIds = make([]int64, 0)
	if err = db.GetContext(ctx, &toUserIds, qStr, userId, 1); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("models relation QueryUserFollowList result is null!")
			err = nil
		}
	}
	zero := utils.SearchZero(toUserIds)
	return toUserIds[:zero], err
}
