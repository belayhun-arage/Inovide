package entity

import (
	// mongodb "github.com/Samuael/Projects/Inovide/DB"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	// "github.com/mongodb/mongo-go-driver/mongo"
)

type Person struct {
	gorm.Model
	ID        uint64 //`gorm:"primary_key;AUTO_INCREMENT" json:"-"`                          //  bson:"_id,omitempty"`
	Firstname string `json:"firstname,omitempty"  gorm:"column:firstname"`                 //  bson:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"  gorm:"column:lastname"`                   //  bson:"lastname,omitempty"`
	Username  string `json:"name,omitempty" sql:"not null;unique"  gorm:"column:username"` //bson:"name,omitempty"`
	Password  string `json:"password,omitempty" gorm:"column:password"`                    // bson:"password,omitempty"`
	Email     string `json:"email,omitempty"  gorm:"column:email"`                         //bson:"email,omitempty"`
	Biography string `json:"biography,omitempty" gorm:"column:biography"`                  // bson:"biography,omitempty"`
	Followers int    `json:"followers,omitempty" gorm:"column:followers" `                 //bson:"followers,omitempty"`
	Ideas     int    `json:"idea,omitempty"  gorm:"column:ideas"`                          // bson:"idea,omitempty"`
	Imagedir  string `json:"imagdire,omitempty" gorm:"column:imagedir"`                    //  bson:"imagedirectory,omitempty"`
	Paid      int    `json:"paid,omitempty"  `                                             // bson:"paid,omitempty"`

}

type SystemMessage struct {
	Message   string `json:"message,omitempty"  bson:"message,omitempty"`
	Succesful bool   `json:"succesfull,omitempty"  bson:"succesfull,omitempty"`
}
type Message struct {
	Id              uint64         ` sql:"DEFAULT:user_gen_id()" json:"id" gorm:"primary_key"`
	Senderid        int            `json:"senderid" `
	Recieverid      int            `json:"recieverid"  `
	Dateofcreation  string         `json:"dateofcreation"  `
	Seen            int            `json:"seen,omitempty"  ` // minus (-1) if not Seen ++++  1-if Seen
	Messagedata     string         `json:"messagedata"  `
	Messageresource pq.StringArray `json:"messageresource"  `
}

type Idea struct {
	Id           int    `json:"id,omitempty"  `
	OwnerId      int    `json:"ownerid,omitempty"  `
	CreationData string `json:"creationdate,omitempty"  `
	Title        string `json:"title,omitempty"  `
	Description  string `json:"description,omitempty"  `
	Visibility   string `json:"visibility,omitempty"  `
	//  Consider the following into consideration while Working on Ideas
	// Use String 	"pu" --For Public  "pv" -- For Private  "pr" -- For Protected
	NumberOfVotes   int `json:"numberofvotes,omitempty"  `
	NumberOfComment int `json:"numberofcomment,omitempty"  `
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
