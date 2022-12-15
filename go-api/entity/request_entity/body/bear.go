package body

import (
	"time"
)

type SendBearBody struct {
	Text string `json:"text" binding:"required"`
	Bot *bool `json:"bot" binding:"required"`
}

type GetHistoryBody struct {
	Start time.Time `json:"start"`
}
