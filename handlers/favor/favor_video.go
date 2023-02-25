package favor

import (
	"douyin/handlers/common"
	"douyin/pkg/utils"
	"douyin/service/favor"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// VideoFavorHandler 点赞
func VideoFavorHandler(c *gin.Context) {
	videoIdStr := c.Query("video_id")
	videoId := utils.AtoI64(videoIdStr)
	action := c.Query("action_type")
	userId := c.GetInt64("user_id")

	if err := favor.FavorVideo(userId, videoId, action); err != nil {
		zap.L().Error("handlers favor FavorVideo method exec fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	common.SuccessWithMsg(c, "点赞成功!")
}
