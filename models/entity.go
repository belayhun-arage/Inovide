package entity

import (
	// mongodb "github.com/Samuael/Projects/Inovide/DB"
	"github.com/jinzhu/gorm"
	// "github.com/mongodb/mongo-go-driver/mongo"
)

type Person struct {
	Model     gorm.Model `gorm:"embedded" `
	ID        uint64     `json:"_id,omitempty"`                                                //  bson:"_id,omitempty"`
	Firstname string     `json:"firstname,omitempty"  gorm:"column:firstname"`                 //  bson:"firstname,omitempty"`
	Lastname  string     `json:"lastname,omitempty"  gorm:"column:lastname"`                   //  bson:"lastname,omitempty"`
	Username  string     `json:"name,omitempty" sql:"not null;unique"  gorm:"column:username"` //bson:"name,omitempty"`
	Password  string     `json:"password,omitempty" gorm:"column:password"`                    // bson:"password,omitempty"`
	Email     string     `json:"email,omitempty"  gorm:"column:email"`                         //bson:"email,omitempty"`
	Biography string     `json:"biography,omitempty" gorm:"column:biography"`                  // bson:"biography,omitempty"`
	Followers int        `json:"followers,omitempty" gorm:"column:followers" `                 //bson:"followers,omitempty"`
	Ideas     int        `json:"idea,omitempty"  gorm:"column:ideas"`                          // bson:"idea,omitempty"`
	Imagedir  string     `json:"imagdire,omitempty" gorm:"column:imagedir"`                    //  bson:"imagedirectory,omitempty"`
	Paid      int        `json:"paid,omitempty"  `                                             // bson:"paid,omitempty"`
}

type SystemMessage struct {
	Message   string `json:"message,omitempty"  bson:"message,omitempty"`
	Succesful bool   `json:"succesfull,omitempty"  bson:"succesfull,omitempty"`
}
type Message struct {
	Id              int      `json:"id,omitempty"  `
	SenderId        int      `json:"senderid,omitempty" `
	RecieverId      int      `json:"recieverid,omitempty"  `
	DateOfCreation  string   `json:"dateofcreation,omitempty"  `
	Seen            int      `json:"seen,omitempty"  `
	MessageData     string   `json:"messagedata,omitempty"  `
	MessageResource []string `json:"messageresource,omitempty"  `
}

type Idea struct {
	Id              int    `json:"id,omitempty"  `
	OwnerId         int    `json:"ownerid,omitempty"  `
	CreationData    string `json:"creationdate,omitempty"  `
	Title           string `json:"title,omitempty"  `
	Description     string `json:"description,omitempty"  `
	Visibility      string `json:"visibility,omitempty"  `
	NumberOfVotes   int    `json:"numberofvotes,omitempty"  `
	NumberOfComment int    `json:"numberofcomment,omitempty"  `
}

type Comment struct {
	Id          int    `json:"id,omitempty"  `
	IdeaId      int    `json:"ideaid,omitempty"  `
	CommentorId int    `json:"commentorid,omitempty"  `
	CommentDate string `json:"commentdate,omitempty"  `
	CommentData string `json:"commentdata,omitempty"  `
}

type Following struct {
	Id          int `json:"id,omitempty"  `
	FollowerId  int `json:"followerid,omitempty"  `
	FollowingId int `json:"followingid,omitempty"  `
}

type Alie struct {
	Id         int `json:"id,omitempty"  `
	UserId     int `json:"userid,omitempty"  `
	AlieId     int `json:"alieid,omitempty"  `
	UserOnline string
	AlieOnline string
}

// var databasemongo = mongodb.InitializeMongo()
// var Users *mongo.Collection

// func (person *Person) RegisterUser() interface{} {
// 	Users := databasemongo.Collection("User")

// 	// filteroption = option
// 	// Row := Users.FindOne(context.TODO() , bson.D{}

// 	insertInfo, erro := Users.InsertOne(context.TODO(), person)

// 	if erro != nil {
// 		return nil
// 	}

// 	return insertInfo.InsertedID
// }
// func (person *Person) FindUser(username string, password string) Message {
// 	Users = databasemongo.Collection("User")
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	errors := Users.FindOne(ctx, person).Decode(person)

// 	message := Message{}
// 	if person.Imagedir == "" {
// 		fmt.Println("There Is No User Named ", username, errors, person.Firstname, person.Lastname, person.Imagedir)
// 		message.Message = "No User Named " + username
// 		message.Succesful = false
// 		return message
// 	}

// 	message.Message = "User Does Exist !!"
// 	message.Succesful = true
// 	return message
// }
