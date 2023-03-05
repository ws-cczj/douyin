package message

import (
	"douyin/consts"
	"douyin/database/models"
	"douyin/database/mongodb"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"go.uber.org/zap"
)

func SendMessage(userId, toUserId int64, action, content string) error {
	return NewSendMessageFlow(userId, toUserId, action, content).Do()
}

func NewSendMessageFlow(userId, toUserId int64, action, content string) *SendMessageFlow {
	return &SendMessageFlow{userId: userId, toUserId: toUserId, action: action, content: content}
}

type SendMessageFlow struct {
	userId, toUserId int64
	action, content  string
}

func (s *SendMessageFlow) Do() error {
	if err := s.checkNum(); err != nil {
		return err
	}
	if err := s.updateData(); err != nil {
		zap.L().Error("service message_send updateData method exec fail!", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	return nil
}

func (s *SendMessageFlow) checkNum() (err error) {
	if s.userId == 0 || s.toUserId == 0 {
		return e.FailNotKnow.Err()
	}
	if s.content == "" {
		return e.FailMessageCantNULL.Err()
	}
	if len(s.content) > consts.MaxMessageLenLimit {
		return e.FailMessageLenLimit.Err()
	}
	if s.action != "1" {
		return e.FailNotKnow.Err()
	}
	var isFriend bool
	if isFriend, err = models.NewRelationDao().IsExistFriend(s.userId, s.toUserId); err != nil {
		zap.L().Error("service message IsExistFriend method exec fail!", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	if !isFriend {
		return e.FailRelationNotFriend.Err()
	}
	return
}

func (s *SendMessageFlow) updateData() error {
	id := utils.GenID()
	if err := mongodb.NewMessageDao().InsertOneMessage(id, s.userId, s.toUserId, s.content, s.action); err != nil {
		zap.L().Error("service message InsertOneMessage method exec fail!", zap.Error(err))
		return e.FailServerBusy.Err()
	}
	return nil
}
