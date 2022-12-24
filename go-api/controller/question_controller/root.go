package question_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/entity/db_entity"
	"back-challe-chara2022/s3"

	"net/http"
	"fmt"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
)

type QuestionController struct {}


type Question struct {
	QuestionId primitive.ObjectID `json:"questionId"`
	Title string `json:"title"`
	Category []string `json:"category"`
	Status string `json:"status"`
	Priority string `json:"priority"`
	Questioner string `json:"questioner"`
	QuestionerIcon []byte `json:"questionerIcon"`
	NumLikes int `json:"numLikes"`
	CreatedAt primitive.DateTime `json:"createdAt"`
}


// GET: /question:<ObjectId: communityId>
// 質問一覧を返すAPI
func (qc QuestionController) GetQuestions(c *gin.Context) {
	
	var err error

	var cursor *mongo.Cursor
	questionCollection := db.MongoClient.Database("insertDB").Collection("questions")
	opts := options.Find().SetProjection(bson.M{"_id": 1, "title": 1, "category": 1, "priority": 1, "status": 1, "questioner": 1, "like": 1, "createdAt": 1}).
		SetSort(bson.D{{"createdAt", -1}})
	communityId, _ := primitive.ObjectIDFromHex(c.Param("communityId"))
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
		var categories []string
		for _, category := range r["category"].(primitive.A) {
			categoryId := category.(primitive.ObjectID)
			filterStatus := bson.M{"_id": categoryId}
			var docCategory bson.M
			categoryCollection := db.MongoClient.Database("insertDB").Collection("categories")
			if err := categoryCollection.FindOne(context.TODO(), filterStatus,
				options.FindOne().SetProjection(bson.M{"_id": 0, "categoryName": 1})).Decode(&docCategory); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusBadRequest,
					"message": err.Error(),
				})
				return
			}
			categories = append(categories, docCategory["categoryName"].(string))
		}

		// 質問の投稿日時を格納
		var timestamp primitive.DateTime = r["createdAt"].(primitive.DateTime)

		// Likeの数
		var cntLikes int = len(r["like"].(primitive.A))

		// ステータスを取得
		statusId := r["status"].(primitive.ObjectID)
		filterStatus := bson.M{"_id": statusId}
		var docStatus bson.M
		statusCollection := db.MongoClient.Database("insertDB").Collection("statuses")
		if err := statusCollection.FindOne(context.TODO(), filterStatus,
			options.FindOne().SetProjection(bson.M{"_id": 0, "statusName": 1})).Decode(&docStatus); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		// 優先度を取得
		priorityId := r["priority"].(primitive.ObjectID)
		filterPriority := bson.M{"_id": priorityId}
		var docPriority bson.M
		priorityCollection := db.MongoClient.Database("insertDB").Collection("priorities")
		if err := priorityCollection.FindOne(context.TODO(), filterPriority,
			options.FindOne().SetProjection(bson.M{"_id": 0, "priorityName": 1})).Decode(&docPriority); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		userCollection := db.MongoClient.Database("insertDB").Collection("users")
		filterUser := bson.M{"_id": r["questioner"].(primitive.ObjectID)}
		var docUser bson.M
		if err := userCollection.FindOne(context.TODO(), filterUser,
			options.FindOne().SetProjection(bson.M{"_id": 0, "userName": 1, "icon": 1})).Decode(&docUser); err != nil {
			fmt.Println("aaaaaa")
				c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		// S3バケットとオブジェクトを指定
		url := docUser["icon"].(string)
		var bucketIndex int = strings.Index(url, "/") // 最初に "/" が出現する位置
		var bucketName, key string = url[:bucketIndex], url[bucketIndex:]
		
		fmt.Println(bucketName, key)
		
		// S3インスタンスを作成
		s3Instance, err := s3.NewS3()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code": 503,
				"message": "Service Unavailable",
			})
		}

		// S3から画像ファイルのダウンロード
		downloadKey := s3.GetObjectInput(bucketName, key)
		buf, err := s3.Download(s3Instance, downloadKey) //[]byte
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code": 404,
				"message": err.Error(),
			})
		}

		question := Question{
			QuestionId: r["_id"].(primitive.ObjectID),
			Title: r["title"].(string),
			Category: categories,
			Priority: docPriority["priorityName"].(string),
			Status: docStatus["statusName"].(string),
			Questioner: docUser["userName"].(string),
			QuestionerIcon: buf,
			NumLikes: cntLikes,
			CreatedAt: timestamp,
		}


		questions = append(questions, question)
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"questions": questions,
	})
	return
}

// POST: /question/<ObjectId: communityId>
// 質問を登録するAPI
func (qc QuestionController) PostQuestion(c *gin.Context) {
	var err error

	// JWTよりuserIdの取得
	claims := jwt.ExtractClaims(c)
	userId, _ := primitive.ObjectIDFromHex(claims["userId"].(string))
	fmt.Println(userId) // debug message
	
	// Bodyを受け取る
	var request body.PostQuestionBody
	if err = c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	questionId := primitive.NewObjectID() // 質問ID

	// データの変換
	communityId, _ := primitive.ObjectIDFromHex(c.Param("communityId"))
	priorityId, _ := primitive.ObjectIDFromHex(request.Priority)
	statusId, _ := primitive.ObjectIDFromHex(request.Status)
	var categoriesId []primitive.ObjectID
	for _, categoryId := range request.Category {
		Id, _ := primitive.ObjectIDFromHex(categoryId)
		categoriesId = append(categoriesId, Id)
	}

	// 画像のアップロード
	var urls []string
	for _, obj := range request.Image {
		urls = append(urls, obj)
	}

	// 質問を登録
	questionCollection := db.MongoClient.Database("insertDB").Collection("questions")
	// 登録データ
	docQuestion := &db_entity.Question{
		Id: questionId,
		Title: request.Title,
		Detail: request.Detail,
		Image: urls,
		CommunityId: communityId,
		Questioner: userId,
		Like: make([]primitive.ObjectID, 0),
		Priority: priorityId, 
		Status: statusId, 
		Category: categoriesId,
		Answer: make([]db_entity.Answer, 0),
	}
	_, err = questionCollection.InsertOne(context.TODO(), docQuestion) // ここでMarshalBSON()される
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
    }

	// Response
	c.JSON(http.StatusOK, gin.H{
		"questionId": questionId,
	})
	return 
}