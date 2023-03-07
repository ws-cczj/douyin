package mongodb

import (
	"context"
	"douyin/conf"
	"douyin/pkg/e"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var ctx = context.Background()
var mdb *mongo.Client

func InitMongodb() {
	var err error
	if mdb, err = mongo.NewClient(options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s", conf.Conf.MgDB.Addr)).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(conf.Conf.MgDB.MaxPoolSize).
		SetMinPoolSize(conf.Conf.MgDB.MinPoolSize).
		SetMaxConnecting(conf.Conf.MgDB.MaxConns)); err != nil {
		panic(fmt.Sprintf("%s, err: %v", e.FailInitMongodb.Msg(), err))
	}
	// connect
	if err = mdb.Connect(ctx); err != nil {
		panic(fmt.Sprintf("%s, err: %v", e.FailInitMongodb.Msg(), err))
	}
}

var clt *mongo.Collection
var cltOnce sync.Once

// message 返回Message集合
func message() *mongo.Collection {
	cltOnce.Do(func() {
		if mdb != nil {
			clt = mdb.Database(conf.Conf.MgDB.Db).Collection(conf.Conf.MgDB.Colt)
		}
	})
	return clt
}

// Close 关闭mongodb数据库连接
func Close() {
	if mdb != nil {
		_ = mdb.Disconnect(ctx)
	}
}
