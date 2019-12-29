package entity

import (
	"context"

	"fmt"
	"time"

	mongodb "github.com/Samuael/Projects/Inovide/DB"
	"go.mongodb.org/mongo-driver/mongo"
)

type Person struct {
	Id        int64  `json:"_id,omitempty"  bson:"_id,omitempty"`
	Firstname string `json:"firstname,omitempty"  bson:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"  bson:"lastname,omitempty"`
	Username  string `json:"name,omitempty"  bson:"name,omitempty"`
	Password  string `json:"password,omitempty"  bson:"password,omitempty"`
	Email     string `json:"email,omitempty"  bson:"email,omitempty"`
	Biography string `json:"biography,omitempty"  bson:"biography,omitempty"`
	Followers int    `json:"followers,omitempty"  bson:"followers,omitempty"`
	Ideas     int    `json:"idea,omitempty"   bson:"idea,omitempty"`
	Imagedir  string `json:"imagdire,omitempty"   bson:"imagedirectory,omitempty"`
	Paid      int    `json:"paid,omitempty"  bson:"paid,omitempty"`
}

type Message struct {
	Message   string `json:"message,omitempty"  bson:"message,omitempty"`
	Succesful bool   `json:"succesfull,omitempty"  bson:"succesfull,omitempty"`
}

var databasemongo = mongodb.InitializeMongo()
var Users *mongo.Collection

func (person *Person) RegisterUser() interface{} {
	Users := databasemongo.Collection("User")

	// filteroption = option
	// Row := Users.FindOne(context.TODO() , bson.D{}

	insertInfo, erro := Users.InsertOne(context.TODO(), person)

	if erro != nil {
		return nil
	}

	return insertInfo.InsertedID
}
func (person *Person) FindUser(username string, password string) Message {
	Users = databasemongo.Collection("User")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	errors := Users.FindOne(ctx, person).Decode(person)

	message := Message{}
	if person.Imagedir == "" {
		fmt.Println("There Is No User Named ", username, errors, person.Firstname, person.Lastname, person.Imagedir)
		message.Message = "No User Named " + username
		message.Succesful = false
		return message
	}

	message.Message = "User Does Exist !!"
	message.Succesful = true
	return message
}
