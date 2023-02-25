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
	lastTime := utils.AtoI64(lastTimeStr) / 1000
	var feedResponse *video.FeedResponse
	var err error
	if userIdStr, exist := c.Get("user_id"); exist {
		userId, ok := userIdStr.(int64)
		if !ok {
			zap.L().Error("handlers video_feed uid invalid")
			common.FailWithCode(c, e.FailTokenInvalid)
			return
		}
		if feedResponse, err = video.UserFeed(lastTime, userId); err != nil {
			zap.L().Error("handlers video_feed UserFeed method exec fail!", zap.Error(err))
			common.FailWithMsg(c, err.Error())
			return
		}
	} else {
		if feedResponse, err = video.VisitorFeed(lastTime); err != nil {
			zap.L().Error("handlers video_feed VisitorFeed method exec fail!", zap.Error(err))
			common.FailWithMsg(c, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, FeedResponse{
		common.Response{
			StatusCode: 0,
		},
		feedResponse,
	})
}
