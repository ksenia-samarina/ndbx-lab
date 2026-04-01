package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sid struct {
	HexString string
	SidString string
}

func NewSid(hexString string) *Sid {
	return &Sid{
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

func NewUser(fullName string, username string, password string, passwordHash string) *User {
	return &User{
		FullName:     fullName,
		Username:     username,
		Password:     password,
		PasswordHash: passwordHash,
	}
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Address     string             `bson:"-" json:"address,omitempty"`
	Location    struct {
		Address string `bson:"address" json:"address"`
	} `bson:"location" json:"location"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy  primitive.ObjectID `bson:"created_by" json:"created_by"`
	StartedAt  time.Time          `bson:"started_at" json:"started_at"`
	FinishedAt time.Time          `bson:"finished_at" json:"finished_at"`
}
