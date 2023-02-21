package consts

const (
	MaxUsernameLimit     = 32 // 最大用户名长度限制
	MaxUserPasswordLimit = 20 // 最大用户密码长度限制

	JWTTokenExpiredAt = 30 * 24 * 60 * 60 // token过期时间 30天
	JWTDouyin         = "douyin"          // 项目名称
	JWTIssuer         = "cczj"            // 签发人
	JWTSecret         = "cczj"            // 密钥

	CacheExpired          = 30 * 24 * 60 * 60     // 缓存统一过期时间 30天
	CacheMaxTryTimes      = 3                     // 缓存最大尝试次数
	CacheDouyin           = "douyin:"             // 抖音缓存
	CacheUser             = "user:"               // 用户缓存
	CacheRelation         = "relation:"           // 关系缓存
	CacheVideo            = "video:"              // 视频缓存
	CacheComment          = "comment:"            // 评论缓存
	CacheSetUserVideo     = "set_user_video:"     // set存储, 用户视频列表: key: userId, val: videoId...
	CacheSetUserFavor     = "set_user_favor:"     // set存储, 用户点赞列表: key: userId, val: videoId...
	CacheSetUserFollow    = "set_user_follow:"    // Set存储, 用户关注列表: key: userId, val: toUserId...
	CacheSetUserFollower  = "set_user_follower:"  // Set存储, 用户粉丝列表: key: userId, val: toUserId...
	CacheSetVideoComment  = "set_video_comment:"  // Set存储, 视频评论列表: key: videoId, val: commentId...
	CacheStringVideoFavor = "string_video_favor:" // String存储, 视频被点赞次数: key: videoId, val: favors
)
