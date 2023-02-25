package user

import (
	"database/sql"
	"douyin/models"
	"errors"

	"go.uber.org/zap"
)

func Info(userId, tkUserId int64) (*models.User, error) {
	return NewUserInfoFlow(userId, tkUserId).Do()
}

func NewUserInfoFlow(userId, tkUserId int64) *InfoFlow {
	return &InfoFlow{userId: userId, tkUserId: tkUserId, data: new(models.User)}
}

type InfoFlow struct {
	userId   int64
	tkUserId int64

	data *models.User
}

func (i *InfoFlow) Do() (*models.User, error) {
	if err := i.checkNum(); err != nil {
		return nil, err
	}
	if err := i.prepareData(); err != nil {
		return nil, err
	}
	return i.data, nil
}

func (i *InfoFlow) checkNum() (err error) {
	if i.userId == 0 || i.tkUserId == 0 {
		return errors.New("用户信息错误")
	}
	// 根据User_id查询数据库获取User信息。
	if err = models.NewUserDao().QueryUserInfoById(i.data, i.userId); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Error("service user_info checkNum UserId not exist!", zap.Error(err))
			return errors.New("无法查询到该用户")
		}
		zap.L().Error("service user_info QueryUserInfoById method exec fail!", zap.Error(err))
	}
	return
}

func (i *InfoFlow) prepareData() (err error) {
	// 判断用户关系
	if i.tkUserId != i.userId {
		var isFollow int
		if isFollow, err = models.NewRelationDao().IsExistRelation(i.tkUserId, i.userId); err != nil {
			zap.L().Error("service user_info IsExistRelation method exec fail!", zap.Error(err))
		}
		if isFollow == 1 {
			i.data.IsFollow = true
		}
	}
	return
}
