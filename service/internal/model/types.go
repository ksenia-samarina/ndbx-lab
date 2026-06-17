package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sid struct {
	HexString string
	SidString string
}

func NewSid(hexString string) Sid {
	return Sid{
		HexString: hexString,
		SidString: "sid:" + hexString,
	}
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName     string             `bson:"full_name" json:"full_name"`
	Username     string             `bson:"username" json:"username"`
	Password     string             `bson:"-" json:"password,omitempty"`
	PasswordHash string             `bson:"password_hash" json:"-"`
}

func NewUser(fullName string, username string, password string, passwordHash string) User {
	return User{
		FullName:     fullName,
		Username:     username,
		Password:     password,
		PasswordHash: passwordHash,
	}
}

type ReactionCounters struct {
	Likes    int64 `json:"likes"`
	Dislikes int64 `json:"dislikes"`
}

type ReviewsSummary struct {
	Count  int     `json:"count"`
	Rating float64 `json:"rating"`
}

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Address     string             `bson:"-" json:"address,omitempty"`
	Location    struct {
		Address string `bson:"address" json:"address"`
		City    string `bson:"city" json:"city"`
	} `bson:"location" json:"location"`
	CreatedAt  string            `bson:"created_at" json:"created_at"`
	CreatedBy  string            `bson:"created_by" json:"created_by"`
	StartedAt  string            `bson:"started_at" json:"started_at"`
	FinishedAt string            `bson:"finished_at" json:"finished_at"`
	Category   string            `bson:"category" json:"category"`
	Price      uint              `bson:"price" json:"price"`
	Reactions  *ReactionCounters `bson:"-" json:"reactions,omitempty"`
	Reviews    *ReviewsSummary   `bson:"-" json:"reviews"`
}

type EventFilter struct {
	ID        string
	Title     string
	Category  string
	PriceFrom int64
	PriceTo   int64
	City      string
	DateFrom  string
	DateTo    string
	User      string
	Offset    int64
	Limit     int64
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Reaction struct {
	EventID   string
	CreatedBy string
	LikeValue int8
	CreatedAt time.Time
}

type Review struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	Rating    int       `json:"rating"`
	UpdatedAt time.Time `json:"updated_at"`
}
