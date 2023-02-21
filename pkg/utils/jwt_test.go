package utils

import (
	"reflect"
	"testing"
)

var userId int64 = 123456789

func TestGenTokenAndVerifyToken(t *testing.T) {
	token, err := GenToken(userId)
	if err != nil {
		t.Logf("gentoken fail!, err: %v", err)
		t.FailNow()
	}
	t.Logf("token: %v", token)
	claim, err := VerifyToken(token)
	if err != nil {
		t.Logf("gentoken fail!, err: %v", err)
		t.FailNow()
	}
	if ok := reflect.DeepEqual(claim.UserID, userId); !ok {
		t.Log("userID not compared!")
		t.FailNow()
	}
	t.Logf("claim res: %v", claim)
	t.Log("GenTokenAndVerifyToken test success!")
}
