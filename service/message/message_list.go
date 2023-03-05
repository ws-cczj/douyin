package message

import (
	"douyin/database/mongodb"
	"douyin/pkg/e"
	"go.uber.org/zap"
)

func FriendMessage(userId, toUserId int64) ([]*mongodb.Message, error) {
	return NewFriendMessageFlow(userId, toUserId).Do()
}

func NewFriendMessageFlow(userId, toUserId int64) *FriendMessageFlow {
	return &FriendMessageFlow{userId: userId, toUserId: toUserId}
}

type FriendMessageFlow struct {
	userId, toUserId int64

	data []*mongodb.Message
}

func (f *FriendMessageFlow) Do() ([]*mongodb.Message, error) {
	if err := f.checkNum(); err != nil {
		return nil, err
	}
	if err := f.prepareData(); err != nil {
		zap.L().Error("service message_list prepareData method exec fail!", zap.Error(err))
		return nil, e.FailServerBusy.Err()
	}
	return f.data, nil
}

func (f *FriendMessageFlow) checkNum() error {
	if f.userId == 0 || f.toUserId == 0 {
		return e.FailNotKnow.Err()
	}
	return nil
}

func (f *FriendMessageFlow) prepareData() (err error) {
	f.data, err = mongodb.NewMessageDao().FindMessage(f.userId, f.toUserId)
	return
}
