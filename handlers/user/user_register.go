package user

import (
	"douyin/handlers/common"
	"douyin/service/user"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type RegisterResponse struct {
	common.Response
	*user.LoginResponse
}

// RegisterHandler 注册处理
func RegisterHandler(c *gin.Context) {
	// 处理参数
	uname := c.Query("username")
	pwd := c.Query("password")
	// 调用流式注册
	registerResponse, err := user.Register(uname, pwd)
	if err != nil {
		zap.L().Error("register method fail!", zap.Error(err))
		common.FailWithMsg(c, err.Error())
		return
	}
	// 返回数据
	c.JSON(http.StatusOK, RegisterResponse{
		common.Response{StatusCode: 0},
		registerResponse,
	})
}
