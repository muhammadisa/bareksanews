package dbc

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

func OpenRedis(conf Config) (redis.Cmdable, error) {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       0,
	}), nil
}

func OpenNoSQL(conf Config) (*mongo.Database, error) {
	ctx := context.Background()
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      conf.Username,
		Password:      conf.Password,
	}
	uri := fmt.Sprintf("mongodb://%s:%s", conf.Host, conf.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(credential))
	if err != nil {
		return nil, err
	}
	mongoDb := client.Database(conf.Name)
	err = mongoDb.CreateCollection(ctx, "wallet")
	err = mongoDb.CreateCollection(ctx, "charity")
	return mongoDb, nil
}

func OpenDB(conf Config) (*sql.DB, error) {
	databaseUrl := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name,
	)
	return sql.Open("mysql", databaseUrl)
}
