package utils

import (
	"douyin/conf"
	"douyin/pkg/e"
	"fmt"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

// 每台机器相当于一个结点
var node *sf.Node

func InitSnowFlake() {
	st, _ := time.Parse("2006-01-02", conf.Conf.StartAt)
	sf.Epoch = st.UnixNano() / 1000000
	var err error
	node, err = sf.NewNode(conf.Conf.Machines)
	if err != nil {
		panic(fmt.Sprintf("%s, err: %v", e.FailInitSnowFlake.Msg(), err))
	}
}

// GenID 生成int64的分布式id
func GenID() int64 {
	return node.Generate().Int64()
}
