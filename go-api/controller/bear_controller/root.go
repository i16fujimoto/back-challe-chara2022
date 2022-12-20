package bear_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/entity/db_entity"
	"back-challe-chara2022/chatGPT"

	"net/http"
	"fmt"
	"time"
	"context"
	"strings"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
)

type BearController struct {}

type BearResponse struct {
	Response string `json:"response"`
}

type History struct {
	Text string `json:"text"`
	Date primitive.DateTime `json:"date"`
}

type BearHistoryResponse struct {
	Histories []History `json:"histories"`
}

// Post: /bear
func (bc BearController) PostNotLoginResponse(c *gin.Context) {
	
	var request body.SendBearBody
	// bodyのjsonデータを構造体にBind
	if err := c.BindJSON(&request); err != nil {
		// bodyのjson形式が合っていない場合
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// クマのレスポンスを返却
	var err error
	var response string

	if *request.Bot {
		response, err = chatGPT.Response(context.TODO(), []string{request.Text})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
		}
		fmt.Println(response) // debug
		var return2Index int = strings.Index(response, "\n\n")
		response = response[return2Index+2:]

	} else {
		// 現在コレクション内に入っている励まし言葉の中から1つを抽出
		bearToneCollection := db.MongoClient.Database("insertDB").Collection("bearTones")
		// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
		var cursor *mongo.Cursor
		pipeline := []bson.D{bson.D{{"$sample", bson.D{{"size", 1}}}}}
		cursor, err = bearToneCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			fmt.Println("aaaaa")
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return 
		} else if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the Responses")
			c.JSON(http.StatusNotFound, gin.H{
				"code": http.StatusNotFound,
				"massage": err.Error(),
			})
			return 
		}
		var result []bson.M
		if err := cursor.All(context.TODO(), &result); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		response = strings.Replace(result[0]["response"].(string), "<name>", "きみ", -1)
		
	}

	// 返り値
	c.JSON(http.StatusOK, BearResponse{Response: response})
	return
}


// POST: /bear/<str: userId>
func (bc BearController) PostResponse(c *gin.Context) {

	var request body.SendBearBody
	// bodyのjsonデータを構造体にBind
	if err := c.BindJSON(&request); err != nil {
		// bodyのjson形式が合っていない場合
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	claims := jwt.ExtractClaims(c)
	userId, _ := primitive.ObjectIDFromHex(claims["userId"].(string))
	fmt.Println(userId) // debug message

	// クマのレスポンスを返却
	var err error
	var response string

	if *request.Bot {

		comCollection := db.MongoClient.Database("insertDB").Collection("communications")
		// 検索条件
		timeNow := time.Now()
		filter := bson.M{
			"userId": userId, 
			"createdAt": bson.D{{"$gte", timeNow.Add(time.Minute * (-3))}},
		}

		fmt.Println(timeNow, timeNow.Add(time.Minute * (-3)))

		var cur *mongo.Cursor
		findOptions := options.Find().SetProjection(bson.M{"_id": 0, "text" : 1, "response": 1}).SetLimit(3).SetSort(bson.D{{"createdAt", -1}})
		// findOptions := options.Find().SetProjection(bson.M{"_id": 0, "messages" : 1}).SetLimit(10).SetSort(bson.M{"messages": bson.M{"createdAt": -1}})
		cur, err = comCollection.Find(context.TODO(), filter, findOptions)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}
		// 検索結果をresultsにデコード
		var results []bson.M
		if err = cur.All(context.TODO(), &results); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		var text string
		for _, r := range results {
			text += r["text"].(string) + "\n"
			text += r["response"].(string) + "\n"
		}
		text += request.Text

		fmt.Println(text)

		response, err = chatGPT.Response(context.TODO(), []string{text})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
		}
		fmt.Println(response) // debug
		var return2Index int = strings.Index(response, "\n\n")
		if return2Index >= 0 {
			response = response[return2Index+2:]
		}

	} else {
		// 現在コレクション内に入っている励まし言葉の中から1つを抽出（ランダム）
		bearToneCollection := db.MongoClient.Database("insertDB").Collection("bearTones")
		// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
		var cursor *mongo.Cursor
		pipeline := []bson.D{bson.D{{"$sample", bson.D{{"size", 1}}}}}
		cursor, err = bearToneCollection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return 
		} else if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the Responses")
			c.JSON(http.StatusNotFound, gin.H{
				"code": http.StatusNotFound,
				"massage": err.Error(),
			})
			return 
		}
		var result []bson.M
		if err := cursor.All(context.TODO(), &result); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		userCollection := db.MongoClient.Database("insertDB").Collection("users")
		var doc bson.M
		// 検索条件
		filter := bson.M{"_id": userId}
		// query
		if err := userCollection.FindOne(context.TODO(), filter, 
			options.FindOne().SetProjection(bson.M{"userName": 1, "_id": 0})).Decode(&doc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		} else if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the userId")
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"message": "No document was found with the userId",
			})
			return
		}

		// ランダムに返ってきた結果を設定
		response = strings.Replace(result[0]["response"].(string), "<name>", doc["userName"].(string), -1)
	}

	// 送られてきた内容（message）はDBに保存
	communicationCollection := db.MongoClient.Database("insertDB").Collection("communications")
	docCommunication := &db_entity.Communication{
		Id: primitive.NewObjectID(),
		UserId: userId,
		Text: request.Text, // ユーザからの入力
		Response: response, // クマからの出力
	}
	_, err = communicationCollection.InsertOne(context.TODO(), docCommunication)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 返り値
	c.JSON(http.StatusOK, BearResponse{Response: response})
	return
}


// GET: /bear/history
func (bc BearController) GetHistory(c *gin.Context) {


	// 指定されたuserIdのユーザのクマとの対話履歴を返す

	var request body.GetHistoryBody
	// bodyのjsonデータを構造体にBind
	if err := c.Bind(&request); err != nil {
		// bodyのjson形式が合っていない場合
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if request.Start.IsZero() {
		request.Start = time.Now()
	}
	
	claims := jwt.ExtractClaims(c)
	userId, _ := primitive.ObjectIDFromHex(claims["userId"].(string))
	fmt.Println(userId) // debug message

	comCollection := db.MongoClient.Database("insertDB").Collection("communications")
	// 検索条件
	filter := bson.M{
		"userId": userId, 
		"createdAt": bson.D{{"$lte", request.Start}},
	}
	var cur *mongo.Cursor
	var err error
	findOptions := options.Find().SetProjection(bson.M{"_id": 0, "text" : 1, "createdAt": 1}).SetLimit(10).SetSort(bson.D{{"createdAt", -1}})
	// findOptions := options.Find().SetProjection(bson.M{"_id": 0, "messages" : 1}).SetLimit(10).SetSort(bson.M{"messages": bson.M{"createdAt": -1}})
	cur, err = comCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Println("No document was found with the userId")
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "No document was found with the stampId",
		})
		return
	}
	// 検索結果をresultsにデコード
	var results []bson.M
	if err = cur.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	var historyArray []History = []History{}
	for _, r := range results {
		history := History{
			Text: r["text"].(string),
			Date:  r["createdAt"].(primitive.DateTime),
		}
		historyArray = append(historyArray, history)
	}

	response := BearHistoryResponse{Histories: historyArray}
	c.JSON(http.StatusOK, response)
	return
}