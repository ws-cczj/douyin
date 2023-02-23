package main

import (
	"context"
	"douyin/conf"
	"douyin/models"
	"douyin/pkg/logger"
	"douyin/pkg/utils"
	"douyin/router"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	InitDevs()

	r := gin.New()

	// 初始化路由
	router.InitRouter(r)

	goAndShutdown(r)
}

// InitDevs 初始化数据
func InitDevs() {
	// 初始化日志
	logger.InitLogger()
	// 初始化雪花生成器
	utils.InitSnowFlake()
	// 初始化数据库
	models.InitMysql()
	// 初始化redis缓存
	//cache.InitRedis()
	// TODO 初始化敏感词拦截器。

}

// goAndShutdown 启动程序和优雅关机
func goAndShutdown(r *gin.Engine) {
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Conf.Ip, conf.Conf.Port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen is fail!", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// -- 创建一个超过5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown fail!", zap.Error(err))
	}
	models.Close()
	//cache.Close()
}
