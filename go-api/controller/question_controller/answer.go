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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
)

type QuestionResponse struct {
	QuestionId primitive.ObjectID `json:"questionId"`
	Questioner string `json:"questioner"`
	Title string `json:"title"`
	Details string `json:"details"`
	Category []string `json:"category"`
	Status string `json:"status"`
	Priority string `json:"priority"`
	Likes []Like
	CreatedAt primitive.DateTime `json:"createdAt"`
}

// AnswerResponseは配列で返す
type AnswerResponse struct {
	AnswerId primitive.ObjectID `json:"answerId"`
	Respondent string `json:"respondent"`
	Details string `json:"details"`
	Likes []Like `json:"like"`
	CreatedAt primitive.DateTime `json:"createdAt"`
}

// GET: /question/answer/<ObjectId: questionId>
// 選択された質問を返すAPI
func (qc QuestionController) GetQuestion(c *gin.Context) {

	// パスパラメータを取得
	questionId, _ := primitive.ObjectIDFromHex(c.Param("questionId"))
	
	// 質問の取得
	questionCollection := db.MongoClient.Database("insertDB").Collection("questions")
	var doc bson.M // クエリ結果を格納
	filter := bson.M{"_id": questionId}
	err := questionCollection.FindOne(context.TODO(), filter,
		options.FindOne().SetProjection(bson.M{"_id": 0, "questioner": 1, "title": 1, "detail": 1, "image": 1, "category": 1, "status": 1, "priority": 1, "like": 1, "answer": 1, "createdAt": 1})).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the questionId")
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "No document was found with the questionId",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 質問の全回答を取得する
	var answers []AnswerResponse
	for _, ans := range doc["answer"].(primitive.A) {

		// 回答投稿時間を格納
		var timestamp primitive.DateTime = ans.(primitive.M)["createdAt"].(primitive.DateTime)

		// 回答の画像の取得
		var bufArray [][]byte
		if ans.(primitive.M)["image"] != nil {
			for _, img := range ans.(primitive.M)["image"].(primitive.A) {
				var url string = img.(string)
				var bucketIndex int = strings.Index(url, "/") // 最初に "/" が出現する位置
				var bucketName, key string = url[:bucketIndex], url[bucketIndex:]
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
				bufArray = append(bufArray, buf)
			}
		}
		// いいね しているユーザの取得
		var likes []Like
		if ans.(primitive.M)["like"] != nil {
			for _, user := range ans.(primitive.M)["like"].(primitive.A) {
				var doc bson.M
				filter := bson.M{"_id": user}
				userCollection := db.MongoClient.Database("insertDB").Collection("users")
				if err := userCollection.FindOne(context.TODO(), filter, 
					options.FindOne().SetProjection(bson.M{"_id": 0, "userName": 1, "icon": 1})).Decode(&doc); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"code": http.StatusBadRequest,
						"message": err.Error(),
					})
					return
				} else if err == mongo.ErrNoDocuments {
					fmt.Printf("No document was found with the userId")
					c.JSON(http.StatusNotFound, gin.H{
						"code": 404,
						"message": "No document was found with the userId",
					})
					return
				}
			}
		}

		// 回答者名の取得
		userCollection := db.MongoClient.Database("insertDB").Collection("users")
		var docUser bson.M
		// 検索条件
		filterUser := bson.D{{"_id", ans.(primitive.M)["respondent"].(primitive.ObjectID)}}
		// query the user collection
		err = userCollection.FindOne(context.TODO(), filterUser,
			options.FindOne().SetProjection(bson.M{"userName": 1, "_id": 0})).Decode(&docUser)
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the stampId")
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"message": "No document was found with the stampId",
			})
			return
		} else if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}


		answers = append(answers, AnswerResponse{
			AnswerId: ans.(primitive.M)["_id"].(primitive.ObjectID),
			Respondent: docUser["userName"].(string),
			Details: ans.(primitive.M)["detail"].(string),
			Likes: likes,
			CreatedAt: timestamp,
		})
	}



	// 質問のカテゴリーを取得する
	var categories []string
	for _, category := range doc["category"].(primitive.A) {
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

	// ステータスを取得
	statusId := doc["status"].(primitive.ObjectID)
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
	priorityId := doc["priority"].(primitive.ObjectID)
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

	// いいね しているユーザの取得
	var likes []Like
	for _, user := range doc["like"].(primitive.A) {
		var doc bson.M
		filter := bson.M{"_id": user}
		userCollection := db.MongoClient.Database("insertDB").Collection("users")
		if err := userCollection.FindOne(context.TODO(), filter, 
			options.FindOne().SetProjection(bson.M{"_id": 0, "userName": 1, "icon": 1})).Decode(&doc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		} else if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the userId")
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"message": "No document was found with the userId",
			})
			return
		}

		// S3バケットとオブジェクトを指定
		url := doc["icon"].(string)
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
		likes = append(likes, Like {
			UserName: doc["userName"].(string),
			Icon: buf,
		})
	}
	
	// 回答者名の取得
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	var docUser bson.M
	// 検索条件
	filterUser := bson.D{{"_id", doc["questioner"].(primitive.ObjectID)}}
	// query the user collection
	err = userCollection.FindOne(context.TODO(), filterUser,
		options.FindOne().SetProjection(bson.M{"userName": 1, "_id": 1})).Decode(&docUser)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the stampId")
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "No document was found with the stampId",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 質問投稿時間を格納
	var timestamp primitive.DateTime = doc["createdAt"].(primitive.DateTime)


	var question QuestionResponse = QuestionResponse{
		QuestionId: docUser["_id"].(primitive.ObjectID),
		Questioner: docUser["userName"].(string),
		Title: doc["title"].(string),
		Details: doc["detail"].(string),
		Category: categories, 
		Status: docStatus["statusName"].(string), 
		Priority: docPriority["priorityName"].(string), 
		Likes: likes,
		CreatedAt: timestamp,
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"question": question,
		"answers": answers,
	})
	return

}

