package question_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/entity/db_entity"
	"back-challe-chara2022/s3"

	"net/http"
	"fmt"
	"context"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
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
	Status string `json:"status"`
	Priority string `json:"priority"`
	Questioner primitive.ObjectID `json:"questioner"`
	NumLikes int `json:"numLikes"`
	CreatedAt primitive.DateTime `json:"createdAt"`
}

type Like struct {
	UserName string `json:"userName"`
	Icon []byte `json:"icon"`
}

type QuestionResponse struct {
	Questioner string `json:"questioner"`
	Title string `json:"title"`
	Details string `json:"details"`
	Image [][]byte `json:"image"`
	Category []string `json:"category"`
	Status string `json:"status"`
	Priority string `json:"priority"`
	Likes []Like
	CreatedAt primitive.DateTime `json:"createdAt"`
}

// AnswerResponseは配列で返す
type AnswerResponse struct {
	Respondent string `json:"respondent"`
	Details string `json:"details"`
	Image [][]byte `json:"image"`
	Likes []Like `json:"like"`
	CreatedAt primitive.DateTime `json:"createdAt"`
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
		var category []string
		for _, c := range r["category"].(primitive.A) {
			category = append(category, c.(string))
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

		question := Question{
			QuestionId: r["_id"].(primitive.ObjectID),
			Title: r["title"].(string),
			Category: category,
			Priority: docPriority["priorityName"].(string),
			Status: docStatus["statusName"].(string),
			Questioner: r["questioner"].(primitive.ObjectID),
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

	// 画像のアップロード
	var urls []string
	for idx, obj := range request.Image {
		var bucketName string = "static"
		var key string = "/" + questionId.Hex() + "_" + strconv.Itoa(idx) + ".png"
		urls = append(urls, bucketName + key)
		// S3インスタンスの作成
		s3_question, _ := s3.NewS3()
		// 画像のアップロード
		err = s3.Upload(s3_question, s3.GetPutObjectInput(bucketName, key, obj))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}
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
		Category: request.Category,
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
			Respondent: docUser["userName"].(string),
			Details: ans.(primitive.M)["detail"].(string),
			Image: bufArray,
			Likes: likes,
			CreatedAt: timestamp,
		})
	}



	// 質問のカテゴリーを取得する
	var categories []string
	for _, category := range doc["category"].(primitive.A) {
		categories = append(categories, category.(string))
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

	// 質問の画像の取得
	var bufArray [][]byte
	if doc["image"] != nil {
		for _, img := range doc["image"].(primitive.A) {
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

	// 質問投稿時間を格納
	var timestamp primitive.DateTime = doc["createdAt"].(primitive.DateTime)


	var question QuestionResponse = QuestionResponse{
		Questioner: docUser["userName"].(string),
		Title: doc["title"].(string),
		Details: doc["detail"].(string),
		Image: bufArray, 
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
	for idx, obj := range request.Images {
		var bucketName string = "static"
		var key string = "/" + questionId.Hex() + "_" + answerId.Hex() + "_" + strconv.Itoa(idx) + ".png"
		urls = append(urls, bucketName + key)
		// S3インスタンスの作成
		s3_answer, _ := s3.NewS3()
		// 画像のアップロード
		err = s3.Upload(s3_answer, s3.GetPutObjectInput(bucketName, key, obj))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}
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