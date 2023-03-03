package user

import (
	"douyin/consts"
	"douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"go.uber.org/zap"
	"sync"
)

type LoginResponse struct {
	UserId int64  `json:"user_id,string"`
	Token  string `json:"token"`
}

func Login(username, password string) (*LoginResponse, error) {
	return NewUserLoginFlow(username, password).Do()
}

func NewUserLoginFlow(username, password string) *LoginFlow {
	return &LoginFlow{username: username, password: password}
}

type LoginFlow struct {
	username string
	password string

	token  string
	userId int64

	data *LoginResponse
}

func (l *LoginFlow) Do() (*LoginResponse, error) {
	if err := l.checkNum(); err != nil {
		return nil, err
	}
	if err := l.prepareData(); err != nil {
		return nil, err
	}
	if err := l.packData(); err != nil {
		return nil, err
	}
	return l.data, nil
}

func (l *LoginFlow) checkNum() (err error) {
	if l.username == "" || l.password == "" {
		return e.FailServerBusy.Err()
	}
	if len(l.username) > consts.MaxUsernameLimit {
		return e.FailUsernameLimit.Err()
	}
	if len(l.password) > consts.MaxUserPasswordLimit {
		return e.FailPasswordLimit.Err()
	}
	var wg sync.WaitGroup
	wg.Add(1)
	userDao := models.NewUserDao()
	go func() {
		defer wg.Done()
		// 1. 判断密码是否正确
		var password string
		if password, err = userDao.QueryPasswordByUsername(l.username); err != nil {
			zap.L().Error("service user_login QueryPasswordByUsername method exec fail!", zap.Error(err))
			err = e.FailServerBusy.Err()
		}
		if utils.SHA1(l.password) != password {
			zap.L().Warn("service user_login user password not compared!")
			err = e.FailPasswordNotCompare.Err()
		}
	}()
	// 2. 检查用户名是否存在
	if l.userId, err = userDao.IsExistUsername(l.username); err != nil {
		zap.L().Error("service user_login isExistUsername method exec fail!", zap.Error(err))
		err = e.FailServerBusy.Err()
	}
	if l.userId == 0 {
		zap.L().Warn("service user_login current Username not found!")
		err = e.FailUsernameNotExist.Err()
	}
	wg.Wait()
	return
}

func (l *LoginFlow) prepareData() (err error) {
	// 生成token
	if l.token, err = utils.GenToken(l.userId); err != nil {
		zap.L().Error("service user_login utils.GenToken method exec fail!", zap.Error(err))
		err = e.FailServerBusy.Err()
	}
	return
}

func (l *LoginFlow) packData() error {
	l.data = &LoginResponse{
		UserId: l.userId,
		Token:  l.token,
	}
	return nil
}
