package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Role     string             `bson:"role" json:"role"`
	Email    string             `bson:"email" json:"email"`
	Status   string             `bson:"status" json:"status"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}
