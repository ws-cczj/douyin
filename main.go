package main

import (
	"douyin/conf"
	"douyin/pkg/utils"
	"douyin/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据
	InitDevs()

	r := gin.New()

	// 初始化路由
	router.InitRouter(r)

	_ = r.Run(fmt.Sprintf("%s:%d", conf.Conf.Ip, conf.Conf.Port))
}

func InitDevs() {
	// 初始化雪花生成器
	utils.InitSnowFlake()
	// 初始化数据库
	// 初始化redis缓存
	// TODO 初始化敏感词拦截器。
}
