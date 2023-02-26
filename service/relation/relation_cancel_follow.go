package relation

import (
	"douyin/models"
	"errors"

	"go.uber.org/zap"
)

func UserCancelFollow(userId, toUserId int64) error {
	return NewUserCancelFollowFlow(userId, toUserId).Do()
}

func NewUserCancelFollowFlow(userId, toUserId int64) *UserCancelFollowFlow {
	return &UserCancelFollowFlow{userId: userId, toUserId: toUserId}
}

type UserCancelFollowFlow struct {
	isFollow, isFollower int

	userId, toUserId int64
}

func (u *UserCancelFollowFlow) Do() error {
	if err := u.checkNum(); err != nil {
		return err
	}
	if err := u.updateData(); err != nil {
		return err
	}
	return nil
}

func (u *UserCancelFollowFlow) checkNum() (err error) {
	if u.userId == 0 || u.toUserId == 0 {
		return errors.New("无效操作")
	}
	if u.isFollow, err = models.NewRelationDao().IsExistRelation(u.userId, u.toUserId); err != nil {
		zap.L().Error("service relation_follow IsExistRelation method exec fail!", zap.Error(err))
		return
	}
	if u.isFollow != 1 {
		return errors.New("无效操作")
	}
	return nil
}

func (u *UserCancelFollowFlow) updateData() (err error) {
	// 判断关系
	if u.isFollower, err = models.NewRelationDao().IsExistRelation(u.toUserId, u.userId); err != nil {
		zap.L().Error("service relation_follow IsExistRelation method exec fail!", zap.Error(err))
		return
	}
	if err = models.NewRelationDao().Action2UserRelation(u.userId, u.toUserId, u.isFollow, u.isFollower); err != nil {
		zap.L().Error("service relation_follow Action2UserRelation method exec fail!", zap.Error(err))
	}
	return nil
}
