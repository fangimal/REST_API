package user

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID           bson.ObjectID `json:"id" bson:"_id, omitempty"`
	Email        string        `json:"email" bson:"email"`
	Username     string        `json:"username" bson:"username"`
	PasswordHash string        `json:"-" bson:"password"`
}

type CreateUserDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
