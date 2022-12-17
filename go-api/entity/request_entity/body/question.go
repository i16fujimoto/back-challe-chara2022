package body

type GetQuestionsBody struct {
	CommunityId string `json:"communityId" binding:"required"`
}

