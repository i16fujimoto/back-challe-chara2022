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


type Status struct {
	StatusId primitive.ObjectID `json:"statusId"`
	StatusName string `json:"statusName"`
}

// GET: /question/status
// 設定可能なステータス一覧を返すAPI
func (qc QuestionController) GetStatus(c *gin.Context) {
	
	var err error
	var cursor *mongo.Cursor

	statusCollection := db.MongoClient.Database("insertDB").Collection("statuses")
	opts := options.Find().SetProjection(bson.M{"_id": 1, "statusName" : 1}).SetSort(bson.D{{"statusName", 1}})
	cursor, err = statusCollection.Find(context.TODO(), bson.D{}, opts)
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

	var statuses []Status
	for _, r := range results {
		status := Status{
			StatusId: r["_id"].(primitive.ObjectID),
			StatusName: r["statusName"].(string),
		}
		statuses = append(statuses, status)
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"statuses": statuses,
	})
	return
}