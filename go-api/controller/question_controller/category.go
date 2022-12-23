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

type Category struct {
	CategoryId primitive.ObjectID `json:"categoryId"`
	CategoryName string `json:"categoryName"`
}

// GET: /question/priority
// 設定可能な優先度一覧を返すAPI
func (qc QuestionController) GetCategory(c *gin.Context) {

	var err error
	var cursor *mongo.Cursor

	categoryCollection := db.MongoClient.Database("insertDB").Collection("categories")
	opts := options.Find().SetProjection(bson.M{"_id": 1, "categoryName" : 1}).SetSort(bson.D{{"categoryName", 1}})
	cursor, err = categoryCollection.Find(context.TODO(), bson.D{}, opts)
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

	var categories []Category
	for _, r := range results {
		category := Category{
			CategoryId: r["_id"].(primitive.ObjectID),
			CategoryName: r["categoryName"].(string),
		}
		categories = append(categories, category)
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
	return
}