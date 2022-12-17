package server

import (
	"back-challe-chara2022/controller/login_controller"
	"back-challe-chara2022/entity/db_entity"
	"back-challe-chara2022/db"
	
	"time"
	"context"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetJWTAuthentication(key string)(interface{}, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm: "Qumatta User SECRET AREA", // Realm name to display to the user. Required.
		SigningAlgorithm: "HS256",  // signing algorithm - possible values are HS256, HS384, HS512
		Key: []byte(key), // Secret key used for signing. Required.
		Timeout: 2 * time.Hour, // Duration that a jwt token is valid.
		// このフィールドは、MaxRefreshが経過するまでクライアントがトークンをリフレッシュできるようする
		// トークンの最大有効期限は MaxRefresh + Timeout となる
		MaxRefresh:  time.Hour, 
		IdentityKey: "userId",

		//「ログイン時」に呼び出されるコールバック関数
		// ログイン情報をもとにユーザの認証を行う
		// ユーザーデータをユーザー識別子として返す必要があり、それは Claim Array に格納される
		Authenticator: func(c *gin.Context)(interface{}, error) {
			if object, err := login_controller.LoginUser(c); err != nil {
				if object == "" {
					return "", jwt.ErrMissingLoginValues // Bodyのvalidationが失敗
				} else {
					return nil, jwt.ErrFailedAuthentication // ログイン失敗
				}
			} else {
				return object, nil
			}
		},

		// 「ログイン時」に呼び出されるコールバック関数
		// この関数を使用するとウェブトークンに追加のペイロードデータを追加することが可能（ペイロードのクレーム設定）
		// このデータはリクエスト時に c.Get("JWT_PAYLOAD") を介して利用可能
		// ペイロードは暗号化されていないことに注意
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*db_entity.User); ok {
				return jwt.MapClaims{
					"userId": v.UserId,	
					"password": v.Password,
				}
			}

			return jwt.MapClaims{}
		}, 

		// 「トークン認証時」に呼び出されるコールバック関数
		// クレームからログインIDを取得する
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			userId, _ := primitive.ObjectIDFromHex(claims["userId"].(string))
			return &db_entity.User{
				UserId: userId,
			}
		},

		//「トークン認証時」に呼び出されるコールバック関数
		// トークンのユーザ情報からの認証
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*db_entity.User); ok && verifyUserId(v.UserId) {
				return true
			}
			return false
		},

		// Case: Unauthorized
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code": code,
				"message": message,
			})
		},

		// TokenLookup は "<source>:<name>" という形式の文字列で、リクエストからトークンを抽出するために利用
    	// リクエストからトークンを抽出するために利用
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		// TimeFunc は現在の時刻を指定する
		// テストやサーバーがトークンと異なるタイムゾーンを使用している場合にも適用
		TimeFunc: time.Now,
	})
}

// userIDがDB中に存在するか
func verifyUserId(userId primitive.ObjectID) bool {
	var err error
	var doc bson.M
	// 検索条件
	filter := bson.D{{"_id", userId}}
	// query the user collection
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	err = userCollection.FindOne(context.TODO(), filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return false
	}else if err != nil {
		return false
	}
	return true
}