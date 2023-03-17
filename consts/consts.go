package consts

import "time"

const (
	CheckMaxUsername     = 32   // 最大用户名长度限制
	CheckMaxUserPassword = 20   // 最大用户密码长度限制
	CheckMaxVideoTitle   = 200  // 最大视频标题长度限制
	CheckMaxCommentLen   = 500  // 最大评论长度限制
	CheckMaxMessageLen   = 2000 // 最大消息长度限制
	CheckMaxFeedVideos   = 5    // 最大Feed取出视频条数

	JWTTokenExpiredAt = 30 * 24 * 60 * 60 // token过期时间 30天
	JWTDouyin         = "douyin"          // 项目名称
	JWTIssuer         = "cczj"            // 签发人
	JWTSecret         = "cczj"            // 密钥

	CacheExpired            = 30 * 24 * 60 * 60 * time.Second      // 缓存统一过期时间 30天, 必须加time.Second
	CacheMaxTryTimes        = 3                                    // 最大重试次数
	CacheSetUserVideo       = "douyin:user:set_user_video"         // set存储, 用户视频列表: key: userId, val: videoId...
	CacheSetUserFavor       = "douyin:favor:set_user_favor:"       // set存储, 用户点赞列表: key: userId, val: videoId...
	CacheSetUserFollow      = "douyin:relation:set_user_follow:"   // Set存储, 用户关注列表: key: userId, val: toUserId...
	CacheSetUserFollower    = "douyin:relation:set_user_follower:" // Set存储, 用户粉丝列表: key: userId, val: toUserId...
	CacheStringVideoComment = "douyin:video:string_video_comment"  // String存储, 视频评论数量: key: videoId, val: comments

	// ffmpeg的参数
	OptionInputVideoPath = "-i"
	OptionStartTime      = "-ss"
	OptionKeepTime       = "-t"
	OptionVideoFilter    = "-vf"
	OptionFormatToImage  = "-f"
	OptionAutoReWrite    = "-y"
	OptionFrames         = "-frames:v"
	DefaultVideoSuffix   = ".mp4"
	DefaultImageSuffix   = ".jpg"
)
