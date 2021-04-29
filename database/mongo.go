package database

import (
	"context"
	"fmt"
	"time"

	"gitlab.jiangxingai.com/asp-hrm/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	log "k8s.io/klog"
)

type MongoDB struct {
	Host      string
	Username  string
	Password  string
	Database  string
	TableName string
	client    *mongo.Client
}

func newMongoDB() *MongoDB {
	host := fmt.Sprintf("mongodb://%s:%s", utils.GetEnv("MONGO_ADDR", "10.56.0.52"),
		utils.GetEnv("MONGO_PORT", "27017"))
	return &MongoDB{
		Host:     host,
		Database: "",
		Username: "",
		Password: "",
		client:   nil,
	}
}

func (db *MongoDB) CloseConnection() error {
	err := db.client.Disconnect(context.TODO())
	if err != nil {
		log.Error("disconnect mongo collection error:", err)
		return err
	}
	return nil
}

func (db *MongoDB) NewConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.Host).SetMaxPoolSize(1024))
	if err != nil {
		return err
	}
	db.client = client
	return nil
}

func (db *MongoDB) CreateCollection(tableName string) (context.Context, *mongo.Collection) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	c := db.client.Database(db.Database).Collection(tableName)
	return ctx, c
}

func (db *MongoDB) UpdateDatabase(database string) {
	db.Database = database
}

func (db *MongoDB) Insert(tableName string, i interface{}) error {
	ctx, c := db.CreateCollection(tableName)
	_, err := c.InsertOne(ctx, i)
	return err
}

func (db *MongoDB) InsertMany(tableName string, i []interface{}) error {
	ctx, c := db.CreateCollection(tableName)
	_, err := c.InsertMany(ctx, i)
	return err
}

func (db *MongoDB) Delete(tableName string, i interface{}) error {
	ctx, c := db.CreateCollection(tableName)
	_, err := c.DeleteOne(ctx, i)
	return err
}

func (db *MongoDB) Update(tableName string, filter interface{}, data interface{}, upsert bool) error {
	ctx, c := db.CreateCollection(tableName)
	_, err := c.UpdateOne(ctx, filter, data, options.Update().SetUpsert(upsert))
	return err
}

func (db *MongoDB) UpdateForResult(tableName string, filter interface{}, data interface{}) (*mongo.UpdateResult, error) {
	ctx, c := db.CreateCollection(tableName)
	result, err := c.UpdateOne(ctx, filter, data)
	return result, err
}

func (db *MongoDB) ReplaceOne(tableName string, filter interface{}, data interface{}) error {
	ctx, c := db.CreateCollection(tableName)
	_, err := c.ReplaceOne(ctx, filter, data)
	return err
}

func (db *MongoDB) Commit() error {
	return nil
}

func (db *MongoDB) QueryCount(tableName string, filter interface{}) (int64, error) {
	ctx, c := db.CreateCollection(tableName)
	count, err := c.CountDocuments(ctx, filter)
	return count, err
}

func (db *MongoDB) QueryOne(tableName string, filter interface{}) *mongo.SingleResult {
	ctx, c := db.CreateCollection(tableName)
	singleResult := c.FindOne(ctx, filter)
	return singleResult
}

func (db *MongoDB) QueryOneByOptions(tableName string, filter interface{}, options *options.FindOneOptions) *mongo.SingleResult {
	ctx, c := db.CreateCollection(tableName)
	singleResult := c.FindOne(ctx, filter, options)
	return singleResult
}

func (db *MongoDB) QueryAll(tableName string, filter interface{}, options *options.FindOptions) (*mongo.Cursor, context.Context, error) {
	ctx, c := db.CreateCollection(tableName)
	cur, err := c.Find(ctx, filter, options)
	return cur, ctx, err
}

func (db *MongoDB) Query(tableName string, filter interface{}) ([]bson.M, error) {
	ctx, c := db.CreateCollection(tableName)
	cur, err := c.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var ret []bson.M
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		ret = append(ret, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (db *MongoDB) RemoveOne(tableName string, i interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	c := db.client.Database(db.Database).Collection(tableName)
	_, err := c.DeleteOne(ctx, i)
	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) RemoveMany(tableName string, i interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	c := db.client.Database(db.Database).Collection(tableName)
	_, err := c.DeleteMany(ctx, i)
	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) UpdateOne(tableName string, i interface{}, condition interface{}) error {
	ctx, c := db.CreateCollection(tableName)
	_, err := c.UpdateOne(ctx, i, condition)
	if err != nil {
		return err
	}
	return nil
}
