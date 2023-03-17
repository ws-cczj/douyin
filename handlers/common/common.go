package common

import (
	"douyin/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode e.Code `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// FailWithCode 带有Code的请求失败响应
func FailWithCode(c *gin.Context, code e.Code) {
	c.JSON(http.StatusOK, Response{StatusCode: code, StatusMsg: code.Msg()})
}

// FailWithMsg 带有Msg的请求失败响应
func FailWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{StatusCode: e.FailServerBusy, StatusMsg: msg})
}

// FailInner 数据库等内部程序错误避免返回给用户，进行统一错误处理
//func FailInner(c *gin.Context) {
//	c.JSON(http.StatusOK, Response{
//		StatusCode: e.FailServerBusy,
//		StatusMsg:  e.FailServerBusy.Msg(),
//	})
//}

// SuccessWithMsg 带有消息的请求成功响应
func SuccessWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: msg})
}
