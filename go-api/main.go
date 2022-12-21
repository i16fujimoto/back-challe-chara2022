package main

import (
	"back-challe-chara2022/db"
    "back-challe-chara2022/server"
)

func main() {
    // DBの初期化
	db.InitDB()

    // サーバの初期化
	server.Init()
}