package message

import (
	"douyin/database/mongodb"
	"errors"
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
	return f.data, nil
}

func (f *FriendMessageFlow) checkNum() error {
	if f.userId == 0 || f.toUserId == 0 {
		return errors.New("服务繁忙")
	}
	return nil
}

func (f *FriendMessageFlow) prepareData() (err error) {
	f.data, err = mongodb.NewMessageDao().FindMessage(f.userId, f.toUserId)
	return
}
