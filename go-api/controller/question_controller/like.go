package question_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"

	"net/http"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
)

type Like struct {
	UserName string `json:"userName"`
	Icon []byte `json:"icon"`
}


// PATCH: /question/answer/like
// 質問・回答のいいね
func(qc QuestionController) PatchLike(c *gin.Context) {
	
	var err error

	// JWTよりuserIdの取得
	claims := jwt.ExtractClaims(c)
	userId, _ := primitive.ObjectIDFromHex(claims["userId"].(string))

	// Bodyの内容を取得
	var request body.PatchLikeBody
	if err = c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// コレクションと接続
	questionCollection := db.MongoClient.Database("insertDB").Collection("questions")

	if id := request.QuestionId; id != "" {
		// 質問にいいねした場合（いいねをはずした場合）
		questionId, _ := primitive.ObjectIDFromHex(id)
		// いいね　が存在するかチェック
		var result bson.M
		checkFilter := bson.M{"_id": questionId, "like": userId}
		isErr := questionCollection.FindOne(context.TODO(), checkFilter).Decode(&result)
		
		update := bson.M{"like": userId}
		filter := bson.M{"_id": questionId}

		// いいね　を更新
		if isErr != nil {
			// いいね　を追加
			if _, err := questionCollection.UpdateOne(context.TODO(), filter, bson.M{"$push": update}); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusBadRequest,
					"message": err.Error(),
				})
				return
			}
		} else {
			// いいね　を削除
			if _, err := questionCollection.UpdateOne(context.TODO(), filter, bson.M{"$pull": update}); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusBadRequest,
					"message": err.Error(),
				})
				return
			}
		}
	} else if id := request.AnswerId; id != "" {
		// 回答にいいねした場合（いいねをはずした場合）
		answerId, _ := primitive.ObjectIDFromHex(id)
		// いいね　が存在するかチェック
		var result bson.M
		checkFilter := bson.M{"answer._id": answerId, "answer.like": userId}
		isErr := questionCollection.FindOne(context.TODO(), checkFilter).Decode(&result)
		
		update := bson.M{"answer.$[element].like": userId}
		filter := bson.M{"answer._id": answerId}
		opts := options.Update().SetArrayFilters(options.ArrayFilters{ // 配列にフィルターをかけれる
			Filters: []interface{}{
				 bson.M{
					"element._id": answerId,
				},
			},
		})

		// いいね　を更新
		if isErr != nil {
			// いいね　を追加
			if _, err := questionCollection.UpdateOne(context.TODO(), filter, bson.M{"$push": update}, opts); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusBadRequest,
					"message": err.Error(),
				})
				return
			}
		} else {
			// いいね　を削除
			if _, err := questionCollection.UpdateOne(context.TODO(), filter, bson.M{"$pull": update}, opts); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusBadRequest,
					"message": err.Error(),
				})
				return
			}
		}
	} else {
		// いいねをしていない場合
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "validation error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "success",
	})
	return

}