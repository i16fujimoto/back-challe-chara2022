package body

type PatchUserStatusBody struct {
	StampId string `json:"stampId" binding:"required"`
}

type SignUpBody struct {
	// Name string `json:"name" binding:"required"`
	EmailAddress string `json:"emailAddress" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginBody struct {
	EmailAddress string `json:"emailAddress" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PatchUserProfileBody struct {
	UserName string `json:"userName" binding:"required"`
	Icon []byte `json:"icon" binding:"required"`
	Profile string `json:"profile" binding:"required"`
}