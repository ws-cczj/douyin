package user

import (
	"douyin/consts"
	"douyin/database/models"
	"douyin/pkg/e"
	"douyin/pkg/utils"
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
		zap.L().Error("service user_register updateData method exec fail!", zap.Error(err))
		return nil, err
	}
	if err := r.packData(); err != nil {
		zap.L().Error("service user_register packData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	return r.data, nil
}

// checkNum 校验参数
func (r *RegisterFlow) checkNum() (err error) {
	// 检查用户名或者密码是否为空
	if r.username == "" || r.password == "" {
		return e.FailNotKnow.Err()
	}
	// 检查用户名长度是否超过限度
	if len(r.username) > consts.CheckMaxUsername {
		return e.FailUsernameLimit.Err()
	}
	// 检查用户密码长度是否超过限度
	if len(r.password) > consts.CheckMaxUserPassword {
		return e.FailPasswordLimit.Err()
	}
	// 检查用户名是否重复
	if r.userId, err = models.NewUserDao().IsExistUsername(r.username); err != nil {
		zap.L().Error("service user_register isExistUsername method exec fail", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	// 检查该用户是否存在
	if r.userId != 0 {
		zap.L().Warn("service user_register current Username already exists!")
		return e.FailUsernameExist.Err()
	}
	return
}

// updateData 更新数据
func (r *RegisterFlow) updateData() (err error) {
	// 1. 生成雪花ID
	r.userId = utils.GenID()
	// 2. 密码加密
	r.password = utils.SHA1(r.password)
	// 3. 注册执行
	if err = models.NewUserDao().AddUser(r.userId, r.username, r.password); err != nil {
		zap.L().Error("service user_register AddUser method exec fail", zap.Error(err))
		return
	}
	// 4. 生成token
	if r.token, err = utils.GenToken(r.userId); err != nil {
		zap.L().Error("service user_register utils.GenToken method exec fail!", zap.Error(err))
	}
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
