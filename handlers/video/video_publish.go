package video

import (
	"douyin/handlers/common"
	"douyin/service/video"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// PublishHandler 发布视频
func PublishHandler(c *gin.Context) {
	userId := c.GetInt64("user_id")
	playUrl := c.GetString("play_url")
	coverUrl := c.GetString("cover_url")
	title := c.PostForm("title")
	zap.String("title", title)

	if err := video.Publish(userId, playUrl, coverUrl, title); err != nil {
		zap.L().Error("handlers video_publish Publish method exec fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	common.SuccessWithMsg(c, "发布成功!")
}
