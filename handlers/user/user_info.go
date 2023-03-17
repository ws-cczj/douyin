package user

import (
	"douyin/database/models"
	"douyin/handlers/common"
	"douyin/pkg/utils"
	"douyin/service/user"
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoResponse struct {
	common.Response
	*models.User `json:"user"`
}

func InfoHandler(c *gin.Context) {
	uidStr := c.Query("user_id")
	// 因为后边会检测是否为0，并且无需根据有无token进行选择服务，所以这里不需要严格判断。
	tkUid := c.GetInt64("user_id")
	uid := utils.AtoI64(uidStr)
	userResponse, err := user.Info(uid, tkUid)
	if err != nil {
		zap.L().Error("handlers user_info InfoHandler Info method exec fail", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, InfoResponse{
		common.Response{StatusCode: 0},
		userResponse,
	})
}
