package message

import (
	"douyin/database/mongodb"
	"douyin/handlers/common"
	"douyin/pkg/utils"
	"douyin/service/message"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type FriendMessageListResponse struct {
	common.Response
	Messages []*mongodb.Message `json:"message_list,omitempty"`
}

// FriendMessageListHandler 消息列表
func FriendMessageListHandler(c *gin.Context) {
	userId := c.GetInt64("user_id")
	toUserIdStr := c.Query("to_user_id")
	toUserId := utils.AtoI64(toUserIdStr)
	// TODO 因为前端问题, 无法接收到前端传来的pre_msg_time, 这里等待前端进行完善后在进行处理
	//preMsgAtStr := c.Query("pre_msg_time")
	//preMsgAt := utils.AtoI64(preMsgAtStr) / 1000
	//zap.L().Debug("pre_msg_time", zap.Int64("pre_msg_at", preMsgAt))

	var list []*mongodb.Message
	var err error
	if list, err = message.FriendMessage(userId, toUserId); err != nil {
		zap.L().Error("handlers message FriendMessage method exec fail!", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, FriendMessageListResponse{
		common.Response{StatusCode: 0},
		list,
	})
}
