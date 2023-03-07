package e

import "errors"

// 以Fail开头统一进行错误码表示
const (
	FailUsernameExist Code = iota + 10000
	FailUsernameNotExist
	FailUsernameLimit
	FailPasswordNotCompare
	FailPasswordLimit
	FailVideoNotExist
	FailVideoIllegal
	FailVideoTitleCantNull
	FailVideoTitleLimit
	FailCommentNotExist
	FailCommentLenLimit
	FailMessageCantNULL
	FailMessageLenLimit
	FailRelationNotFriend
	FailCantFollowYourself
	FailRepeatAction
	FailTokenExpired
	FailTokenVerify
	FailTokenInvalid
	FailInitFilter
	FailInitLogger
	FailInitMongodb
	FailInitMysql
	FailInitRedis
	FailInitSnowFlake
	FailServerBusy
	FailCacheExpired
	FailNotKnow
)

var failMsg = map[Code]string{
	FailUsernameExist:      "username already exist!",
	FailUsernameNotExist:   "user not register!",
	FailUsernameLimit:      "username len overflow! (should < 32)",
	FailPasswordNotCompare: "user password not compare!",
	FailPasswordLimit:      "password len overflow! (should < 20)",
	FailVideoNotExist:      "video already delete!",
	FailVideoIllegal:       "video format is incorrect!",
	FailVideoTitleCantNull: "video title cant null!",
	FailVideoTitleLimit:    "title len overflow! (should < 200)",
	FailCommentNotExist:    "comment already delete!",
	FailCommentLenLimit:    "comment len overflow! (should < 500)",
	FailMessageCantNULL:    "message content cant null!",
	FailMessageLenLimit:    "message len overflow! (should < 2000)",
	FailRelationNotFriend:  "he/she still not you friend!",
	FailCantFollowYourself: "cant to follow yourself!",
	FailRepeatAction:       "action operator repeat!",
	FailTokenExpired:       "user token time expired!",
	FailTokenVerify:        "user token verify fail!",
	FailTokenInvalid:       "token invalid!",
	FailInitFilter:         "init filter fail!",
	FailInitLogger:         "init logger fail!",
	FailInitMongodb:        "init mongoDb fail!",
	FailInitMysql:          "init mysql fail!",
	FailInitRedis:          "init redis fail!",
	FailInitSnowFlake:      "init snowFlake fail!",
	FailServerBusy:         "server busy!",
	FailCacheExpired:       "cache expired! (ttl < 0)",
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
