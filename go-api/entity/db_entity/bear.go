package db_entity

import (
	"time"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BearTone struct {
	ToneId primitive.ObjectID `json:"toneId" bson:"_id"`
	Response string `json:"response" bson:"response"`
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
}

// Auto Create CreatedAt, UpdatedAt
func (b *BearTone) MarshalBSON() ([]byte, error) {
    if b.CreatedAt.IsZero() {
        b.CreatedAt = time.Now()
    }
    b.UpdatedAt = time.Now()
    
    type my BearTone
    return bson.Marshal((*my)(b))
}

type Communication struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	UserId primitive.ObjectID `json:"userId" bson:"userId"`
	Text string `json:"messages" bson:"text"`
	Response string `json:"response" bson:"response"`
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
}

// Auto Create CreatedAt, UpdatedAt
func (c *Communication) MarshalBSON() ([]byte, error) {
    if c.CreatedAt.IsZero() {
        c.CreatedAt = time.Now()
    }
    c.UpdatedAt = time.Now()
    
    type my Communication
    return bson.Marshal((*my)(c))
}
