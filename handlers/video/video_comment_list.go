package video

import (
	"douyin/database/models"
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"douyin/service/video"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentListResponse struct {
	common.Response
	Comment []*models.Comment `json:"comment_list,omitempty"`
}

func CommentListHandler(c *gin.Context) {
	tkUserIdStr, exist := c.Get("user_id")
	videoIdStr := c.Query("video_id")
	videoId := utils.AtoI64(videoIdStr)
	var tkUserId int64
	if exist {
		var ok bool
		if tkUserId, ok = tkUserIdStr.(int64); !ok {
			zap.L().Error("handlers video_comment_list CommentListHandler tokenId format fail!", zap.Any("tokenId", tkUserIdStr))
			common.FailWithCode(c, e.FailTokenInvalid)
			return
		}
	}
	list, err := video.CommentList(videoId, tkUserId)
	if err != nil {
		zap.L().Error("handlers video_comment_list CommentList method exec fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		common.Response{
			StatusCode: 0,
		},
		list,
	})
}
