package user_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	// "back-challe-chara2022/entity/db_entity"

	"net/http"
	"fmt"
	"io/ioutil"
	"context"
	
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

// PATCH: /user/status/<uuid: user_id>
func (uc UserController) PatchUserStatus(c *gin.Context) {

	// スタンプが押された際に，userのステータスを更新

	user_id, _ := primitive.ObjectIDFromHex(c.Param("user_id"))
	fmt.Println(user_id) // debug message

	var request body.PatchUserStatusBody
	// bodyのjsonデータを構造体にBind
	if err := c.Bind(&request); err != nil {
		// bodyのjson形式が合っていない場合
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}
	fmt.Println(request.StampId) // debug print

	var doc_stamp bson.M
	// 検索条件
	stampId, _ := primitive.ObjectIDFromHex(request.StampId)
	filter_stamp := bson.D{{"_id", stampId}}
	fmt.Println(filter_stamp)
	// query to stampCollection
	stampCollection := db.MongoClient.Database("insertDB").Collection("stamps")
	if err := stampCollection.FindOne(context.TODO(), filter_stamp, 
		options.FindOne().SetProjection(bson.M{"status": 1, "_id": 0})).Decode(&doc_stamp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
	} else if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the stamp_id")
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	fmt.Println(doc_stamp)
	// update raw data
	update_fields := bson.M{
		"$set": bson.M{
			"status": doc_stamp["status"].(string),
		},
	}
	filter := bson.M{"_id": user_id}
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	result, err := userCollection.UpdateOne(context.TODO(), filter, update_fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the user_id")
		c.JSON(http.StatusOK, gin.H{})
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

// GET: /user/community/<uuid: user_id>
// $inを用いることで1つのクエリでいけるかも
func (uc UserController) GetUserCommunity(c *gin.Context) {

	// user_idが所属するコミュニティのcommunity_nameを全て返す

	user_id, _ := primitive.ObjectIDFromHex(c.Param("user_id"))
	fmt.Println(user_id) // debug message

	var err error

	userCollection := db.MongoClient.Database("insertDB").Collection("users")

	var doc_filter bson.Raw
	// 検索条件
	filter := bson.D{{"_id", user_id}}
	// query the user collection
	err = userCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(bson.M{"communityId": 1})).Decode(&doc_filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the user_id")
		c.JSON(http.StatusOK, gin.H{})
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

		var doc_community bson.M
		// 検索条件
		id, _ := primitive.ObjectIDFromHex(doc)
 		filter_community := bson.M{"_id": id}
		// query the community collection
		communityCollection := db.MongoClient.Database("insertDB").Collection("communities")
		err = communityCollection.FindOne(context.TODO(), filter_community, options.FindOne().SetProjection(bson.M{"communityName": 1, "_id": 0})).Decode(&doc_community)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		} else if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the user_id")
			c.JSON(http.StatusOK, gin.H{"userCommunity": make([]string, 0)})
			return
		}
		fmt.Println(doc_community["communityName"].(string))
		fmt.Println("aaa")
		response.UserCommunity = append(response.UserCommunity, doc_community["communityName"].(string))

	}

	c.JSON(http.StatusOK, response)
	return


}

// GET: /user/icon/<uuid: user_id>
func (uc UserController) GetUserIcon(c *gin.Context) {

	// user_idのユーザのiconを返す

	user_id, _ := primitive.ObjectIDFromHex(c.Param("user_id"))
	fmt.Println(user_id) // debug message

	var err error
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	var doc bson.M
	// 検索条件
	filter := bson.M{"_id": user_id}
	// query
	if err := userCollection.FindOne(context.TODO(), filter, 
		options.FindOne().SetProjection(bson.M{"icon": 1, "_id": 0})).Decode(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the user_id")
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	// 画像のbyteデータ読み込み
	url := doc["icon"].(string)
	buf, err := ioutil.ReadFile(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	} else {
		response := UserIconResponse{UserIcon: buf}
		c.JSON(http.StatusOK, response)
	}

}