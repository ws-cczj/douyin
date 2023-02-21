package video

import (
	"douyin/handlers/common"
	"douyin/service/video"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	common.Response
	*video.FeedResponse
}

func FeedHandler(c *gin.Context) {
	c.JSON(http.StatusOK, FeedResponse{
		common.Response{
			StatusCode: 0,
		},
		&video.FeedResponse{},
	})
}
