package user

import (
	"douyin/database/models"
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"douyin/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VideoListResponse struct {
	common.Response
	Videos []*models.Video `json:"video_list,omitempty"`
}

// PublishVideoListHandler 用户视频列表
func PublishVideoListHandler(c *gin.Context) {
	tkUserIdStr, exist := c.Get("user_id")
	userIdStr := c.Query("user_id")
	userId := utils.AtoI64(userIdStr)
	var tkUserId int64
	if exist {
		var ok bool
		if tkUserId, ok = tkUserIdStr.(int64); !ok {
			zap.L().Error("handlers user_info_list PublishVideoListHandler tokenId format fail!", zap.Any("tokenId", tkUserIdStr))
			common.FailWithCode(c, e.FailTokenInvalid)
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
		common.Response{
			StatusCode: 0,
		},
		list,
	})
}

// FavorVideoListHandler 用户视频列表
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
