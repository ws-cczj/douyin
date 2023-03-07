package relation

import (
	models "douyin/database/models"
	"douyin/pkg/e"
	"go.uber.org/zap"
)

// UserFriendList 用户朋友列表
func UserFriendList(userId int64) ([]*models.User, error) {
	return NewUserFriendListFlow(userId).Do()
}

func NewUserFriendListFlow(userId int64) *UserFriendListFlow {
	return &UserFriendListFlow{userId: userId}
}

type UserFriendListFlow struct {
	userId int64

	data []*models.User
}

func (u *UserFriendListFlow) Do() ([]*models.User, error) {
	if err := u.checkNum(); err != nil {
		return nil, err
	}
	if err := u.prepareData(); err != nil {
		zap.L().Error("service relation_friend_list prepareData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	return u.data, nil
}

func (u *UserFriendListFlow) checkNum() error {
	if u.userId == 0 {
		return e.FailNotKnow.Err()
	}
	return nil
}

func (u *UserFriendListFlow) prepareData() (err error) {
	// 获取朋友数目
	var friends int64
	relationDao := models.NewRelationDao()
	if friends, err = relationDao.QueryUserFriendsById(u.userId); err != nil {
		zap.L().Error("service relation_friend_list method exec fail!", zap.Error(err))
	}
	u.data = make([]*models.User, friends)
	// 查询朋友信息
	if err = relationDao.QueryUserFriendsList(u.data, u.userId); err != nil {
		zap.L().Error("service relation_friend_list method exec fail!", zap.Error(err))
	}
	return
}
