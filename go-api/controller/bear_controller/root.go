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
	if err := c.Bind(&request); err != nil {
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

	if request.Bot {
		response = "chatGPT" 
	} else {
		// 現在コレクション内に入っている励まし言葉の中から1つを抽出
		bearToneCollection := db.MongoClient.Database("insertDB").Collection("bearTones")
		// filter
		filter := bson.M{ 
			"createdAt": bson.D{{"$lte", time.Now()}},
		}
		// query result
		var result bson.M
		// query
		if err = bearToneCollection.FindOne(context.TODO(), filter, 
			options.FindOne().SetProjection(bson.M{"response": 1, "_id": 0})).Decode(&result); err != nil {
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
		response = result["response"].(string)
	}

	// 返り値
	c.JSON(http.StatusOK, BearResponse{Response: response})
	return

	// // ランダムにresponseを返却
	// var err error
	// bearToneCollection := db.MongoClient.Database("insertDB").Collection("bearTones")
	// toneId, _ := primitive.ObjectIDFromHex("633ee9f501830d402ce385c3")
	// var doc_bearTone bson.Raw
	// if err = bearToneCollection.FindOne(context.TODO(), bson.M{"_id": toneId}, 
	// 	options.FindOne().SetProjection(bson.M{"talk.response": 1, "_id": 0})).Decode(&doc_bearTone); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	// 	return
	// } else if err == mongo.ErrNoDocuments {
	// 	fmt.Printf("No document was found with the toneId")
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"code": 404,
	// 		"message": "No document was found with the toneId",
	// 	})
	// 	return
	// }

	// var d_tmp DocTalk
	// // 配列の型を確定させるためにbsonを構造体に変換
	// err = bson.Unmarshal(doc_bearTone, &d_tmp)

	// var response []string
	// for _, v := range d_tmp.Talk {
	// 	response = append(response, v.Response)
	// }

	// rand.Seed(time.Now().UnixNano())
    // var idx int = rand.Intn(8)
	// talk := BearResponse{Response: response[idx]}

	// c.JSON(http.StatusOK, talk)
}


// POST: /bear/<str: userId>
func (bc BearController) PostResponse(c *gin.Context) {

	var request body.SendBearBody
	// bodyのjsonデータを構造体にBind
	if err := c.Bind(&request); err != nil {
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

	if request.Bot {
		response = "chatGPT" 
	} else {
		// 現在コレクション内に入っている励まし言葉の中から1つを抽出
		bearToneCollection := db.MongoClient.Database("insertDB").Collection("bearTones")
		// filter
		filter := bson.M{ 
			"createdAt": bson.D{{"$lte", time.Now()}},
		}
		// query result
		var result bson.M
		// query
		if err = bearToneCollection.FindOne(context.TODO(), filter, 
			options.FindOne().SetProjection(bson.M{"response": 1, "_id": 0})).Decode(&result); err != nil {
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
		response = result["response"].(string)
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
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
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