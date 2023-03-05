package relation

import (
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"douyin/service/relation"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// UserActionHandler 用户关系操作
func UserActionHandler(c *gin.Context) {
	userId := c.GetInt64("user_id")
	toUserIdStr := c.Query("to_user_id")
	action := c.Query("action_type")
	switch action {
	case "1":
		if err := relation.UserFollow(userId, utils.AtoI64(toUserIdStr)); err != nil {
			zap.L().Error("handlers relation_action UserFollow method exec fail!", zap.Error(err))
			common.FailWithMsg(c, err.Error())
			return
		}
	case "2":
		if err := relation.UserCancelFollow(userId, utils.AtoI64(toUserIdStr)); err != nil {
			zap.L().Error("handlers relation_action UserCancelFollow method exec fail!", zap.Error(err))
			common.FailWithMsg(c, err.Error())
			return
		}
	default:
		common.FailWithCode(c, e.FailServerBusy)
		return
	}
	common.SuccessWithMsg(c, "")
}
