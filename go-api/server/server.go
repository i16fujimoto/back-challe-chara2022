package server

import (
	"back-challe-chara2022/controller/bear_controller"
	"back-challe-chara2022/controller/user_controller"
	"back-challe-chara2022/controller/login_controller"
	
	"time"
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

// 初期化
func Init() {

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
	r.Use(cors.New(cors.Config{
		// アクセスを許可したいアクセス元
		AllowOrigins: []string{
			"http://localhost", 
			"http://localhost:3000",
			"https://qmatta.vercel.app",
		},
		// アクセスを許可したいHTTPメソッド
		AllowMethods: []string{
			"POST",
			"GET",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		// 許可したいHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: true,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	  }))

	//ルーティング
	bear_group := r.Group("bear")
	{
		ctrl := bear_controller.BearController{}
		// 熊の返答を返す
		bear_group.GET("/", ctrl.GetNotLoginRresponse) // not required login user 
		bear_group.POST(":user_id", ctrl.PostResponse) // required login user
		// クマとの対話履歴を返す
		bear_group.GET("history/:user_id", ctrl.GetHistory)
		// クマのiconデータ，口調idを返す
		bear_group.GET("custom/:user_id", ctrl.GetCustom)
	}

	user_group := r.Group("user")
	{
		ctrl := user_controller.UserController{}
		// userのステータスを更新
		user_group.PATCH("status/:user_id", ctrl.PatchUserStatus)
		// userの所属するコミュニティを全て取得
		user_group.GET("community/:user_id", ctrl.GetUserCommunity)
		// userのアイコンを取得
		user_group.GET("icon/:user_id", ctrl.GetUserIcon)	
	}

	r.POST("/signup", login_controller.CreateUser)

	return r
}

