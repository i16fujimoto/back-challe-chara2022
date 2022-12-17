package login_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/crypto"
	"back-challe-chara2022/entity/db_entity"

	"errors"
	"net/http"
	"fmt"
	"context"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
 	// "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

// POST: /signup
func CreateUser(c *gin.Context) {

	// 新規ユーザ登録

	var form body.SignUpBody
	var err error

	// バリデーション処理
	if err = c.BindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	userCollection := db.MongoClient.Database("insertDB").Collection("users")

	var doc bson.M
	// 検索条件
	filter := bson.D{{"emailAddress", form.EmailAddress}}
	// query the user collection
	err = userCollection.FindOne(context.TODO(), filter).Decode(&doc)
	if err == mongo.ErrNoDocuments {

		// パスワードのハッシュ化
		passwordEncrypt, _ := crypto.PasswordEncrypt(form.Password)

		var questionIdArray []primitive.ObjectID = make([]primitive.ObjectID, 0)
		var likeIdArray []primitive.ObjectID = make([]primitive.ObjectID, 0)
		var communityIdArray []primitive.ObjectID = make([]primitive.ObjectID, 0)

		docUser := &db_entity.User{
			UserId: primitive.NewObjectID(),
			UserName: "べあ",
			EmailAddress: form.EmailAddress,
			Password: passwordEncrypt,
			Icon: "static/icon.jpg",
			Profile: "test",
			CommunityId: communityIdArray,
			Status: "スッキリ",
			Role: db_entity.Role{RoleName: "admin", Permission: 7},
			Question: questionIdArray,
			Like: likeIdArray,
		}
		fmt.Println(*docUser)

		// Insert処理
		_, err = userCollection.InsertOne(context.TODO(), docUser) // ここでMarshalBSON()される
		if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": true, "msg": "Qmattaに登録されました"})
		return

	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}else {
		err = errors.New("同一のメールアドレスが既に登録されています。")
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "同一のメールアドレスが既に登録されています"})
		return
	}
}

// POST: /login
func LoginUser(c *gin.Context)(interface{}, error) {

	// ログイン機能

	var form body.LoginBody
	var err error

	// バリデーション処理
	if err = c.BindJSON(&form); err != nil {
		return "", err
	} else {

		// メールアドレスの照合
		userCollection := db.MongoClient.Database("insertDB").Collection("users")

		var doc bson.M
		// 検索条件
		filter := bson.D{{"emailAddress", form.EmailAddress}}
		// query the user collection
		err = userCollection.FindOne(context.TODO(), filter).Decode(&doc)
		if err != nil {
			return nil, err
		}

		// ユーザーパスワードの比較
		crypto.PasswordEncrypt(form.Password)
		err = crypto.CompareHashAndPassword(doc["password"].(string), form.Password)
		if err != nil {
			fmt.Println("パスワードが一致しませんでした。：", err)
			return nil, err
		}

		var user db_entity.User
		var docBson []byte
		// bsonにエンコード
		docBson, err = bson.Marshal(doc)
		// 構造体にデコード
		bson.Unmarshal(docBson, &user)

		return &user, nil
	}
}