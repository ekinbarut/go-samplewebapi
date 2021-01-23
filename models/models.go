package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Object struct {
	Id    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string             `json:"title,omitempty" bson:"title,omitempty"`
	Desc  string             `json:"desc" bson:"desc,omitempty"`
	Body  string             `json:"body" bson:"body,omitempty"`
}
