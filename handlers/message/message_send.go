package message

import (
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"douyin/service/message"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SendMessageRequest struct {
	ToUserId string `form:"to_user_id"`
	Action   string `form:"action_type"`
	Content  string `form:"content"`
}

// SendMessageHandler 发送消息
func SendMessageHandler(c *gin.Context) {
	msg := new(SendMessageRequest)
	if err := c.ShouldBind(msg); err != nil {
		zap.L().Error("handlers message param invalid!", zap.Error(err))
		common.FailWithCode(c, e.FailParamInvalid)
		return
	}
	userId := c.GetInt64("user_id")
	if err := message.SendMessage(userId, utils.AtoI64(msg.ToUserId), msg.Action, msg.Content); err != nil {
		zap.L().Error("handlers message SendMessage method exec fail!", zap.Error(err))
		return
	}
	common.SuccessWithMsg(c, "")
}
