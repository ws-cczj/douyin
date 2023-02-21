package user

import (
	"douyin/handlers/common"
	"douyin/service/user"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// LoginHandler 登录处理
func LoginHandler(c *gin.Context) {
	// 处理参数
	uname := c.Query("username")
	pwd := c.Query("password")
	loginResponse, err := user.Login(uname, pwd)
	if err != nil {
		zap.L().Error("handlers user_Login Login method exec fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, RegisterResponse{
		common.Response{StatusCode: 0},
		loginResponse,
	})
}
