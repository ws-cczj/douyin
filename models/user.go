package models

import (
	"database/sql"
	"errors"
	"sync"

	"go.uber.org/zap"
)

type User struct {
	UserId      int64  `json:"user_id,string" db:"user_id"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password,omitempty" db:"password"`
	Avatar      string `json:"avatar" db:"avatar"`
	Description string `json:"signature" db:"description"`
	BgImage     string `json:"background_image" db:"bg_image"`
}

type UserDao struct {
}

var (
	userDao  *UserDao
	userOnce sync.Once
)

// NewUserDao 使用饿汉式单例模式初始化UserDao对象
func NewUserDao() *UserDao {
	userOnce.Do(func() {
		userDao = new(UserDao)
	})
	return userDao
}

// AddUser 注册一个新用户
func (*UserDao) AddUser(user_id int64, username, password string) (err error) {
	iStr := `insert into users(user_id,username,password) values(?,?,?)`
	_, err = db.ExecContext(ctx, iStr, user_id, username, password)
	if err != nil {
		zap.L().Error("models AddUser method exec fail", zap.Error(err))
	}
	return
}

// IsExistUsername 查询是否存在该用户名
func (*UserDao) IsExistUsername(username string) (user_id int64, err error) {
	qStr := `select user_id from users where username = ?`
	if err = db.GetContext(ctx, &user_id, qStr, username); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		zap.L().Error("models GetContext method exec fail!", zap.Error(err))
	}
	return
}

// QueryPasswordByUsername 通过用户名查找密码
func (*UserDao) QueryPasswordByUsername(username string) (password string, err error) {
	qStr := `select password from users where username = ?`
	if err = db.GetContext(ctx, &password, qStr, username); err != nil {
		zap.L().Error("models QueryPasswordByUsername method exec fail!", zap.Error(err))
	}
	return
}

// QueryUserInfoById 根据Id查询用户信息
func (*UserDao) QueryUserInfoById(user *User, userId int64) (err error) {
	if user == nil {
		return errors.New("null pointer error")
	}
	qStr := `select user_id,username,avatar,bg_image,description from users where user_id = ?`
	if err = db.GetContext(ctx, user, qStr, userId); err != nil {
		zap.L().Error("models QueryUserInfoById method exec fail!", zap.Error(err))
	}
	return
}
