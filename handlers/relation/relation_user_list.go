package relation

import (
	"douyin/database/models"
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/pkg/utils"
	"douyin/service/relation"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserListResponse struct {
	common.Response
	User []*models.User `json:"user_list"`
}

// UserFollowListHandler 用户关注列表
func UserFollowListHandler(c *gin.Context) {
	tkUserIdStr, exist := c.Get("user_id")
	userIdStr := c.Query("user_id")
	userId := utils.AtoI64(userIdStr)
	var tkUserId int64
	if exist {
		var ok bool
		if tkUserId, ok = tkUserIdStr.(int64); !ok {
			zap.L().Error("handlers relation_user_list UserFollowListHandler tokenId format fail!", zap.Any("tokenId", tkUserIdStr))
			common.FailWithCode(c, e.FailTokenInvalid)
			return
		}
	}
	list, err := relation.UserFollowList(userId, tkUserId, true)
	if err != nil {
		zap.L().Error("handlers relation_user_list UserFollowListHandler method exec fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		common.Response{
			StatusCode: 0,
		},
		list,
	})
}

// UserFollowerListHandler 用户粉丝列表
func UserFollowerListHandler(c *gin.Context) {
	tkUserIdStr, exist := c.Get("user_id")
	userIdStr := c.Query("user_id")
	userId := utils.AtoI64(userIdStr)
	var tkUserId int64
	if exist {
		var ok bool
		if tkUserId, ok = tkUserIdStr.(int64); !ok {
			zap.L().Error("handlers relation_user_list UserFollowerListHandler tokenId format fail!", zap.Any("tokenId", tkUserIdStr))
			common.FailWithCode(c, e.FailTokenInvalid)
			return
		}
	}
	list, err := relation.UserFollowList(userId, tkUserId, false)
	if err != nil {
		zap.L().Error("handlers relation_user_list UserFollowerListHandler method exec fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		common.Response{
			StatusCode: 0,
		},
		list,
	})
}

// UserFriendListHandler 用户朋友列表
func UserFriendListHandler(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId := utils.AtoI64(userIdStr)
	list, err := relation.UserFriendList(userId)
	if err != nil {
		zap.L().Error("handlers relation_user_list UserFriendListHandler method exec fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		common.Response{
			StatusCode: 0,
		},
		list,
	})
}
