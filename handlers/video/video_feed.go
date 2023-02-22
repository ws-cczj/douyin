package video

import (
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"douyin/service/video"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	common.Response
	*video.FeedResponse
}

func FeedHandler(c *gin.Context) {
	lastTimeStr := c.Query("latest_time")
	lastTime := utils.AtoI64(lastTimeStr)
	var feedFlow *video.FeedResponse
	var err error
	if userIdStr, exist := c.Get("user_id"); exist {
		userId, ok := userIdStr.(int64)
		if !ok {
			zap.L().Error("handlers video uid invalid")
			common.FailWithCode(c, e.FailTokenInvalid)
			return
		}
		if feedFlow, err = video.UserFeed(lastTime, userId); err != nil {
			return
		}
	} else {
		if feedFlow, err = video.VisitorFeed(lastTime); err != nil {
			return
		}
	}
	c.JSON(http.StatusOK, FeedResponse{
		common.Response{
			StatusCode: 0,
		},
		feedFlow,
	})
}
