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
	Id       int64     `json:"id" bson:"_id"`
	UserId   int64     `json:"from_user_id" bson:"user_id"`
	ToUserId int64     `json:"to_user_id" bson:"to_user_id"`
	Content  string    `json:"content" bson:"content"`
	Action   string    `json:"-" bson:"action"`
	CreateAt time.Time `json:"create_time" bson:"create_at"`
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
			"create_at":    time.Now().Format("2006-01-02 15:04:05")})
	return err
}

// FindMessage 查询多条消息
func (*MessageDao) FindMessage(userId, toUserId int64) (messages []*Message, err error) {
	var cursor *mongo.Cursor
	defer cursor.Close(ctx)
	if cursor, err = message().Find(ctx,
		bson.M{"$or": []bson.M{{"user_id": userId, "to_user_id": toUserId},
			{"user_id": toUserId, "to_user_id": userId}}},
		options.Find().SetSort(bson.M{"create_at": 1})); err != nil {
		zap.L().Error("mongodb message FindMessage method exec fail!", zap.Error(err))
		return nil, err
	}
	messages = make([]*Message, 0)
	for cursor.Next(ctx) {
		var msg *Message
		if err = cursor.Decode(&msg); err != nil {
			zap.L().Warn("mongodb message Decode fail!", zap.Error(err))
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
