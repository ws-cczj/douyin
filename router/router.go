package router

import (
	"douyin/conf"
	"douyin/handlers/user"
	"douyin/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	if conf.Conf.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	v1 := r.Group("/douyin")
	// 全局使用自定义中间件 logger日志 recovery 异常恢复
	v1.Use(middleware.GinLogger(), middleware.GinRecovery())
	{
		// basic apis
		v1.POST("/user/login", user.LoginHandler)
	}
}
