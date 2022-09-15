package mongokv

import (
	"context"
	"errors"

	"github.com/kh9543/koala/domain/kv"
	"github.com/kh9543/koala/domain/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoKvType struct {
	db string
}

func NewMongoKv(db string) kv.Kv {
	return &MongoKvType{
		db: db,
	}
}

func (m *MongoKvType) Add(col, key string, val interface{}) error {
	switch val.(type) {
	case string:
		mongo.MongoDB.Upsert(context.Background(), m.db, col, bson.M{
			"key": key,
		}, bson.M{
			"$set": bson.M{
				"value": val,
			},
		})
	default:
		return errors.New("mongo kv val type should be string")
	}
	return nil
}

func (m *MongoKvType) Get(col, key string) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (m *MongoKvType) GetAll(col string) (map[string]interface{}, error) {
	results := []struct {
		Key   string `bson:"key" json:"key"`
		Value string `bson:"value" json:"value"`
	}{}
	if err := mongo.MongoDB.FindAll(context.Background(), m.db, col, bson.M{}, &results); err != nil {
		return nil, err
	}
	mp := make(map[string]interface{}, len(results))
	for _, r := range results {
		mp[r.Key] = r.Value
	}
	return mp, nil
}

func (m *MongoKvType) Delete(col, key string) error {
	return errors.New("not implemented")
}
