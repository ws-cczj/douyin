package models

import (
	"context"
	"douyin/conf"
	"douyin/pkg/e"
	"fmt"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var ctx = context.Background()
var db *sqlx.DB

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Conf.MDB.Username,
		conf.Conf.MDB.Password,
		conf.Conf.MDB.Host,
		conf.Conf.MDB.Port,
		conf.Conf.MDB.Dbname)
	var err error
	if db, err = sqlx.Connect("mysql", dsn); err != nil {
		zap.L().Debug("mysql dsn", zap.String("dsn", dsn))
		panic(fmt.Sprintf("%s, err: %v", e.FailInitMysql.Msg(), err))
	}
	// 配置连接池,最大空闲连接数,最大同时连接数
	db.SetMaxIdleConns(conf.Conf.MDB.MaxIdles)
	db.SetMaxOpenConns(conf.Conf.MDB.MaxOpens)
}

func Close() {
	if db != nil {
		_ = db.Close()
	}
}
