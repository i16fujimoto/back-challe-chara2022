package body

import (
	"time"
)

type SendBearBody struct {
	Text string `json:"text" binding:"required"`
}

type SendBearSentimentBody struct {
	Text string `json:"text" binding:"required"`
}

type GetHistoryBody struct {
	Start time.Time `json:"start"`
}
