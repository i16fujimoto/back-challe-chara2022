package body

type PostQuestionBody struct {
	Title string `json:"title" binding:"required"`
	Detail string `json:"detail" binding:"required"`
	Image []string `json:"image"`
	Priority string `json:"priority" binding:"required"`
	Status string `json:"status" binding:"required"`
	Category []string `json:"category" binding:"required"`
}

type PostAnswerBody struct {
	Detail string `json:"detail" binding:"required"`
	Images []string `json:"images"`
}

type PatchLikeBody struct {
	QuestionId string `json:"questionId"`
	AnswerId string `json:"answerId"`
}