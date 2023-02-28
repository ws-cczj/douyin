package router

import (
	"douyin/conf"
	"douyin/handlers/comment"
	"douyin/handlers/favor"
	"douyin/handlers/message"
	"douyin/handlers/relation"
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
		// 游客判断
		{
			v1.GET("/feed/", middleware.VisitorAuth(), video.FeedHandler)
			v1.GET("/publish/list/", middleware.VisitorAuth(), user.PublishVideoListHandler)
			v1.GET("/favorite/list/", middleware.VisitorAuth(), user.FavorVideoListHandler)
			v1.GET("/comment/list/", middleware.VisitorAuth(), video.CommentListHandler)
			v1.GET("/relation/follow/list/", middleware.VisitorAuth(), relation.UserFollowListHandler)
			v1.GET("/relation/follower/list/", middleware.VisitorAuth(), relation.UserFollowerListHandler)
		}
		// 加入JWT用户认证
		v1.Use(middleware.Auth())
		{
			v1.GET("/user/", user.InfoHandler)
			v1.GET("/relation/friend/list/", relation.UserFriendListHandler)
			v1.POST("/publish/action/", middleware.Ffmpeg(true), video.PublishHandler)
			v1.POST("/favorite/action/", favor.VideoFavorHandler)
			v1.POST("/comment/action/", comment.VideoCommentHandler)
			v1.POST("/relation/action/", relation.UserActionHandler)
			v1.POST("/message/action/", message.SendMessageHandler)
			v1.GET("/message/chat/", message.FriendMessageListHandler)
		}
	}
}
