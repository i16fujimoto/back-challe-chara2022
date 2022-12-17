package question_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	// "back-challe-chara2022/entity/db_entity"

	"net/http"
	"fmt"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	// jwt "github.com/appleboy/gin-jwt/v2"
)

type QuestionController struct {}

type Priority struct {
	PriorityId primitive.ObjectID `json:"priorityId"`
	PriorityName string `json:"priorityName"`
}

type Status struct {
	StatusId primitive.ObjectID `json:"statusId"`
	StatusName string `json:"statusName"`
}

type Question struct {
	QuestionId primitive.ObjectID `json:"questionId"`
	Title string `json:"title"`
	Category []string `json:"category"`
	Questioner primitive.ObjectID `json:"questioner"`
	NumLikes int `json:"numLikes"`
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

// GET: /question
// 質問一覧を返すAPI
func (qc QuestionController) GetQuestions(c *gin.Context) {
	
	var err error

	var request body.GetQuestionsBody
	if err = c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	var cursor *mongo.Cursor
	questionCollection := db.MongoClient.Database("insertDB").Collection("questions")
	opts := options.Find().SetProjection(bson.M{"_id": 1, "title": 1, "category": 1, "questioner": 1, "like": 1, "createdAt": 1}).
		SetSort(bson.D{{"createdAt", -1}})
	communityId, _ := primitive.ObjectIDFromHex(request.CommunityId)
	fmt.Println(communityId)
	cursor, err = questionCollection.Find(context.TODO(), bson.M{"communityId": communityId}, opts)
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

	var questions []Question
	for _, r := range results {

		// カテゴリーを格納
		var category []string
		for _, c := range r["category"].(primitive.A) {
			category = append(category, c.(string))
		}

		// Likeの数
		var cntLikes int = len(r["like"].(primitive.A))

		question := Question{
			QuestionId: r["_id"].(primitive.ObjectID),
			Title: r["title"].(string),
			Category: category,
			Questioner: r["questioner"].(primitive.ObjectID),
			NumLikes: cntLikes,
		}
		questions = append(questions, question)
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"questions": questions,
	})
	return
}

