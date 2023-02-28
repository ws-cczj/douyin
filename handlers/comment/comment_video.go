package comment

import (
	"douyin/database/models"
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"douyin/service/comment"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type VideoCommentRequest struct {
	VideoId   string `form:"video_id"`
	Action    string `form:"action_type"`
	Content   string `form:"comment_text"`
	CommentId string `form:"comment_id"`
}

type VideoCommentResponse struct {
	common.Response
	*models.Comment `json:"comment"`
}

// VideoCommentHandler 评论
func VideoCommentHandler(c *gin.Context) {
	userId := c.GetInt64("user_id")
	vcQ := new(VideoCommentRequest)
	var err error
	if err = c.ShouldBind(vcQ); err != nil {
		zap.L().Error("handlers comment_video param binding fail!", zap.Error(err))
		common.FailWithCode(c, e.FailParamInvalid)
		return
	}
	var mc *models.Comment
	if mc, err = comment.VideoComment(userId,
		utils.AtoI64(vcQ.VideoId),
		utils.AtoI64(vcQ.CommentId),
		vcQ.Action, vcQ.Content); err != nil {
		zap.L().Error("handlers comment_video VideoComment method exec fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	if vcQ.Action == "2" {
		common.SuccessWithMsg(c, "删除评论成功!")
		return
	}
	c.JSON(http.StatusOK, VideoCommentResponse{
		common.Response{StatusCode: 0, StatusMsg: "发布评论成功!"},
		mc,
	})
}
