package middleware

import (
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// Auth 对前端传入的token进行认证
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Query("token")
		if auth == "" {
			auth = c.DefaultPostForm("token", "")
		}
		if len(auth) == 0 {
			zap.L().Error("middleware Auth token invalid!", zap.String("auth", auth))
			common.FailWithCode(c, e.FailTokenInvalid)
			c.Abort()
			return
		}
		//token := strings.Fields(auth)[1]
		claim, err := utils.VerifyToken(auth)
		if err != nil {
			zap.L().Error("middleware VerifyToken method fail!", zap.Error(err))
			c.Abort()
			return
		}
		c.Set("user_id", claim.UserID)
		c.Next()
	}
}

// VisitorAuth 对Feed流进行游客和用户的特殊判断
func VisitorAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Query("token")
		if auth == "" {
			auth = c.DefaultPostForm("token", "")
		}
		//token := strings.Fields(auth)[1]
		if len(auth) != 0 {
			claim, err := utils.VerifyToken(auth)
			if err != nil {
				zap.L().Error("middleware VerifyToken method fail!", zap.Error(err))
				c.Abort()
				return
			}
			c.Set("user_id", claim.UserID)
		}
		c.Next()
	}
}
