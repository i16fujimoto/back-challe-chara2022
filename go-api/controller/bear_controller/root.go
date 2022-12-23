package bear_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/entity/db_entity"
	"back-challe-chara2022/chatGPT"
	"back-challe-chara2022/nlpAPI"

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
	Response string `json:"response"`
	Date primitive.DateTime `json:"date"`
}

type BearHistoryResponse struct {
	Histories []History `json:"histories"`
}

//POST: /bear-notlogin/sentiment
func (bc BearController) PostNotLoginSentimentResponse(c *gin.Context) {
	
	var request body.SendBearSentimentBody
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

	// NLP API
	negPhrase, sentiment, err := nlpAPI.GetTextSentiment(request.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	// 現在コレクション内に入っている中から1つを抽出
	bearToneCollection := db.MongoClient.Database("insertDB").Collection("bearTones")
	// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
	var cursor *mongo.Cursor
	matchStage := bson.D{{"$match", bson.D{{"sentiment", sentiment}}}}
	sampleStage := bson.D{{"$sample", bson.D{{"size", 1}}}}
	// pipeline := []bson.D{bson.D{{"$match", bson.D{{"sentiment", sentiment}}}, {"$sample", bson.D{{"size", 1}}}}}
	cursor, err = bearToneCollection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, sampleStage})
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
	// 返り値
	c.JSON(http.StatusOK, gin.H{
		"negPhrase": negPhrase,
		"response": response,
	})
	return

}

// POST: /bear-notlogin
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

	response, _, err = askToOthersResponse(request.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	response += "\n\n" + "Qmattaに登録すると僕がもっと話し相手になれるよ！\nコミュニティのみんなにも質問して悩みを解決しよう！"

	// 返り値
	c.JSON(http.StatusOK, BearResponse{Response: response})
	return
}

//POST: /bear/sentiment
func (bc BearController) PostSentimentResponse(c *gin.Context) {
	
	var request body.SendBearSentimentBody
	// bodyのjsonデータを構造体にBind
	if err := c.BindJSON(&request); err != nil {
		// bodyのjson形式が合っていない場合
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// userIdをJWTから取得
	claims := jwt.ExtractClaims(c)
	userId, _ := primitive.ObjectIDFromHex(claims["userId"].(string))
	fmt.Println(userId) // debug message

	// クマのレスポンスを返却
	var err error
	var response string

	// NLP API
	negPhrase, sentiment, err := nlpAPI.GetTextSentiment(request.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	// 現在コレクション内に入っている中から1つを抽出
	bearToneCollection := db.MongoClient.Database("insertDB").Collection("bearTones")
	// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
	var cursor *mongo.Cursor
	matchStage := bson.D{{"$match", bson.D{{"sentiment", sentiment}}}}
	sampleStage := bson.D{{"$sample", bson.D{{"size", 1}}}}
	// pipeline := []bson.D{bson.D{{"$match", bson.D{{"sentiment", sentiment}}}, {"$sample", bson.D{{"size", 1}}}}}
	cursor, err = bearToneCollection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, sampleStage})
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

	response = strings.Replace(result[0]["response"].(string), "<name>", doc["userName"].(string), -1)

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
	c.JSON(http.StatusOK, gin.H{
		"negPhrase": negPhrase,
		"response": response,
	})
	return

}



