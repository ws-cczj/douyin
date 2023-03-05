package middleware

import (
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SlackAuth 对前端传入的token进行选择性认证
func SlackAuth(slack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Query("token")
		if auth == "" {
			auth = c.DefaultPostForm("token", "")
		}
		// 如果松紧度比较高说明需要严格验证，否则不需要严格验证
		if slack {
			if len(auth) == 0 {
				zap.L().Error("middleware Auth token invalid!", zap.String("auth", auth))
				common.FailWithCode(c, e.FailTokenInvalid)
				c.Abort()
				return
			}
		}
		// 如果有token就验证，否则放行
		if len(auth) != 0 {
			claim, err := utils.VerifyToken(auth)
			if err != nil {
				zap.L().Error("middleware VerifyToken method fail!", zap.Error(err))
				common.FailWithCode(c, e.FailTokenVerify)
				c.Abort()
				return
			}
			c.Set("user_id", claim.UserID)
		}
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
				common.FailWithCode(c, e.FailTokenVerify)
				c.Abort()
				return
			}
			c.Set("user_id", claim.UserID)
		}
		c.Next()
	}
}
