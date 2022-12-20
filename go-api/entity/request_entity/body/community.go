package body

type PostAddCommunityBody struct {
	CommunityId string `json:"communityId" binding:"required"`
}

