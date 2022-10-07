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