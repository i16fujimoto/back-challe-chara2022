package db_entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 質問ページ

// QuestionのSubCollection
type Answer struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	Detail string `json:"detail" bson:"detail"`
	Image []string `json:"image" bson:"image"`
	Respondent primitive.ObjectID `json:"respondent" bson:"respondent"` // User ObjectID
	Like []primitive.ObjectID `json:"like" bson:"like"` // User.ObjectID
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
}

// Auto Create CreatedAt, UpdatedAt
func (a *Answer) MarshalBSON() ([]byte, error) {
    if a.CreatedAt.IsZero() {
        a.CreatedAt = time.Now()
    }
    a.UpdatedAt = time.Now()
    
    type my Answer
    return bson.Marshal((*my)(a))
}

type Status struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	StatusName string `json:"statusName" bson:"statusName"`
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
}

// Auto Create CreatedAt, UpdatedAt
func (c *Status) MarshalBSON() ([]byte, error) {
    if c.CreatedAt.IsZero() {
        c.CreatedAt = time.Now()
    }
    c.UpdatedAt = time.Now()
    
    type my Status
    return bson.Marshal((*my)(c))
}

type Priority struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	PriorityName string `json:"priorityName" bson:"priorityName"`
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
}

// Auto Create CreatedAt, UpdatedAt
func (c *Priority) MarshalBSON() ([]byte, error) {
    if c.CreatedAt.IsZero() {
        c.CreatedAt = time.Now()
    }
    c.UpdatedAt = time.Now()
    
    type my Priority
    return bson.Marshal((*my)(c))
}

type Category struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	CategoryName string `json:"categoryName" bson:"categoryName"`
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
}

// Auto Create CreatedAt, UpdatedAt
func (c *Category) MarshalBSON() ([]byte, error) {
    if c.CreatedAt.IsZero() {
        c.CreatedAt = time.Now()
    }
    c.UpdatedAt = time.Now()
    
    type my Category
    return bson.Marshal((*my)(c))
}

type Question struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Detail string `json:"detail" bson:"detail"`
	Image []string `json:"image" bson:"image"` // 質問内に挿入する画像のパス
	CommunityId primitive.ObjectID `json:"communityId" bson:"communityId"` // Community ObjectID
	Questioner primitive.ObjectID `json:"questioner" bson:"questioner"` // User ObjectID
	Like []primitive.ObjectID `json:"like" bson:"like"` // USer.ObjectID
	Priority primitive.ObjectID `json:"priority" bson:"priority"` // 緊急 or なるはや or まったり 等
	Status primitive.ObjectID `json:"status" bson:"status"` // 回答募集中 or 沼り中 or 解決済 等
	Category []primitive.ObjectID `json:"category" bson:"category"`
	// Language string `json:"language" bson:"language"` // 言語テーブル内の言語名を挿入
	Answer []Answer `json:"answer" bson:"answer"` // SubCollection
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
}

// Auto Create CreatedAt, UpdatedAt
func (q *Question) MarshalBSON() ([]byte, error) {
    if q.CreatedAt.IsZero() {
        q.CreatedAt = time.Now()
    }
    q.UpdatedAt = time.Now()
    
    type my Question
    return bson.Marshal((*my)(q))
}

// type Language struct {
// 	Id primitive.ObjectID `json:"id" bson:"_id"`
// 	LanguageName string `json:"languageName" bson:"languageName"`
// 	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
// 	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
// }

// // Auto Create CreatedAt, UpdatedAt
// func (l *Language) MarshalBSON() ([]byte, error) {
//     if l.CreatedAt.IsZero() {
//         l.CreatedAt = time.Now()
//     }
//     l.UpdatedAt = time.Now()
    
//     type my Language
//     return bson.Marshal((*my)(l))
// }

// tipsページ

// type Tips struct {
// 	Id primitive.ObjectID `json:"id" bson:"_id"`
// 	Details string `json:"details" bson:"details"`
// 	Image string `json:"image" bson:"image"`
// 	Category []string `json:"category" bson:"category"`
// 	Author User `json:"author" bson:"author"`
// 	Like []User `json:"like" bson:"like"`
// 	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
// 	UpdatedAt  time.Time  `json:"updatedAt" bson:"updatedAt"`
// 	DeletedAt  *time.Time `json:"deletedAt" bson:"deletedAt"`
// }

// // Auto Create CreatedAt, UpdatedAt
// func (t *Tips) MarshalBSON() ([]byte, error) {
//     if t.CreatedAt.IsZero() {
//         t.CreatedAt = time.Now()
//     }
//     t.UpdatedAt = time.Now()
    
//     type my Tips
//     return bson.Marshal((*my)(t))
// }