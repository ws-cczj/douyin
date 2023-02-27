package models

import (
	"database/sql"
	"errors"
	"sync"

	"go.uber.org/zap"
)

type User struct {
	UserId          int64  `json:"id,string" db:"user_id"`
	FollowCount     int64  `json:"follow_count" db:"follow_count"`
	FollowerCount   int64  `json:"follower_count" db:"follower_count"`
	WorkCount       int64  `json:"work_count" db:"work_count"`
	FavorCount      int64  `json:"favorite_count" db:"favor_count"`
	TotalFavorCount int64  `json:"total_favorited,string" db:"total_favor_count"`
	IsFollow        bool   `json:"is_follow"`
	Username        string `json:"name" db:"username"`
	Password        string `json:"password,omitempty" db:"password"`
	Avatar          string `json:"avatar" db:"avatar"`
	Signature       string `json:"signature" db:"signature"`
	BackgroundImage string `json:"background_image" db:"background_image"`
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
	qStr := `select user_id,username,avatar,background_image,signature,
       follow_count,follower_count,work_count,favor_count,total_favor_count
       from users where user_id = ?`
	if err = db.GetContext(ctx, user, qStr, userId); err != nil {
		zap.L().Error("models QueryUserInfoById method exec fail!", zap.Error(err))
	}
	return
}

// QueryUserFollows 获取用户关注数
func (*UserDao) QueryUserFollows(userId int64) (follows int64, err error) {
	qStr := `select follow_count from users where user_id = ?`
	if err = db.GetContext(ctx, &follows, qStr, userId); err != nil {
		zap.L().Error("models user GetContext method exec fail!", zap.Error(err))
	}
	return
}

// QueryUserFollowers 获取用户粉丝数
func (*UserDao) QueryUserFollowers(userId int64) (followers int64, err error) {
	qStr := `select follower_count from users where user_id = ?`
	if err = db.GetContext(ctx, &followers, qStr, userId); err != nil {
		zap.L().Error("models user GetContext method exec fail!", zap.Error(err))
	}
	return
}

// QueryUserFavorVideos 查询用户的点赞视频数量
func (*UserDao) QueryUserFavorVideos(userId int64) (favors int64, err error) {
	qStr := `select favor_count from users where user_id = ?`
	if err = db.GetContext(ctx, &favors, qStr, userId); err != nil {
		zap.L().Error("models user GetContext method exec fail!", zap.Error(err))
	}
	return
}
