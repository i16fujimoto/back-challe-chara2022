package server

import (
	"back-challe-chara2022/controller/bear_controller"
	"back-challe-chara2022/controller/user_controller"
	"back-challe-chara2022/controller/login_controller"
	"back-challe-chara2022/controller/question_controller"
	"back-challe-chara2022/controller/community_controller"
	
	"os"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	jwt "github.com/appleboy/gin-jwt/v2"
)

// 初期化
func Init() {

	// 環境変数の読み込み
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// ルーティング
	r := setRouter()
	// Server Run (Port 8080)
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		panic(err)
	}
}

// ルーティング設定
func setRouter() *gin.Engine {
	
	r := gin.Default()

	// ミドルウェアの設定

	// CORSミドルウェアの定義
	r.Use(GetCORSMiddleware())
	
	// JWT認証ミドルウェアの定義
	var key string = os.Getenv("SECRET_KEY")
	authMiddleware, err := GetJWTAuthentication(key)
	if err != nil {
		panic(err)	
	}

	//ルーティング

	// user登録
	r.POST("/signup", login_controller.CreateUser)
	// ユーザ認証
	r.POST("/login", authMiddleware.(*jwt.GinJWTMiddleware).LoginHandler)

	r.NoRoute(authMiddleware.(*jwt.GinJWTMiddleware).MiddlewareFunc(), func(c *gin.Context) {
		c.JSON(404, gin.H{"code": http.StatusNotFound, "message": "Page not found"})
	})

	// ログインなしでクマを利用する
	notLogin := r.Group("/bear-notlogin")
	{
		ctrl := bear_controller.BearController{}
		// クマの返答
		notLogin.POST("", ctrl.PostNotLoginResponse)
		// Textのネガポジ判定とクマの返答
		notLogin.POST("/sentiment", ctrl.PostNotLoginSentimentResponse)
	}

	// JWT認証のミドルウェアを通すAPIを設定
	auth := r.Group("/")
	auth.GET("/refresh_token", authMiddleware.(*jwt.GinJWTMiddleware).RefreshHandler)
	auth.Use(authMiddleware.(*jwt.GinJWTMiddleware).MiddlewareFunc())
	{	
		bearGroup := auth.Group("/bear")
		{
			ctrl := bear_controller.BearController{}
			// 熊の返答を返す
			bearGroup.POST("", ctrl.PostResponse) // required login user
			// クマの定型文に必要なネガポジ判定
			bearGroup.POST("/sentiment", ctrl.PostSentimentResponse)
			// クマとの対話履歴を返す
			bearGroup.GET("/history", ctrl.GetHistory)
		}

		userGroup := auth.Group("/user")
		{
			ctrl := user_controller.UserController{}
			// user情報を返す
			userGroup.GET("", ctrl.GetUser)
			// userのステータスを更新
			userGroup.PATCH("/status", ctrl.PatchUserStatus)
			// userのアイコンを取得
			userGroup.GET("/icon", ctrl.GetUserIcon)
			// userのプロフィールの編集
			userGroup.PATCH("/profile", ctrl.PatchUserProfile)
		}

		questionGroup := auth.Group("/question")
		{
			ctrl := question_controller.QuestionController{}
			// 質問の一覧を取得
			questionGroup.GET(":communityId", ctrl.GetQuestions)
			// 質問の登録
			questionGroup.POST(":communityId", ctrl.PostQuestion)
			// 質問の取得
			questionGroup.GET("answer/:questionId", ctrl.GetQuestion)
			// 質問に回答を追加
			questionGroup.POST("answer/:questionId", ctrl.PostAnswer)
			// 優先度一覧を取得
			questionGroup.GET("/priority", ctrl.GetPriority)
			// ステータス一覧を取得
			questionGroup.GET("/status", ctrl.GetStatus)
			// カテゴリー一覧を取得
			questionGroup.GET("/category", ctrl.GetCategory)
			// 質問・回答のいいねの更新
			questionGroup.PATCH("/answer/like", ctrl.PatchLike)
		}

		communityGroup := auth.Group("/community")
		{
			ctrl := community_controller.CommunityController{}
			// userの所属するコミュニティを全て取得
			communityGroup.GET("", ctrl.GetCommunity)
			// userがコミュニティを追加（参加）
			communityGroup.POST("", ctrl.PostAddCommunity)
			// コミュニティを作成する
			communityGroup.POST("/make", ctrl.PostMakeCommunity)
			// コミュニティ内のユーザ情報一覧を取得する
			communityGroup.GET("/users/:communityId", ctrl.GetUsersInCommunity)
		}
	}
	return r
}
