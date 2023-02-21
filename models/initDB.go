package models

import (
	"context"
	"douyin/conf"
	"fmt"

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
		panic(fmt.Errorf("mysql connect fail, err: %s", err))
	}
	// 配置连接池,最大空闲连接数,最大同时连接数
	db.SetMaxIdleConns(conf.Conf.MDB.MaxIdles)
	db.SetMaxOpenConns(conf.Conf.MDB.MaxOpens)
}

func Close() {
	_ = db.Close()
}