package relation

import (
	models2 "douyin/database/models"
	"errors"
	"sync"

	"go.uber.org/zap"
)

func UserFollowList(userId, tkUserId int64, action bool) ([]*models2.User, error) {
	return NewUserRelationListFlow(userId, tkUserId, action).Do()
}

func NewUserRelationListFlow(userId, tkUserId int64, action bool) *UserRelationListFlow {
	return &UserRelationListFlow{userId: userId, tkUserId: tkUserId, action: action}
}

type UserRelationListFlow struct {
	userId, tkUserId int64

	action bool

	data []*models2.User
}

func (f *UserRelationListFlow) Do() ([]*models2.User, error) {
	if err := f.checkNum(); err != nil {
		return nil, err
	}
	if err := f.prepareData(); err != nil {
		return nil, err
	}
	if err := f.packData(); err != nil {
		return nil, err
	}
	return f.data, nil
}

func (f *UserRelationListFlow) checkNum() (err error) {
	if f.userId == 0 || f.tkUserId == 0 {
		return errors.New("用户无效")
	}
	return
}

func (u *UserRelationListFlow) prepareData() (err error) {
	if u.action {
		var follows int64
		if follows, err = models2.NewUserDao().QueryUserFollows(u.userId); err != nil {
			zap.L().Error("service relation_user_list QueryUserFollows method exec fail!", zap.Error(err))
			return
		}
		u.data = make([]*models2.User, follows)
		// 获取关注列表
		if err = models2.NewRelationDao().QueryUserFollowList(u.data, u.userId); err != nil {
			zap.L().Error("service relation_user_list QueryUserFollowList method exec fail!", zap.Error(err))
			return
		}
	} else {
		// 获取粉丝数
		var followers int64
		if followers, err = models2.NewUserDao().QueryUserFollowers(u.userId); err != nil {
			zap.L().Error("service relation_user_list QueryUserFollowers method exec fail!", zap.Error(err))
			return
		}
		u.data = make([]*models2.User, followers)
		// 获取粉丝列表
		if err = models2.NewRelationDao().QueryUserFollowerList(u.data, u.userId); err != nil {
			zap.L().Error("service relation_user_list QueryUserFollowerList method exec fail!", zap.Error(err))
			return
		}
	}
	return
}

func (u *UserRelationListFlow) packData() (err error) {
	var wg sync.WaitGroup
	wg.Add(len(u.data))
	for _, data := range u.data {
		user := data
		go func() {
			defer wg.Done()
			var isFollow int
			if isFollow, err = models2.NewRelationDao().IsExistRelation(u.tkUserId, user.UserId); err != nil {
				zap.L().Error("service relation_user_list IsExistRelation method exec fail!", zap.Error(err))
			}
			if isFollow == 1 {
				user.IsFollow = true
			}
		}()
	}
	return
}
