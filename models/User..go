package usermodel

import (
	"context"

	mongodb "github.com/Samuael/Projects/Inovide/DB"
)

type Person struct {
	Id             int    `json:"id,omitempty"  bson:"id,omitempty"`
	Firstname      string `json:"firstname,omitempty"  bson:"firstname,omitempty"`
	Lastname       string `json:"lastname,omitempty"  bson:"lastname,omitempty"`
	Username       string `json:"name,omitempty"  bson:"name,omitempty"`
	Password       string `json:"password,omitempty"  bson:"password,omitempty"`
	Email          string `json:"email,omitempty"  bson:"email,omitempty"`
	Biography      string `json:"biography,omitempty"  bson:"biography,omitempty"`
	Followers      int    `json:"followers,omitempty"  bson:"followers,omitempty"`
	Ideas          int    `json:"idea,omitempty"   bson:"idea,omitempty"`
	ImageDirectory string `json:"imagedirectory,omitempty"   bson:"imagedirectory,omitempty"`
}

// var Users *mongo.Collection

// var Users *mongo.Collection
// var Users *mongo.Collection
// var Users *mongo.Collection
var databasemongo = mongodb.InitializeMongo()

func (person *Person) RegisterUser() interface{} {
	Users := databasemongo.Collection("User")
	insertInfo, erro := Users.InsertOne(context.TODO(), person)

	if erro != nil {
		return nil
	}

	return insertInfo.InsertedID
}
