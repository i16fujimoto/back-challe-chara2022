package db

import (
	"context"
	"os"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/joho/godotenv"
)

var MongoClient *mongo.Client


// Databaseの初期化
func InitDB() {
	
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// 空のコンテキストを作成
	ctx := context.Background()
	fmt.Println(os.Getenv("MONGO_URI")) // debug msg

	// Create a new client and connect to the server
	// 接続先の設定
	opt := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	// timezoneの設定
	rb := bson.NewRegistryBuilder()
	rb.RegisterTypeDecoder(reflect.TypeOf(time.Time{}), bsoncodec.NewTimeCodec(bsonoptions.TimeCodec().SetUseLocalTimeZone(true)))
	opt.SetRegistry(rb.Build())

	MongoClient, err = mongo.Connect(ctx, opt)
	if err != nil {
		panic(err)
	}
	// defer MongoClient.Disconnect(ctx)

	// Ping the server to DB
    err = MongoClient.Ping(ctx, readpref.Primary())
    if err != nil {
        fmt.Println("connection error:", err)
    } else {
        fmt.Println("connection success:")
    }
}