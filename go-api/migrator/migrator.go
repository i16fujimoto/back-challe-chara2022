package main

import(
	"context"
	"fmt"
	"time"

	"back-challe-chara2022/db"
	"back-challe-chara2022/entity/db_entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// DBマイグレート
func main() {

	// DBの初期化
	db.InitDB()

	// test success
	
	// bearToneCollection init
	bearToneCollection := db.MongoClient.Database("insertDB").Collection("bearTones")
	// 仮レスポンス
	responses := [8]string{"そうだよね", "もう一回最初から教えてよ", "そこってどういうことなの？", 
						"頑張ってるやん！", "もうちょっと詰めてみようよ！", "君がそんなに考えてわかんないなら，誰もわかんないよ！", 
						"今，頭が回らないだけで少し時間を空けて考えたらわかる時もあるよ！", "それはもう心が１回休めって言ってるんだよ"}

	for _, response := range responses {
		var err error
		docTone := &db_entity.BearTone{
			ToneId: primitive.NewObjectID(),
			Response: response,
		}
	
		_, err = bearToneCollection.InsertOne(context.TODO(), docTone) // ここでMarshalBSON()される
		if err != nil {
			fmt.Println("Error inserting bear")
			panic(err)
		} else {
			fmt.Println("Successfully inserting bear_tone")
		}
	}
	
	// communityCollection init
	CommunityCollection := db.MongoClient.Database("insertDB").Collection("communities")
	var user_id_array []primitive.ObjectID = make([]primitive.ObjectID, 0)
	communityId	:= primitive.NewObjectID()
	docCom := &db_entity.Community{
		CommunityId: communityId,
		CommunityName: "test",
		Member: user_id_array,
		Icon: "icon.jpg",
	}

	_, err3 := CommunityCollection.InsertOne(context.TODO(), docCom) // ここでMarshalBSON()される
	if err3 != nil {
		fmt.Println("Error inserting bear")
        panic(err3)
    } else {
		fmt.Println("Successfully inserting community")
	}

	// stampCollection init
	stampCollection := db.MongoClient.Database("insertDB").Collection("stamps")
	numaId, _ := primitive.ObjectIDFromHex("633ee6747701fa78b2af8486")
	sukkiriId, _ := primitive.ObjectIDFromHex("633ee6747701fa78b2af8487")
	docStamp := []interface{}{
		&db_entity.Stamp{
			StampId: numaId,
			StampName: "沼った",
			StampImg: "static/numa.png",
			Status: "ぬまった",
		},
		&db_entity.Stamp{
			StampId: sukkiriId,
			StampName: "スッキリ",
			StampImg: "static/sukkiri.png",
			Status: "スッキリ",
		},

	}
	_, err6 := stampCollection.InsertMany(context.TODO(), docStamp) // ここでMarshalBSON()される
	if err6 != nil {
		fmt.Println("Error inserting bear")
        panic(err6)
    } else {
		fmt.Println("Successfully inserting stamp")
	}

	// userCollection init
	// Test User
	userCollection := db.MongoClient.Database("insertDB").Collection("users")
	var question_id_array []primitive.ObjectID = make([]primitive.ObjectID, 0)
	var like_id_array []primitive.ObjectID = make([]primitive.ObjectID, 0)
	var community_id_array []primitive.ObjectID = make([]primitive.ObjectID, 0)
	userId, _ := primitive.ObjectIDFromHex("633ee9f501830d402ce385c5")
	// community_id_array = append(community_id_array, communityId)
	docUser := &db_entity.User{
		UserId: userId,
		UserName: "test",
		EmailAddress: "test@example.com",
		Password: "password",
		Icon: "img_dir/test.png",
		Profile: "test",
		CommunityId: community_id_array,
		Status: "スッキリ",
		Role: db_entity.Role{RoleName: "admin", Permission: 7},
		Question: question_id_array,
		Like: like_id_array,
	}
	fmt.Println(*docUser)

	_, err4 := userCollection.InsertOne(context.TODO(), docUser) // ここでMarshalBSON()される
	if err4 != nil {
		fmt.Println("Error inserting bear")
        panic(err4)
    } else {
		fmt.Println("Successfully inserting users")
	}

	user_id_array = append(user_id_array, docUser.UserId)
	fmt.Println(docUser.UserId)
	result, err5 := CommunityCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": communityId},
		bson.D{
			{"$set", bson.D{{"member", user_id_array}, {"updatedAt", time.Now()}}},
		},
	)
	if err5 != nil {
		panic(err5)
	} else {
		fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
	}

	// communicationCollection init
	communicationCollection := db.MongoClient.Database("insertDB").Collection("communications")


	docCommunication := []interface{}{
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "Hello",
			Response: "頑張ってるやん！",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "H",
			Response: "頑張ってるやん！",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "He",
			Response: "頑張ってるやん！",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "Hel",
			Response: "頑張ってるやん！",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "Hell",
			Response: "きいちろう",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "aiueo",
			Response: "古谷くんくらい真面目やん",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "zcv",
			Response: "君のそういうところが好きいちろう",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "sdfg",
			Response: "頑張ってるやん！",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "dfgh",
			Response: "頑張ってるやん！",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "osjdfo",
			Response: "頑張ってるやん！",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "lfssk",
			Response: "頑張ってるやん！",
		},
		&db_entity.Communication {
			Id: primitive.NewObjectID(),
			UserId: docUser.UserId,
			Text: "wert",
			Response: "頑張ってるやん！",
		},
	}

	_, err7 := communicationCollection.InsertMany(context.TODO(), docCommunication) // ここでMarshalBSON()される
	if err7 != nil {
		fmt.Println("Error inserting Communication")
        panic(err7)
    } else {
		fmt.Println("Successfully inserting communications")
	}
	
}
