package user_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/s3"
	// "back-challe-chara2022/entity/db_entity"

	"net/http"
	"fmt"
	"io/ioutil"
	"context"
	"strings"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

type UserController struct {}

type UserIconResponse struct {
	UserIcon []byte `json:"userIcon"`
}

type UserCommunityResponse struct {
	UserCommunity []string `json:"userCommunity"`
}

type UserStatusResponse struct {
	IsUpdated bool `json:"isUpdated"`
}

type DocCommunity struct {
	Id primitive.ObjectID `json:"id"`
	CommunityId []string `json:"communityId"`
}

// PATCH: /user/status/<uuid: userId>
func (uc UserController) PatchUserStatus(c *gin.Context) {

	// スタンプが押された際に，userのステータスを更新

	userId, _ := primitive.ObjectIDFromHex(c.Param("userId"))
	fmt.Println(userId) // debug message

	var request body.PatchUserStatusBody
	// bodyのjsonデータを構造体にBind
	if err := c.Bind(&request); err != nil {
		// bodyのjson形式が合っていない場合
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	fmt.Println(request.StampId) // debug print

	var docStamp bson.M
	// 検索条件
	stampId, _ := primitive.ObjectIDFromHex(request.StampId)
	filterStamp := bson.D{{"_id", stampId}}
	fmt.Println(filterStamp)
	// query to stampCollection
	stampCollection := db.MongoClient.Database("insertDB").Collection("stamps")
	if err := stampCollection.FindOne(context.TODO(), filterStamp, 
		options.FindOne().SetProjection(bson.M{"status": 1, "_id": 0})).Decode(&docStamp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	} else if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the stampId")
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"message": "No document was found with the stampId",
		})
		return
	}

	fmt.Println(docStamp)
	// update raw data
	updateFields := bson.M{
		"$set": bson.M{
			"status": docStamp["status"].(string),
		},
	}
	filter := bson.M{"_id": userId}
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	result, err := userCollection.UpdateOne(context.TODO(), filter, updateFields)
	if err != nil {
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

	fmt.Println(result)

	var response UserStatusResponse

	if result.ModifiedCount > 0 {
		response.IsUpdated = true
	} else {
		response.IsUpdated = false
	}

	c.JSON(http.StatusOK, response)
	return
}

// GET: /user/community/<uuid: userId>
// $inを用いることで1つのクエリでいけるかも
func (uc UserController) GetUserCommunity(c *gin.Context) {

	// userIdが所属するコミュニティのcommunity_nameを全て返す

	userId, _ := primitive.ObjectIDFromHex(c.Param("userId"))
	fmt.Println(userId) // debug message

	var err error

	userCollection := db.MongoClient.Database("insertDB").Collection("users")

	var doc_filter bson.Raw
	// 検索条件
	filter := bson.D{{"_id", userId}}
	// query the user collection
	err = userCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(bson.M{"communityId": 1})).Decode(&doc_filter)
	if err != nil {
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
	var d_tmp DocCommunity

	// 配列の型を確定させるためにbsonを構造体に変換
	err = bson.Unmarshal(doc_filter, &d_tmp)
	if len(d_tmp.CommunityId) < 1 {
		c.JSON(http.StatusOK, gin.H{"userCommunity": d_tmp.CommunityId})
		return
	}

	var response UserCommunityResponse

	for _, doc := range d_tmp.CommunityId {

		var docCommunity bson.M
		// 検索条件
		id, _ := primitive.ObjectIDFromHex(doc)
 		filterCommunity := bson.M{"_id": id}
		// query the community collection
		communityCollection := db.MongoClient.Database("insertDB").Collection("communities")
		err = communityCollection.FindOne(context.TODO(), filterCommunity, options.FindOne().SetProjection(bson.M{"communityName": 1, "_id": 0})).Decode(&docCommunity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		} else if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the userId")
			// c.JSON(http.StatusOK, gin.H{"userCommunity": make([]string, 0)})
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"message": "No document was found with the userId",
			})
			return
		}
		fmt.Println(docCommunity["communityName"].(string))
		response.UserCommunity = append(response.UserCommunity, docCommunity["communityName"].(string))

	}

	c.JSON(http.StatusOK, response)
	return


}

// GET: /user/icon/<uuid: userId>
func (uc UserController) GetUserIcon(c *gin.Context) {

	// userIdのユーザのiconを返す

	userId, _ := primitive.ObjectIDFromHex(c.Param("userId"))
	fmt.Println(userId) // debug message

	var err error
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	var doc bson.M
	// 検索条件
	filter := bson.M{"_id": userId}
	// query
	if err := userCollection.FindOne(context.TODO(), filter, 
		options.FindOne().SetProjection(bson.M{"icon": 1, "_id": 0})).Decode(&doc); err != nil {
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

	// S3バケットとオブジェクトを指定
	url := doc["icon"].(string)
	var bucketIndex int = strings.Index(url, "/") // 最初に "/" が出現する位置
	var bucketName, key string = url[:bucketIndex], url[bucketIndex:]
	
	// S3インスタンスを作成
	s3Instance, err := s3.newS3()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code": 503,
			"message": "Service Unavailable",
		})
	}

	// S3から画像ファイルのダウンロード
	downloadKey := s3.getObjectInput(bucketName, key)
	imageData := s3.download(s3Instance, downloadKey) //[]byte



	buf, err := ioutil.ReadFile(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else {
		response := UserIconResponse{UserIcon: buf}
		c.JSON(http.StatusOK, response)
	}

}