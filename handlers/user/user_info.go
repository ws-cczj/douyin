package user

import (
	"douyin/database/models"
	"douyin/handlers/common"
	"douyin/pkg/e"
	"douyin/service/user"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type InfoResponse struct {
	common.Response
	*models.User `json:"user"`
}

func InfoHandler(c *gin.Context) {
	uidStr := c.Query("user_id")
	tkUidStr, _ := c.Get("user_id")
	tkUid, ok := tkUidStr.(int64)
	if !ok {
		zap.L().Error("handlers user_info InfoHandler uid invalid")
		common.FailWithCode(c, e.FailTokenVerify)
		return
	}
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		zap.L().Error("handlers user_info InfoHandler param uid invalid")
		common.FailWithCode(c, e.FailNotKnow)
		return
	}
	var userResponse *models.User
	if userResponse, err = user.Info(uid, tkUid); err != nil {
		zap.L().Error("handlers user_info InfoHandler Info method exec fail", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, InfoResponse{
		common.Response{StatusCode: 0},
		userResponse,
	})
}
