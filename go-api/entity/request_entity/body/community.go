package body

type PostAddCommunityBody struct {
	CommunityId string `json:"communityId" binding:"required"`
}

type PostMakeCommunityBody struct {
	CommunityName string `json:"communityName" binding:"required"`
	Icon []byte `json:"icon" binding:"required"`
}
