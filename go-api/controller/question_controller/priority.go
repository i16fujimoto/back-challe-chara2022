package question_controller

import (
	"back-challe-chara2022/db"

	"net/http"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

type Priority struct {
	PriorityId primitive.ObjectID `json:"priorityId"`
	PriorityName string `json:"priorityName"`
}

// GET: /question/priority
// 設定可能な優先度一覧を返すAPI
func (qc QuestionController) GetPriority(c *gin.Context) {

	var err error
	var cursor *mongo.Cursor

	priorityCollection := db.MongoClient.Database("insertDB").Collection("priorities")
	opts := options.Find().SetProjection(bson.M{"_id": 1, "priorityName" : 1}).SetSort(bson.D{{"priorityName", 1}})
	cursor, err = priorityCollection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": err.Error(),
		})
		return
	}
	// 検索結果をresultsにデコード
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	var priorities []Priority
	for _, r := range results {
		priority := Priority{
			PriorityId: r["_id"].(primitive.ObjectID),
			PriorityName: r["priorityName"].(string),
		}
		priorities = append(priorities, priority)
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"priorities": priorities,
	})
	return
}