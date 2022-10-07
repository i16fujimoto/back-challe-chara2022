package bear_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/entity/db_entity"


	"net/http"
	"fmt"
	"time"
	"math/rand"
	"io/ioutil"
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

type BearHistoryResponse struct {
	Messages []string `json:"message"`
	Dates []primitive.DateTime `json:"date"`
}

type BearCustomResponse struct {
	BearIcon []byte `json:"bearIcon"`
	BearTone primitive.ObjectID `json:"bearTone"`
}

type DocTalk struct {
	Talk []db_entity.Talk `json:"talk"`
}

// GET:  /bear
func (bc BearController) GetNotLoginRresponse(c *gin.Context) {
	
	// ランダムにresponseを返却

	var err error
	bearToneCollection := db.MongoClient.Database("insertDB").Collection("bearTones")
	toneId, _ := primitive.ObjectIDFromHex("633ee9f501830d402ce385c3")
	var doc_bearTone bson.Raw
	if err = bearToneCollection.FindOne(context.TODO(), bson.M{"_id": toneId}, 
		options.FindOne().SetProjection(bson.M{"talk.response": 1, "_id": 0})).Decode(&doc_bearTone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the tone_id")
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	var d_tmp DocTalk
	// 配列の型を確定させるためにbsonを構造体に変換
	err = bson.Unmarshal(doc_bearTone, &d_tmp)

	var response []string
	for _, v := range d_tmp.Talk {
		response = append(response, v.Response)
	}

	rand.Seed(time.Now().UnixNano())
    var idx int = rand.Intn(8)
	talk := BearResponse{Response: response[idx]}

	c.JSON(http.StatusOK, talk)

}


// POST: /bear/<str: user_id>
func (bc BearController) PostResponse(c *gin.Context) {
	
	// 送られてきた内容（message）はDBに保存
	// ランダムにresponseを返却

	var request body.PostSendBearBody
	// bodyのjsonデータを構造体にBind
	if err := c.Bind(&request); err != nil {
		// bodyのjson形式が合っていない場合
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	user_id, _ := primitive.ObjectIDFromHex(c.Param("user_id"))
	fmt.Println(user_id) // debug message

	// 送られてきた内容（message）はDBに保存
	var err error
	communicationCollection := db.MongoClient.Database("insertDB").Collection("communications")
	docCommunication := &db_entity.Communication{
		Id: primitive.NewObjectID(),
		UserId: user_id,
		Messages: request.Message,
	}
	_, err = communicationCollection.InsertOne(context.TODO(), docCommunication)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	// ランダムにresponseを返却
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	var doc bson.M
	// 検索条件
	filter := bson.M{"_id": user_id}
	// query
	if err := userCollection.FindOne(context.TODO(), filter, 
		options.FindOne().SetProjection(bson.M{"bearTone": 1, "_id": 0})).Decode(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the user_id")
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	bearToneCollection := db.MongoClient.Database("insertDB").Collection("bearTones")
	var doc_bearTone bson.Raw
	if err := bearToneCollection.FindOne(context.TODO(), bson.M{"_id": doc["bearTone"]}, 
		options.FindOne().SetProjection(bson.M{"talk.response": 1, "_id": 0})).Decode(&doc_bearTone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the tone_id")
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	var d_tmp DocTalk
	// 配列の型を確定させるためにbsonを構造体に変換
	err = bson.Unmarshal(doc_bearTone, &d_tmp)

	var response []string
	for _, v := range d_tmp.Talk {
		response = append(response, v.Response)
	}

	rand.Seed(time.Now().UnixNano())
    var idx int = rand.Intn(8)
	talk := BearResponse{Response: response[idx]}

	c.JSON(http.StatusOK, talk)

}


// GET: /bear/history/<uuid:user_id>
func (bc BearController) GetHistory(c *gin.Context) {


	// 指定されたuser_idのユーザのクマとの対話履歴を返す

	var request body.GetHistoryBody
	// bodyのjsonデータを構造体にBind
	if err := c.Bind(&request); err != nil {
		// bodyのjson形式が合っていない場合
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if request.Start.IsZero() {
		request.Start = time.Now()
	}
	
	user_id, _ := primitive.ObjectIDFromHex(c.Param("user_id"))
	fmt.Println(user_id) // debug message

	comCollection := db.MongoClient.Database("insertDB").Collection("communications")
	// 検索条件
	filter := bson.M{
		"userId": user_id, 
		"createdAt": bson.D{{"$lte", request.Start}},
	}
	var cur *mongo.Cursor
	var err error
	findOptions := options.Find().SetProjection(bson.M{"_id": 0, "messages" : 1, "createdAt": 1}).SetLimit(10).SetSort(bson.D{{"createdAt", -1}})
	// findOptions := options.Find().SetProjection(bson.M{"_id": 0, "messages" : 1}).SetLimit(10).SetSort(bson.M{"messages": bson.M{"createdAt": -1}})
	cur, err = comCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Println("No document was found with the user_id")
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	var results []bson.M
	if err = cur.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	var messages []string
	var dates []primitive.DateTime
	for _, r := range results {
		fmt.Printf("%T\n", r["createdAt"])
		messages = append(messages, r["messages"].(string))
		dates = append(dates, r["createdAt"].(primitive.DateTime))
	}

	response := BearHistoryResponse{Messages: messages, Dates: dates}

	c.JSON(http.StatusOK, response)
}


// GET: /bear/custom/<uuid: user_id>
func (bc BearController) GetCustom(c *gin.Context) {
	
	// user_idのユーザの，クマのiconデータ，口調idを返す

	user_id, _ := primitive.ObjectIDFromHex(c.Param("user_id"))
	fmt.Println(user_id) // debug message

	var err error
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	var doc bson.M
	// 検索条件
	filter := bson.M{"_id": user_id}
	// query
	if err := userCollection.FindOne(context.TODO(), filter, 
		options.FindOne().SetProjection(bson.M{"bearIcon": 1, "bearTone": 1, "_id": 0})).Decode(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the user_id")
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	
	bearCollection := db.MongoClient.Database("insertDB").Collection("bears")
	var doc_bearIcon bson.M
	if err := bearCollection.FindOne(context.TODO(), bson.M{"_id": doc["bearIcon"]}, 
		options.FindOne().SetProjection(bson.M{"bearIcon": 1, "_id": 0})).Decode(&doc_bearIcon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the bear_id")
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	// 画像のbyteデータ読み込み
	url := doc_bearIcon["bearIcon"].(string)
	buf, err := ioutil.ReadFile(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else {
		// primitive.ObjectID型にcast
		response := BearCustomResponse{BearIcon: buf, BearTone: doc["bearTone"].(primitive.ObjectID)}
		c.JSON(http.StatusOK, response)
		return
	}
}

