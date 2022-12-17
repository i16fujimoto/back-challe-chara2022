package server

import (
	"time"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetCORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
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
			"Access-Control-Allow-Origin",
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
	})
}