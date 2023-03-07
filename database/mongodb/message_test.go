package mongodb

import (
	"douyin/pkg/utils"
	"testing"
)

func TestMongoInit(t *testing.T) {
	InitMongodb()
	utils.InitSnowFlake()
	TestInsertMessage(t)
	TestFindMessage(t)
}

func TestInsertMessage(t *testing.T) {
	if err := NewMessageDao().InsertOneMessage(utils.GenID(), 8037037554798592, 8579925648871424, "新朋友你好!", "1"); err != nil {
		t.Errorf("insertMessage test fail!, err:%v", err)
		t.FailNow()
	}
}

func TestFindMessage(t *testing.T) {
	data, err := NewMessageDao().FindMessage(8037037554798592, 8579925648871424)
	if err != nil {
		t.Errorf("FindMessage test fail!, err:%v", err)
		t.FailNow()
	}
	t.Log(data)
}
