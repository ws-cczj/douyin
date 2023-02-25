package comment

import (
	"douyin/handlers/common"
	"douyin/pkg/e"

	"github.com/gin-gonic/gin"
)

// VideoCommentHandler 评论
func VideoCommentHandler(c *gin.Context) {
	common.FailWithCode(c, e.FailServerBusy)
}
