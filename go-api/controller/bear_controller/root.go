package bear_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/entity/db_entity"


	"net/http"
	"fmt"
	"time"
	"context"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
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

type DocTalk struct {
	Talk []db_entity.BearTone `json:"talk"`
}

// GET:  /bear
func (bc BearController) GetNotLoginResponse(c *gin.Context) {
	
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
		response = "chatGPT" 
	} else {
		// 現在コレクション内に入っている励まし言葉の中から1つを抽出
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

		response = result[0]["response"].(string)
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

	userId, _ := primitive.ObjectIDFromHex(c.Param("userId"))
	fmt.Println(userId) // debug message

	// クマのレスポンスを返却
	var err error
	var response string

	if *request.Bot {
		response = "chatGPT" 
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
		// ランダムに返ってきた結果を設定
		response = result[0]["response"].(string)
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


// GET: /bear/history/<uuid:userId>
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
	
	userId, _ := primitive.ObjectIDFromHex(c.Param("userId"))
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
	// var messages []string
	// var dates []primitive.DateTime
	for _, r := range results {
		history := History {
			Text: r["messages"].(string),
			Date:  r["createdAt"].(primitive.DateTime),
		}
		historyArray = append(historyArray, history)
		// fmt.Printf("%T\n", r["createdAt"])
		// messages = append(messages, r["messages"].(string))
		// dates = append(dates, r["createdAt"].(primitive.DateTime))
	}

	response := BearHistoryResponse{Histories: historyArray}
	c.JSON(http.StatusOK, response)
	return
}