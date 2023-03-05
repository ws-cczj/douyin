package consts

import "time"

const (
	MaxUsernameLimit     = 32   // 最大用户名长度限制
	MaxUserPasswordLimit = 20   // 最大用户密码长度限制
	MaxVideoTileLimit    = 200  // 最大视频标题长度限制
	MaxCommentLenLimit   = 500  // 最大评论长度限制
	MaxMessageLenLimit   = 2000 // 最大消息长度限制
	MaxFeedVideos        = 5    // 最大Feed取出视频条数

	JWTTokenExpiredAt = 30 * 24 * 60 * 60 // token过期时间 30天
	JWTDouyin         = "douyin"          // 项目名称
	JWTIssuer         = "cczj"            // 签发人
	JWTSecret         = "cczj"            // 密钥

	CacheExpired            = 30 * 24 * 60 * 60 * time.Second // 缓存统一过期时间 30天
	CacheMaxTryTimes        = 3                               // 缓存最大尝试次数
	CacheDouyin             = "douyin:"                       // 抖音缓存
	CacheUser               = "user:"                         // 用户缓存
	CacheRelation           = "relation:"                     // 关系缓存
	CacheVideo              = "video:"                        // 视频缓存
	CacheFavor              = "favor"                         // 点赞缓存
	CacheSetUserVideo       = "set_user_video:"               // set存储, 用户视频列表: key: userId, val: videoId...
	CacheSetUserFavor       = "set_user_favor:"               // set存储, 用户点赞列表: key: userId, val: videoId...
	CacheSetUserFollow      = "set_user_follow:"              // Set存储, 用户关注列表: key: userId, val: toUserId...
	CacheSetUserFollower    = "set_user_follower:"            // Set存储, 用户粉丝列表: key: userId, val: toUserId...
	CacheStringVideoComment = "string_video_comment:"         // String存储, 视频评论数量: key: videoId, val: comments
	CacheStringVideoFavor   = "string_video_favor:"           // String存储, 视频被点赞数量: key: videoId, val: favors

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
