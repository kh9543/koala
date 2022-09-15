package mongo

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo interface {
	FindOne(ctx context.Context, db, table string, filter interface{}, result interface{}) error
}

type Mongo struct {
	client *mongo.Client
}

func (m *Mongo) Upsert(ctx context.Context, db, table string, filter interface{}, update interface{}) error {
	opts := options.Update().SetUpsert(true)
	if _, err := m.client.Database(db).Collection(table).UpdateOne(context.TODO(), filter, update, opts); err != nil {
		return err
	}
	return nil
}

func (m *Mongo) FindOne(ctx context.Context, db, table string, filter interface{}, result interface{}) error {
	if err := m.client.Database(db).Collection(table).FindOne(ctx, filter).Decode(result); err != nil {
		return err
	}
	return nil
}

func (m *Mongo) FindAll(ctx context.Context, db, table string, filter interface{}, result interface{}) error {
	cursor, err := m.client.Database(db).Collection(table).Find(ctx, filter)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, result)
	if err != nil {
		return err
	}
	return nil
}

var MongoDB Mongo

func init() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("should set MONGODB_URI")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	MongoDB = Mongo{
		client: client,
	}
}
