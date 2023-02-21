package e

// 以Fail开头统一进行错误码表示
const (
	FailUsernameExist Code = iota + 10000
	FailPasswordNotCompare
	FailParamInvalid
	FailTokenExpired
	FailTokenVerify
	FailTokenInvalid
	FailServerBusy
	FailNotKnow
)

var failMsg = map[Code]string{
	FailUsernameExist:      "username already exist!",
	FailPasswordNotCompare: "user password not compare!",
	FailParamInvalid:       "param invalid!",
	FailTokenExpired:       "user token time expired!",
	FailTokenVerify:        "user token verify fail!",
	FailTokenInvalid:       "token invalid!",
	FailServerBusy:         "server busy!",
	FailNotKnow:            "not know error!",
}

type Code int32

func (c Code) Msg() string {
	if msg, ok := failMsg[c]; ok {
		return msg
	}
	return failMsg[FailNotKnow]
}
