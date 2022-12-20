package community_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/s3"
	"back-challe-chara2022/entity/db_entity"

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

type CommunityController struct {}

type CommunityResponse struct {
	CommunityId primitive.ObjectID `json:"communityId"`
	CommunityName string `json:"communityName"`
	Icon []byte `json:"icon"`
}

type DocCommunity struct {
	Id primitive.ObjectID `json:"id"`
	CommunityId []string `json:"communityId"`
}

// GET: /user/community
// Userが属するコミュニティを取得
func (cc CommunityController) GetCommunity(c *gin.Context) {

	// userIdが所属するコミュニティのcommunity_nameを全て返す

	claims := jwt.ExtractClaims(c)
	userId, _ := primitive.ObjectIDFromHex(claims["userId"].(string))
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

	// コミュニティに属していない場合
	err = bson.Unmarshal(doc_filter, &d_tmp)
	if len(d_tmp.CommunityId) < 1 {
		c.JSON(http.StatusOK, gin.H{"communities": d_tmp.CommunityId})
		return
	}

	var response []CommunityResponse

	for _, doc := range d_tmp.CommunityId {

		var docCommunity bson.M
		// 検索条件
		id, _ := primitive.ObjectIDFromHex(doc)
 		filterCommunity := bson.M{"_id": id}
		// query the community collection
		communityCollection := db.MongoClient.Database("insertDB").Collection("communities")
		err = communityCollection.FindOne(context.TODO(), filterCommunity, options.FindOne().SetProjection(bson.M{"communityName": 1, "icon": 1, "_id": 0})).Decode(&docCommunity)
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
		
		// 画像の処理
		var url string = docCommunity["icon"].(string)
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

		// Responseデータに追加
		response = append(response, CommunityResponse{
			CommunityId: id,
			CommunityName: docCommunity["communityName"].(string),
			Icon: buf,
		})

	}

	c.JSON(http.StatusOK, gin.H{
		"communities": response,
	})
	return
}

// POST: /community
// ユーザがコミュニティに参加（追加）する
func(cc CommunityController) PostAddCommunity(c *gin.Context) {

	var err error

	// JWTからuserIdを取得
	claims := jwt.ExtractClaims(c)
	userId, _ := primitive.ObjectIDFromHex(claims["userId"].(string))
	fmt.Println(userId)

	// bodyの内容を取得
	var request body.PostAddCommunityBody
	if err = c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	// communityIdをbodyから取得
	communityId, _ := primitive.ObjectIDFromHex(request.CommunityId)

	// User情報のアップデート
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	filter := bson.M{"_id": userId}
	update := bson.M{"communityId": communityId}
	if _, err = userCollection.UpdateOne(context.TODO(), filter, bson.M{"$push": update}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// Community情報のアップデート
	communityCollection := db.MongoClient.Database("insertDB").Collection("communities")
	filter = bson.M{"_id": communityId}
	update = bson.M{"member": userId}
	var result bson.M
	err = communityCollection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$push": update}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"communityName": result["communityName"].(string),
	})
	return
}

// POST: /community/make
// ユーザがコミュニティを作成 ＆ 作成したコミュニティに参加する
func(cc CommunityController) PostMakeCommunity(c *gin.Context) {

	var err error

	// JWTからuserIdを取得
	claims := jwt.ExtractClaims(c)
	userId, _ := primitive.ObjectIDFromHex(claims["userId"].(string))
	fmt.Println(userId)

	// bodyの内容を取得
	var request body.PostMakeCommunityBody
	if err = c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// communityIdの作成
	communityId := primitive.NewObjectID()

	// User情報のアップデート
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	filter := bson.M{"_id": userId}
	update := bson.M{"communityId": communityId}
	var docUser bson.M
	if err = userCollection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$push": update}).Decode(&docUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}
	// 構造体にデコード
	var user db_entity.User
	var buf []byte
	// エンコード
	buf, err = bson.Marshal(docUser)
	bson.Unmarshal(buf, &user)


	// 画像をS3に登録
	var bucketName string = "static"
	var key string = "/" + communityId.Hex() + ".png"
	var url string = bucketName + key
	// S3インスタンスの作成
	s3_community, _ := s3.NewS3()
	// 画像のアップロード
	err = s3.Upload(s3_community, s3.GetPutObjectInput(bucketName, key, request.Icon))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	docCommunity := &db_entity.Community{
		CommunityId: communityId,
		CommunityName: request.CommunityName,
		Member: []db_entity.User{user},
		Icon: url,
	}
	// コミュニティコレクションに登録
	communityCollection := db.MongoClient.Database("insertDB").Collection("communities")
	if _, err = communityCollection.InsertOne(context.TODO(), docCommunity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"communityId": communityId,
	})
	return
}