// POST: /qustion/answer/<objectID:questionId>
// 質問に対する回答を追加するAPI
// タグをうまく使うことでUpdateOneで対応可能
func (qc QuestionController) PostAnswer(c *gin.Context) {
	var err error

	// JWTよりuserIdの取得
	claims := jwt.ExtractClaims(c)
	userId, _ := primitive.ObjectIDFromHex(claims["userId"].(string))

	// パスパラメータを取得
	questionId, _ := primitive.ObjectIDFromHex(c.Param("questionId"))

	// Bodyの内容を取得
	var request body.PostAnswerBody
	if err = c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 回答を追加
	answerId := primitive.NewObjectID() // 質問ID

	// 画像のアップロード & URIの指定
	var urls []string = make([]string, 0)
	for _, obj := range request.Images {
		urls = append(urls, obj)
	}
	
	docAnswer := db_entity.Answer{
		Id: answerId,
		Detail: request.Detail,
		Image: urls,
		Respondent: userId,
		Like: make([]primitive.ObjectID, 0), // スライスを作成
	}

	questionCollection := db.MongoClient.Database("insertDB").Collection("questions")
	var doc bson.M
	filter := bson.M{"_id": questionId}
	err = questionCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(bson.M{"answer": 1, "_id": 0})).Decode(&doc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	var docAnswerArray []db_entity.Answer
	for _, ans := range doc["answer"].(primitive.A) {
		// ゼロ値処理
		var images []string = make([]string, 0)
		if ans.(primitive.M)["image"] != nil {
			for _, image := range ans.(primitive.M)["image"].(primitive.A) {
				images = append(images, image.(string))
			}
		}
		var likes []primitive.ObjectID = make([]primitive.ObjectID,0)
		if ans.(primitive.M)["like"] != nil {
			for _, like := range ans.(primitive.M)["like"].(primitive.A) {
				likes = append(likes, like.(primitive.ObjectID))
			}
		}
		docAnswerArray = append(docAnswerArray, db_entity.Answer{
			Id: ans.(primitive.M)["_id"].(primitive.ObjectID),
			Detail: ans.(primitive.M)["detail"].(string),
			Image: images,
			Respondent: ans.(primitive.M)["respondent"].(primitive.ObjectID),
			Like: likes,
		})
	}

	docAnswerArray = append(docAnswerArray, docAnswer)
	_, err = questionCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": questionId},
		bson.D{
			{"$set", bson.D{{"answer", docAnswerArray}, {"updatedAt", time.Now()}}},
		},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"answerId": answerId,
	})
	return

}