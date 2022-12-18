package server

import (
	"back-challe-chara2022/controller/bear_controller"
	"back-challe-chara2022/controller/user_controller"
	"back-challe-chara2022/controller/login_controller"
	"back-challe-chara2022/controller/question_controller"
	
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
	notLogin := r.Group("/bear_notlogin")
	notLogin.GET("", bear_controller.BearController{}.GetNotLoginResponse)

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
			// userの所属するコミュニティを全て取得
			userGroup.GET("/community", ctrl.GetUserCommunity)
			// userのアイコンを取得
			userGroup.GET("/icon", ctrl.GetUserIcon)	
		}

		questionGroup := auth.Group("/question")
		{
			ctrl := question_controller.QuestionController{}
			// 質問の一覧を取得
			questionGroup.GET("", ctrl.GetQuestions)
			// 質問の登録
			questionGroup.POST("", ctrl.PostQuestion)
			// 質問の取得
			questionGroup.GET(":questionId", ctrl.GetQuestion)
			// 優先度一覧を取得
			questionGroup.GET("/priority", ctrl.GetPriority)
			// ステータス一覧を取得
			questionGroup.GET("/status", ctrl.GetStatus)
		}
	}
	return r
}
