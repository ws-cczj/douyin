package e

import "errors"

// 以Fail开头统一进行错误码表示
const (
	FailUsernameExist Code = iota + 10000
	FailUsernameNotExist
	FailUsernameLimit
	FailPasswordNotCompare
	FailPasswordLimit
	FailParamInvalid
	FailTokenExpired
	FailTokenVerify
	FailTokenInvalid
	FailVideoIllegal
	FailServerBusy
	FailCacheExpried
	FailNotKnow
)

var failMsg = map[Code]string{
	FailUsernameExist:      "username already exist!",
	FailUsernameNotExist:   "user not register!",
	FailUsernameLimit:      "username len overflow! (should < 32)",
	FailPasswordNotCompare: "user password not compare!",
	FailPasswordLimit:      "password len overflow! (should < 20)",
	FailParamInvalid:       "param invalid!",
	FailTokenExpired:       "user token time expired!",
	FailTokenVerify:        "user token verify fail!",
	FailTokenInvalid:       "token invalid!",
	FailVideoIllegal:       "video format is incorrect!",
	FailServerBusy:         "server busy!",
	FailCacheExpried:       "cache ttl < 0",
	FailNotKnow:            "not know error!",
}

type Code int32

// Err 根据code创建对应错误
func (c Code) Err() error {
	if msg, ok := failMsg[c]; ok {
		return errors.New(msg)
	}
	return errors.New(failMsg[FailNotKnow])
}

// Msg 根据code返回对应错误消息
func (c Code) Msg() string {
	if msg, ok := failMsg[c]; ok {
		return msg
	}
	return failMsg[FailNotKnow]
}
