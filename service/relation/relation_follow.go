package relation

import (
	"douyin/models"
	"errors"

	"go.uber.org/zap"
)

func UserFollow(userId, toUserId int64) error {
	return NewUserFollowFlow(userId, toUserId).Do()
}

func NewUserFollowFlow(userId, toUserId int64) *UserFollowFlow {
	return &UserFollowFlow{userId: userId, toUserId: toUserId}
}

type UserFollowFlow struct {
	isFollow, isFollower int

	userId, toUserId int64
}

func (u *UserFollowFlow) Do() error {
	if err := u.checkNum(); err != nil {
		return err
	}
	if err := u.updateData(); err != nil {
		return err
	}
	return nil
}

func (u *UserFollowFlow) checkNum() (err error) {
	if u.userId == u.toUserId {
		return errors.New("不能关注自己")
	}
	if u.userId == 0 || u.toUserId == 0 {
		return errors.New("无效操作")
	}
	if u.isFollow, err = models.NewRelationDao().IsExistRelation(u.userId, u.toUserId); err != nil {
		zap.L().Error("service relation_follow IsExistRelation method exec fail!", zap.Error(err))
		return
	}
	if u.isFollow == 1 {
		return errors.New("请勿重复操作")
	}
	return nil
}

func (u *UserFollowFlow) updateData() (err error) {
	// 判断关系
	if u.isFollower, err = models.NewRelationDao().IsExistRelation(u.toUserId, u.userId); err != nil {
		zap.L().Error("service relation_follow IsExistRelation method exec fail!", zap.Error(err))
		return
	}
	if u.isFollower != 1 {
		u.isFollower = 0
	}
	if err = models.NewRelationDao().Action1UserRelation(u.userId, u.toUserId, u.isFollow, u.isFollower); err != nil {
		zap.L().Error("service relation_follow Action1UserRelation method exec fail!", zap.Error(err))
	}
	return
}
