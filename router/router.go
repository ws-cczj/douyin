package router

import (
	"douyin/conf"
	"douyin/handlers/comment"
	"douyin/handlers/favor"
	"douyin/handlers/user"
	"douyin/handlers/video"
	"douyin/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	if conf.Conf.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Static("/static", "./public")
	v1 := r.Group("/douyin")
	// 全局使用自定义中间件 logger日志 recovery 异常恢复
	v1.Use(middleware.GinLogger(), middleware.GinRecovery())
	{
		// basic apis
		v1.POST("/user/register/", user.RegisterHandler)
		v1.POST("/user/login/", user.LoginHandler)
		v1.GET("/feed/", middleware.VisitorAuth(), video.FeedHandler)
		v1.GET("/publish/list/", middleware.VisitorAuth(), user.PublishVideoListHandler)
		v1.GET("/favorite/list/", middleware.VisitorAuth(), user.FavorVideoListHandler)
		// 加入JWT认证中间件
		v1.Use(middleware.Auth())
		{
			v1.GET("/user/", user.InfoHandler)
			v1.POST("/publish/action/", middleware.Ffmpeg(true), video.PublishHandler)
			v1.POST("/favorite/action/", favor.VideoFavorHandler)
			v1.POST("/comment/action/", comment.VideoCommentHandler)
		}
	}
}
