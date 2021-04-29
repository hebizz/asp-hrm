package database

import (
	"context"

	"gitlab.jiangxingai.com/asp-hrm/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	log "k8s.io/klog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db MongoDatabase

//初始化数据库
func Setup() {
	dbClient := NewDatabase()
	err := dbClient.NewConnection()
	if err != nil {
		log.Error(err)
	}
	dbClient.UpdateDatabase("asp-hrm")
	Db = dbClient.(MongoDatabase)
	err = Db.Update("department", bson.M{"_id": "1"},
		bson.M{"$set": &interfaces.Department{Id: "1", Title: "未分组", Type: "department"}},
		true)
	if err != nil {
		log.Error(err)
	}
	err = Db.Update("countHuman", bson.M{"_id": "1"},
		bson.M{"$setOnInsert": bson.M{"count": 0}}, true)
	if err != nil {
		log.Error(err)
	}
}

type Database interface {
	NewConnection() error
	UpdateDatabase(string)
	Commit() error
	Insert(string, interface{}) error
	Query(string, interface{}) ([]bson.M, error)
}

type MongoDatabase interface {
	Database
	CloseConnection() error
	InsertMany(string, []interface{}) error
	Update(string, interface{}, interface{}, bool) error
	UpdateForResult(string, interface{}, interface{}) (*mongo.UpdateResult, error)
	ReplaceOne(string, interface{}, interface{}) error
	Delete(string, interface{}) error
	QueryCount(string, interface{}) (int64, error)
	QueryOne(string, interface{}) *mongo.SingleResult
	QueryOneByOptions(string, interface{}, *options.FindOneOptions) *mongo.SingleResult
	QueryAll(string, interface{}, *options.FindOptions) (*mongo.Cursor, context.Context, error)
	RemoveOne(string, interface{}) error
	RemoveMany(string, interface{}) error
	UpdateOne(string, interface{}, interface{}) error
}

func NewDatabase() Database {
	return newMongoDB()
}