// POST: /bear
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
	var history string


	comCollection := db.MongoClient.Database("insertDB").Collection("communications")
	// 検索条件
	timeNow := time.Now()
	filter := bson.M{
		"userId": userId, 
		"createdAt": bson.D{{"$lte", timeNow}},
	}

	fmt.Println(timeNow, timeNow.Add(time.Minute * (-3)))

	var doc bson.M
	findOptions := options.FindOne().SetProjection(bson.M{"_id": 0, "createdAt": 1}).SetSort(bson.D{{"createdAt", 1}})
	// findOptions := options.Find().SetProjection(bson.M{"_id": 0, "messages" : 1}).SetLimit(10).SetSort(bson.M{"messages": bson.M{"createdAt": -1}})
	err = comCollection.FindOne(context.TODO(), filter, findOptions).Decode(&doc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 最後の履歴の時間
	lastMessage := doc["createdAt"].(primitive.DateTime).Time()
	// 比較対象の定義
	min4Before := time.Now().Add(time.Minute * (-4))
	min9Before := time.Now().Add(time.Minute * (-9))
	min15Before := time.Now().Add(time.Minute * (-15))
	min20Before := time.Now().Add(time.Minute * (-20))

	// fmt.Printf("%T\n", doc["createdAt"].(primitive.DateTime).Time())
	// 入力文
	fmt.Println(request.Text)

	switch {
	case lastMessage.Before(min4Before):
		response, history, err = adviceResponse(request.Text)
	case lastMessage.Before(min9Before):
		response, history, err = hintResponse(request.Text)
	case lastMessage.Before(min15Before):
		response, history, err = answerResponse(request.Text)
	case lastMessage.Before(min20Before):
		response, history, err = askToOthersResponse(request.Text)
	default:
		response, history, err = adviceResponse(request.Text)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	// 送られてきた内容（message）はDBに保存
	communicationCollection := db.MongoClient.Database("insertDB").Collection("communications")
	docCommunication := &db_entity.Communication{
		Id: primitive.NewObjectID(),
		UserId: userId,
		Text: history, // 履歴に載せる言葉
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
	findOptions := options.Find().SetProjection(bson.M{"_id": 0, "text" : 1, "response": 1, "createdAt": 1}).SetLimit(10).SetSort(bson.D{{"createdAt", -1}})
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
			Response: r["response"].(string),
			Date:  r["createdAt"].(primitive.DateTime),
		}
		historyArray = append(historyArray, history)
	}

	response := BearHistoryResponse{Histories: historyArray}
	c.JSON(http.StatusOK, response)
	return
}

// アドバイスを返す
// 1 ~ 4 min
func adviceResponse(text string) (string, string, error) {

	fmt.Println("advice")
	
	// prefix
	var prefix string = "Give me your advice on the following statement in Japanese.\n"

	var response string
	var err error

	response, err = chatGPT.Response(context.TODO(), []string{prefix + text})
	if err != nil {
		return "", "", err
	}
	fmt.Println(response) // debug
	var return2Index int = strings.Index(response, "\n\n")
	if return2Index >= 0 {
		response = response[return2Index+2:]
	}

	return response, "\\\\ クマからのアドバイス ! //", nil
}

// ヒントを返す
// 4 ~ 9 min
func hintResponse(text string) (string, string, error) {

	fmt.Println("hint")
	
	// prefix
	var prefix string = "以下の悩みを解消するヒントを教えてください。\n"

	var response string
	var err error

	response, err = chatGPT.Response(context.TODO(), []string{prefix + text})
	if err != nil {
		return "", "", err
	}
	fmt.Println(response) // debug
	var return2Index int = strings.Index(response, "\n\n")
	if return2Index >= 0 {
		response = response[return2Index+2:]
	}

	return response, "\\\\ クマからのヒント ! //", nil
}


// 答えを返す
// 9 ~ 15 min
func answerResponse(text string) (string, string, error) {
	
	fmt.Println("answer")

	var response string
	var err error
	
	response, err = chatGPT.Response(context.TODO(), []string{text})
	if err != nil {
		return "", "", err
	}
	fmt.Println(response) // debug
	var return2Index int = strings.Index(response, "\n\n")
	if return2Index >= 0 {
		response = response[return2Index+2:]
	}

	return response, "\\\\ クマの答え ! //", nil
}

// 人に聞くことを勧める
// 15 ~ 20 min
func askToOthersResponse(text string) (string, string, error) {
	
	fmt.Println("ask to others")

	response, _, err := adviceResponse(text)
	if err != nil {
		return "", "", err
	}
	var return2Index int = strings.Index(response, "\n\n")
	if return2Index >= 0 {
		response = response[return2Index+2:]
	}

	response += "\n\n" + "15分ルールっていうのがあって、僕が力になれるのはここまでかな。\n今僕に話してくれた内容から得られたヒントをもとにコミュニティのみんなに質問してみよう！"
	fmt.Println(response) // debug

	return response, "\\\\ コミュニティに聞いてみて ！ //", nil
}