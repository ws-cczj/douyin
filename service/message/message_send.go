package message

import (
	"douyin/consts"
	"douyin/database/models"
	"douyin/database/mongodb"
	"douyin/pkg/utils"
	"errors"
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
		return err
	}
	return nil
}

func (s *SendMessageFlow) checkNum() (err error) {
	if s.userId == 0 || s.toUserId == 0 {
		return errors.New("服务繁忙")
	}
	if s.content == "" {
		return errors.New("内容不能为空")
	}
	if len(s.content) > consts.MaxCommentLenLimit {
		return errors.New("内容不能超过500字")
	}
	if s.action != "1" {
		return errors.New("无效操作")
	}
	var isFriend bool
	if isFriend, err = models.NewRelationDao().IsExistFriend(s.userId, s.toUserId); err != nil {
		zap.L().Error("service message IsExistFriend method exec fail!", zap.Error(err))
		return
	}
	if !isFriend {
		return errors.New("你与对方还不是朋友关系")
	}
	return nil
}

func (s *SendMessageFlow) updateData() error {
	id := utils.GenID()
	if err := mongodb.NewMessageDao().InsertOneMessage(id, s.userId, s.toUserId, s.content, s.action); err != nil {
		zap.L().Error("service message InsertOneMessage method exec fail!", zap.Error(err))
		return err
	}
	return nil
}
