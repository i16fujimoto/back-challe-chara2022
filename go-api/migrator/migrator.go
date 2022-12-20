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
	responses := []string{"そうだよね", "もう一回最初から教えてよ", "そこってどういうことなの？", 
						"<name> は頑張ってるやん！", "もうちょっと詰めてみようよ！", "<name> がそんなに考えてわかんないなら，誰もわかんないよ！", 
						"今，頭が回らないだけで少し時間を空けて考えたらわかる時もあるよ！", "それはもう心が１回休めって言ってるんだよ",
						"僕は <name> のことすごいと思ってるよ", "あきらめないで！", "深呼吸をして一旦落ち着いてみよ！",
						"問題の発端を見つけないとだね．それを解決するためにはどうしたらいいんだろ？", "僕にもっと詳しく教えて！",
						"苦しいときこそ僕を頼ってよ", "んー，ちょっと僕も考えてみるね", "甘いもの食べるといいかも！とっておきのハチミツわけてあげるよ！", 
						"そんなことまでやってるの！偉いなぁ", "自分が思ってるより，<name>はすごい人だよ！", "10分くらい休憩して考えてみたら？", "後もうちょっとだよ！"}

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
	var user_id_array []db_entity.User = make([]db_entity.User, 0)
	communityId, _	:= primitive.ObjectIDFromHex("639e1e8803161570622d5263")
	docCom := &db_entity.Community{
		CommunityId: communityId,
		CommunityName: "test",
		Member: user_id_array,
		Icon: "static/icon.jpg",
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
	userId := primitive.NewObjectID()
	// community_id_array = append(community_id_array, communityId)
	docUser := db_entity.User{
		UserId: userId,
		UserName: "test",
		EmailAddress: "test@example.com",
		Password: "password",
		Icon: "static/test.png",
		Profile: "test",
		CommunityId: community_id_array,
		Status: "スッキリ",
		Role: db_entity.Role{RoleName: "admin", Permission: 7},
		Question: question_id_array,
		Like: like_id_array,
	}
	fmt.Println(docUser)

	_, err4 := userCollection.InsertOne(context.TODO(), docUser) // ここでMarshalBSON()される
	if err4 != nil {
		fmt.Println("Error inserting bear")
        panic(err4)
    } else {
		fmt.Println("Successfully inserting users")
	}

	// $push : 配列に追加
	// $pull : 配列から削除
	// $set : 変更更新
	if result, err5 := CommunityCollection.UpdateOne(context.TODO(), bson.M{"_id": communityId}, bson.D{{"$push", bson.M{"member": docUser}}, {"$set", bson.M{"updatedAt": time.Now()}}}); err5 != nil {
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

	// statusCollection init
	statusCollection := db.MongoClient.Database("insertDB").Collection("statuses")
	answerWaitId := primitive.NewObjectID()
	numariId := primitive.NewObjectID()
	finishId := primitive.NewObjectID()
	docStatus := []interface{}{
		&db_entity.Status {
			Id: answerWaitId,
			StatusName: "回答募集中",
		},
		&db_entity.Status {
			Id: numariId,
			StatusName: "ぬまり中",
		},
		&db_entity.Status {
			Id: finishId,
			StatusName: "解決済み",
		},
	}
	_, err78 := statusCollection.InsertMany(context.TODO(), docStatus) // ここでMarshalBSON()される
	if err78 != nil {
		fmt.Println("Error inserting Status")
        panic(err78)
    } else {
		fmt.Println("Successfully inserting Statuses")
	}

	// priorityCollection init
	priorityCollection := db.MongoClient.Database("insertDB").Collection("priorities")
	urgentId := primitive.NewObjectID()
	slowId := primitive.NewObjectID()
	moreSlowId := primitive.NewObjectID()
	docPriority := []interface{}{
		&db_entity.Priority {
			Id: urgentId,
			PriorityName: "緊急！",
		},
		&db_entity.Priority {
			Id: slowId,
			PriorityName: "なるはや",
		},
		&db_entity.Priority {
			Id: moreSlowId,
			PriorityName: "まったり",
		},
	}
	_, err89 := priorityCollection.InsertMany(context.TODO(), docPriority) // ここでMarshalBSON()される
	if err89 != nil {
		fmt.Println("Error inserting Priority")
        panic(err89)
    } else {
		fmt.Println("Successfully inserting Priorities")
	}

	// questionCollection init
	questionCollection := db.MongoClient.Database("insertDB").Collection("questions")
	docQuestion := []interface{}{
		&db_entity.Question {
			Id: primitive.NewObjectID(),
			Title: "pythonがなんもわからん",
			Detail: "#実装すること\n- pythonの環境構築\n- Goの環境構築\n",
			Image: []string{"static/icon.jpg", "static/test.png"},
			CommunityId: communityId,
			Questioner: docUser.UserId,
			Like: []primitive.ObjectID{docUser.UserId},
			Priority: urgentId, 
			Status: answerWaitId, 
			Category: []string{"環境構築", "pythonがまずわからん"},
			Answer: []db_entity.Answer {
				db_entity.Answer {
					Id: primitive.NewObjectID(),
					Detail: "## なんもわからん\n- pythonがまずわからん\n- Goとかもっとしらん\n- Swift神!!",
					Image: []string{"static/icon.jpg", "static/test.png"},
					Respondent: docUser.UserId,
					Like: []primitive.ObjectID{docUser.UserId},
				},	
			},
		},
		&db_entity.Question {
			Id: primitive.NewObjectID(),
			Title: "質問の取得",
			Detail: "#実装すること\n- swiftの環境構築\n- Javaの環境構築\n",
			Image: []string{"static/icon.jpg", "static/test.png"},
			CommunityId: communityId,
			Questioner: docUser.UserId,
			Like: []primitive.ObjectID{docUser.UserId},
			Priority: moreSlowId, 
			Status: numariId, 
			Category: []string{"環境構築", "Golangがまずわからん"},
			Answer: []db_entity.Answer {
				db_entity.Answer {
					Id: primitive.NewObjectID(),
					Detail: "## なんもわからん\n- pythonがまずわからん\n- Goとかもっとしらん\n- Swift神!!",
					Image: []string{"static/icon.jpg", "static/test.png"},
					Respondent: docUser.UserId,
					Like: []primitive.ObjectID{docUser.UserId},
				},
			},
		},
	}
	_, err8 := questionCollection.InsertMany(context.TODO(), docQuestion) // ここでMarshalBSON()される
	if err8 != nil {
		fmt.Println("Error inserting Question")
        panic(err8)
    } else {
		fmt.Println("Successfully inserting questions")
	}

	
}
