package db_entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	RoleName string `json:"roleName"` // [admin, manager, general_user]
	Permission uint `json:"permission"` // 権限を３桁 ２進数 → 10進数にしたもの (rwx)
}

type User struct {
	UserId primitive.ObjectID `json:"userId" bson:"_id"`
	UserName string `json:"userName" bson:"userName"`
	EmailAddress string `json:"emailAddress" bson:"emailAddress"`
	Password string `json:"password" bson:"password"`
	Icon string	 `json:"icon" bson:"icon"`
	Profile string `json:"profile" bson:"profile"`
	CommunityId []primitive.ObjectID `json:"communityId" bson:"communityId"`
	Status string `json:"status" bson:"status"`
	Role Role `json:"role" bson:"role"`
	Question []primitive.ObjectID `json:"question" bson:"question"`
	Like []primitive.ObjectID `json:"like" bson:"like"` // QuestionにしたLikeを保持
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
}

// Auto Create CreatedAt, UpdatedAt
func (u *User) MarshalBSON() ([]byte, error) {
    if u.CreatedAt.IsZero() {
        u.CreatedAt = time.Now()
    }
    u.UpdatedAt = time.Now()
    
    type my User
    return bson.Marshal((*my)(u))
}

type Community struct {
	CommunityId primitive.ObjectID `json:"communityId" bson:"_id"`
	CommunityName string `json:"communityName" bson:"communityName"`
	Member []primitive.ObjectID `json:"member" bson:"member"`
	Icon string `json:"icon" bson:"icon"`
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
}

// Auto Create CreatedAt, UpdatedAt
func (c *Community) MarshalBSON() ([]byte, error) {
    if c.CreatedAt.IsZero() {
        c.CreatedAt = time.Now()
    }
    c.UpdatedAt = time.Now()
    
    type my Community
    return bson.Marshal((*my)(c))
}

