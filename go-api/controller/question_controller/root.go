package question_controller

import (
	"back-challe-chara2022/entity/request_entity/body"
	"back-challe-chara2022/db"
	"back-challe-chara2022/entity/db_entity"

	"net/http"
	"fmt"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"
)

type QuestionController struct {}

// GET: /question/priority
// 設定可能な優先度一覧を返すAPI
func (qc QuestionController) GetPriority(c *gin.Context) {
	
}

// GET: /question/status
// 設定可能なステータス一覧を返すAPI
func (qc QuestionController) GetStatus(c *gin.Context) {
	
}

// GET: /question
// 質問一覧を返すAPI
func (qc QuestionController) GetQuestions(c *gin.Context) {
	
}

