package user

import (
	"douyin/database/models"
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"douyin/service/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type VideoListResponse struct {
	common.Response
	Videos []*models.Video `json:"video_list,omitempty"`
}

// PublishVideoListHandler 用户发布视频列表
func PublishVideoListHandler(c *gin.Context) {
	tkUserIdStr, exist := c.Get("user_id")
	userIdStr := c.Query("user_id")
	userId := utils.AtoI64(userIdStr)
	// 由于这里有游客登陆和用户登录两种情况，因此要在外侧进行严格的token id验证
	var tkUserId int64
	if exist {
		var ok bool
		if tkUserId, ok = tkUserIdStr.(int64); !ok {
			zap.L().Error("handlers user_info_list PublishVideoListHandler tokenId format fail!", zap.Any("tokenId", tkUserIdStr))
			common.FailWithCode(c, e.FailTokenVerify)
			return
		}
	}
	list, err := user.PublishVideoList(userId, tkUserId)
	if err != nil {
		zap.L().Error("handlers user_info_list PublishVideoListHandler PublishVideoList method exec fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		common.Response{StatusCode: 0},
		list,
	})
}

// FavorVideoListHandler 用户点赞视频列表
func FavorVideoListHandler(c *gin.Context) {
	tkUserIdStr, exist := c.Get("user_id")
	userIdStr := c.Query("user_id")
	userId := utils.AtoI64(userIdStr)
	var tkUserId int64
	if exist {
		var ok bool
		if tkUserId, ok = tkUserIdStr.(int64); !ok {
			zap.L().Error("handlers user_video_list FavorVideoListHandler tokenId format fail!", zap.Any("tokenId", tkUserIdStr))
			common.FailWithCode(c, e.FailTokenInvalid)
			return
		}
	}
	list, err := user.FavorVideoList(userId, tkUserId)
	if err != nil {
		zap.L().Error("handlers user_video_lsit FavorVideoListHandler VisitorVideoList method exec fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		common.Response{
			StatusCode: 0,
		},
		list,
	})
}
