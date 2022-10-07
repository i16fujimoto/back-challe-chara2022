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

		var question_id_array []primitive.ObjectID = make([]primitive.ObjectID, 0)
		var like_id_array []primitive.ObjectID = make([]primitive.ObjectID, 0)
		var community_id_array []primitive.ObjectID = make([]primitive.ObjectID, 0)

		// 仮データ
		toneId, _ := primitive.ObjectIDFromHex("633ee9f501830d402ce385c3")
		bearId, _ := primitive.ObjectIDFromHex("633f7fc114dae75d9e701e24")

		docUser := &db_entity.User{
			UserId: primitive.NewObjectID(),
			UserName: "べあ",
			EmailAddress: form.EmailAddress,
			Password: passwordEncrypt,
			Icon: "img_dir/icon.png",
			Profile: "test",
			CommunityId: community_id_array,
			Status: "スッキリ",
			Role: db_entity.Role{RoleName: "admin", Permission: 7},
			BearIcon: bearId,
			BearTone: toneId,
			Question: question_id_array,
			Like: like_id_array,
		}
		fmt.Println(*docUser)

		// Insert処理
		_, err = userCollection.InsertOne(context.TODO(), docUser) // ここでMarshalBSON()される
		if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		}
		c.JSON(200, gin.H{"result": true})
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
func LoginUser(c *gin.Context) {

	// ログイン機能

	var form body.LoginBody
	var err error

	// バリデーション処理
	if err = c.BindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else {

		// メールアドレスの照合
		userCollection := db.MongoClient.Database("insertDB").Collection("users")

		var doc bson.M
		// 検索条件
		filter := bson.D{{"emailAddress", form.EmailAddress}}
		// query the user collection
		err = userCollection.FindOne(context.TODO(), filter).Decode(&doc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		} else if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the user_id")
			c.JSON(http.StatusOK, gin.H{"result": false, "user": nil})
			return
		}

		// ユーザーパスワードの比較
		passwordEncrypt, _ := crypto.PasswordEncrypt(form.Password)
		fmt.Println(doc["password"].(string), passwordEncrypt) // debug msg
		err = crypto.CompareHashAndPassword(doc["password"].(string), form.Password)
		if err != nil {
			fmt.Println("パスワードが一致しませんでした。：", err)
			c.JSON(http.StatusOK, gin.H{"result": false, "user": nil})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": true, "user": doc})
		return
	}
}