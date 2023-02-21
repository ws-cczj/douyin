package user

import (
	"douyin/cache"
	"douyin/consts"
	"douyin/models"
	"douyin/pkg/utils"
	"errors"

	"go.uber.org/zap"
)

// Register 流式处理注册操作
func Register(username, password string) (*LoginResponse, error) {
	return NewUserRegisterFlow(username, password).Do()
}

func NewUserRegisterFlow(username, password string) *RegisterFlow {
	return &RegisterFlow{username: username, password: password}
}

type RegisterFlow struct {
	username string
	password string

	token  string
	userId int64

	data *LoginResponse
}

// Do 集中处理流式操作，打包返回需要数据
func (r *RegisterFlow) Do() (*LoginResponse, error) {
	if err := r.checkNum(); err != nil {
		return nil, err
	}
	if err := r.updateData(); err != nil {
		return nil, err
	}
	if err := r.packData(); err != nil {
		return nil, err
	}
	return r.data, nil
}

// CheckNum 校验参数
func (r *RegisterFlow) checkNum() error {
	if r.username == "" {
		return errors.New("用户名为空")
	}
	if len(r.username) > consts.MaxUsernameLimit {
		return errors.New("超出用户名字数上限")
	}
	if r.password == "" {
		return errors.New("密码为空")
	}
	if len(r.password) > consts.MaxUserPasswordLimit {
		return errors.New("超出用户密码字数上限")
	}
	return nil
}

// updateData 更新数据
func (r *RegisterFlow) updateData() (err error) {
	userDao := models.NewUserDao()
	// 1. 检查用户名是否重复
	if r.userId, err = userDao.IsExistUsername(r.username); err != nil {
		zap.L().Error("service user_register isExistUsername method exec fail", zap.Error(err))
		return
	}
	if r.userId != 0 {
		zap.L().Error("service user_register current Username already exists!")
		return errors.New("用户名已经存在")
	}
	// 2. 生成雪花ID
	r.userId = utils.GenID()
	// 3. 密码加密
	r.password = utils.SHA1(r.password)
	// 4. 注册执行
	if err = userDao.AddUser(r.userId, r.username, r.password); err != nil {
		zap.L().Error("service user_register AddUser method exec fail", zap.Error(err))
		return
	}
	// 5. 生成token
	if r.token, err = utils.GenToken(r.userId); err != nil {
		zap.L().Error("service user_register utils.GenToken method exec fail!", zap.Error(err))
	}
	// 6. 初始化用户缓存数据
	go cache.NewRelationCache().SAddRegisterActionUserFollowAndFollower(r.userId)
	go cache.NewFavorCache().SAddUserFavorVideo(r.userId, -1)
	return
}

// packData 打包数据
func (r *RegisterFlow) packData() error {
	r.data = &LoginResponse{
		UserId: r.userId,
		Token:  r.token,
	}
	return nil
}
