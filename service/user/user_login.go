package user

import (
	"douyin/consts"
	"douyin/models"
	"douyin/pkg/utils"
	"errors"

	"go.uber.org/zap"
)

type LoginResponse struct {
	UserId int64  `json:"user_id"`
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
	if l.username == "" {
		return errors.New("用户名为空")
	}
	if l.password == "" {
		return errors.New("密码为空")
	}
	if len(l.username) > consts.MaxUsernameLimit {
		return errors.New("超出用户名字数上限")
	}
	if len(l.password) > consts.MaxUserPasswordLimit {
		return errors.New("超出用户密码字数上限")
	}
	// 1. 检查用户名是否存在
	userDao := models.NewUserDao()
	if l.userId, err = userDao.IsExistUsername(l.username); err != nil {
		zap.L().Error("service user_login isExistUsername method exec fail!", zap.Error(err))
		return
	}
	if l.userId == 0 {
		zap.L().Warn("service user_login current Username not found!")
		return errors.New("该用户还未注册")
	}
	// 2. 判断密码是否正确
	var password string
	if password, err = userDao.QueryPasswordByUsername(l.username); err != nil {
		zap.L().Error("service user_login QueryPasswordByUsername method exec fail!", zap.Error(err))
		return
	}
	if utils.SHA1(l.password) != password {
		zap.L().Warn("service user_login password not compared!")
		return errors.New("用户密码不匹配")
	}
	return nil
}

func (l *LoginFlow) prepareData() (err error) {
	if l.token, err = utils.GenToken(l.userId); err != nil {
		zap.L().Error("service user_login utils.GenToken method exec fail!", zap.Error(err))
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
