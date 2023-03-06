package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Message struct {
	Id       int64  `json:"id" bson:"_id"`
	UserId   int64  `json:"from_user_id" bson:"user_id"`
	ToUserId int64  `json:"to_user_id" bson:"to_user_id"`
	CreateAt int64  `json:"create_time" bson:"create_at"`
	Content  string `json:"content" bson:"content"`
	Action   string `json:"-" bson:"action"`
}

var messageDao *MessageDao
var messageOnce sync.Once

type MessageDao struct {
}

func NewMessageDao() *MessageDao {
	messageOnce.Do(func() {
		messageDao = &MessageDao{}
	})
	return messageDao
}

// InsertOneMessage 插入一条消息
func (*MessageDao) InsertOneMessage(id, userId, toUserId int64, content, action string) error {
	_, err := message().InsertOne(ctx,
		bson.M{"_id": id,
			"from_user_id": userId,
			"to_user_id":   toUserId,
			"content":      content,
			"action":       action,
			"create_at":    time.Now().Unix()})
	return err
}

// FindMessage 查询多条消息
func (*MessageDao) FindMessage(userId, toUserId int64) (messages []*Message, err error) {
	var cursor *mongo.Cursor
	defer func() {
		if cursor != nil {
			cursor.Close(ctx)
		}
	}()
	if cursor, err = message().Find(ctx,
		bson.M{"$or": []bson.M{{"from_user_id": userId, "to_user_id": toUserId, "action": "1"},
			{"from_user_id": toUserId, "to_user_id": userId, "action": "1"}}},
		options.Find().SetSort(bson.M{"create_at": 1})); err != nil {
		zap.L().Error("mongodb message FindMessage method exec fail!", zap.Error(err))
		return
	}
	messages = []*Message{}
	if err = cursor.All(ctx, &messages); err != nil {
		zap.L().Error("mongodb message All method exec fail!", zap.Error(err))
	}
	return
}
