package conf

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Server    `mapstructure:"server"`
	SnowFlake `mapstructure:"snowflake"`
	Logger    `mapstructure:"logger"`
	MDB       Mysql `mapstructure:"mysql"`
	RDB       Redis `mapstructure:"redis"`
}

type Server struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
	Ip   string `mapstructure:"ip"`
}

type SnowFlake struct {
	StartAt  string `mapstructure:"start_at"`
	Machines int64  `mapstructure:"machines"`
}

type Logger struct {
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Layout     string `mapstructure:"layout"`
	Filename   string `mapstructure:"filename"`
}

type Mysql struct {
	MaxIdles int    `mapstructure:"max_idles_conns"`
	MaxOpens int    `mapstructure:"max_opens_conns"`
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	Dbname   string `mapstructure:"dbname"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Redis struct {
	UserDB     int    `mapstructure:"user_db"`
	RelationDB int    `mapstructure:"relation_db"`
	VideoDB    int    `mapstructure:"video_db"`
	FavorDB    int    `mapstructure:"favor_db"`
	CommentDB  int    `mapstructure:"comment_db"`
	PoolSize   int    `mapstructure:"pool_size"`
	Addr       string `mapstructure:"addr"`
	Password   string `mapstructure:"password"`
}

func init() {
	Conf = new(Config)
	var err error
	// 统一处理初始化错误,避免多 error 场景出现
	defer func() {
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}()

	viper.SetConfigFile("./conf/config.yaml")
	err = viper.ReadInConfig()
	// 将配置信息反序列化到 Conf 全局变量中去
	err = viper.Unmarshal(Conf)

	//fmt.Println(Conf)
	// 监视文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 一旦发生变化就进行重新赋值
		if err = viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	})
}
