package body

import (
	"time"
)

type SendBearBody struct {
	Text string `json:"Text" binding:"required"`
	Bot bool `json:"Bot" binding:"required"`
}

type GetHistoryBody struct {
	Start time.Time `json:"Start"`
}