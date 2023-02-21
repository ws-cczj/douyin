package user

import (
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
	*user.InfoResponse
}

func InfoHandler(c *gin.Context) {
	uidStr := c.Query("user_id")
	tkUidStr, _ := c.Get("user_id")
	tkUid, ok := tkUidStr.(int64)
	if !ok {
		zap.L().Error("handlers InfoHandler token.uid invalid")
		common.FailWithCode(c, e.FailTokenInvalid)
		return
	}
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		zap.L().Error("handlers InfoHandler param uid invalid")
		common.FailWithCode(c, e.FailParamInvalid)
		return
	}
	infoResponse, err := user.Info(uid, tkUid)
	if err != nil {
		zap.L().Error("handlers InfoHandler user.Info method exec fail", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, InfoResponse{
		common.Response{StatusCode: 0},
		infoResponse,
	})
}