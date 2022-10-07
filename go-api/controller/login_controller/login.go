package login_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/crypto"
	"back-challe-chara2022/entity/db_entity"

	"net/http"
	"fmt"
	"context"
	
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
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
	} else {

		// パスワードのハッシュ化
		passwordEncrypt, _ := crypto.PasswordEncrypt(form.Password)

		userCollection := db.MongoClient.Database("insertDB").Collection("users")
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
		c.JSON(200, gin.H{"result": "success create user"})
		return
	}
}