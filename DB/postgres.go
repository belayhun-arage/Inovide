package config

import (
	"fmt"

	// _ "github.com/lib/pq"
	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

/***************
**
Migration and  Instantiation Of The Databasee

*****************/

var db *gorm.DB
var dbs *sql.DB

var postgresStatmente string
var errors error

const (
	username = "samuael"
	password = "samuaelfirst"
	host     = "localhost"
	dbname   = "inovide"
)

func InitializPosts() (*sql.DB, error) {
	postgresStatmente = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, host, dbname)

	dbs, errors = sql.Open("postgres", postgresStatmente)

	if errors != nil {
		panic(errors)

	}

	if errors = dbs.Ping(); errors != nil {
		panic(errors)

	}
	fmt.Println("Succesfull Registration ")

	return dbs, nil

}

func InitializPostgres() (*gorm.DB, error) {
	// Preparing the statmente
	postgresStatmente = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, host, dbname)

	db, errors = gorm.Open("postgres", postgresStatmente)

	if errors != nil {
		panic(errors)

	}

	fmt.Println("Succesfull Registration ")

	return db, nil

}
func migrate() {
	db, err := gorm.Open("postgres", "user=postgres password=faniman093864 dbname=faniinovide sslmode=disable")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	db.CreateTable(Users{})
	db.CreateTable(Idea{})
	db.CreateTable(Comment{})
	db.CreateTable(Following{})
	db.CreateTable(Message{})
	db.CreateTable(Alies{})
	db.CreateTable(Votee{})
	db.CreateTable(Session{})
}

type Users struct {
	ID        int `gorm : "primary_key"`
	UserName  string
	Email     string `gorm : "not_null"`
	FirstName string `gorm : "not_null"`
	Password  string `gorm : "not_null"`
	LastName  string
	ImageDir  string `gorm:"default:'/public/img/UsersImage/thisistheRandomImage.png'"`
	Paid      int    `gorm:"default:18"`
	Biography string // `sql : Text`
	Followers int    `gorm:"default:0"`
	Ideas     int    `gorm:"default:0"`
}

type Idea struct {
	// ID SERIAL PRIMARY key not null ,
	// IdeaOwnerId serial references Users(ID) not null ,
	// createdDate Date ,
	// Title Varchar(400) not null ,
	// Description Text not null ,
	// visibility char(2) not null ,
	// NumberOfVotes INTEGER default 0 NOT NULL ,
	// NumberOfComment INTEGER DEFAULT 0 NOT NULL ,
	// Resources json

	ID              int `gorm:"primary_key;AUTO_INCREMENT;not null"`
	IdeaOwnerId     int `gorm:"not null"`
	createdDate     string
	Title           string `gorm:"type:varchar(100);not null"`
	Description     string `gorm : "not_null"`
	visibility      string `gorm : "not_null"`
	NumberOfVotes   int    `gorm:"default:0;not null"`
	NumberOfComment int    `gorm:"default:0;not null"`
	Resources       string
}

type Comment struct {
	// ID SERIAL PRIMARY KEY NOT NULL ,
	// IdeaId Serial References Idea(ID) not null,
	// CommentorId serial References Users(ID) not null ,
	// CommentData text not null ,
	// CommentedDate Date not null
	ID            int    `gorm:"primary_key;AUTO_INCREMENT;not null"`
	IdeaId        int    `gorm:"AUTO_INCREMENT;not null"`
	CommentorId   int    `gorm:"AUTO_INCREMENT;not null"`
	CommentData   string `gorm : "not_null"`
	CommentedDate string `gorm : "not_null"`
}

type Following struct {
	Id          int `gorm:"primary_key;AUTO_INCREMENT;not null"`
	FOllowerId  int `gorm:"AUTO_INCREMENT;not null"`
	FollowingId int `gorm:"AUTO_INCREMENT;not null"`
}

type Message struct {
	Id              int      `gorm:"primary_key;AUTO_INCREMENT;not null"`
	SenderId        int      `gorm:"AUTO_INCREMENT;not null"`
	RecieverId      int      `gorm:"AUTO_INCREMENT;not null"`
	DateOFCreation  string   `gorm : "not_null"`
	seen            int      `gorm:"default:0;not null"`
	MessageData     string   `gorm : "not_null"`
	MessageResource []string `gorm : "not_null"`
}

type Alies struct {
	// Id serial Primary Key not null ,
	// UserId serial References Users(ID) not null ,
	// AlieId serial References Users(ID) not null ,
	// UserOnline char(1)  default 'n',
	// AlieOnline  char(1) default 'n'

	Id int `gorm:"primary_key;AUTO_INCREMENT;not null"`
	//UserId     int    `gorm:"AUTO_INCREMENT;not null"`
	UserId     int    `sql:"type:int REFERENCES Users(id)"`
	AlieId     int    `sql:"AUTO_INCREMENT;not null;type:int REFERENCES Users(id)"`
	UserOnline string `gorm:"default:'n'"`
	AlieOnline string `gorm:"default:'n'"`
}

type Votee struct {
	// ID SERIAL PRIMARY KEY Not Null ,
	// Ideaid serial References Idea(ID) not null ,
	// Ownerid serial References Users(ID) not null

	ID      int `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Ideaid  int `gorm:"AUTO_INCREMENT;not null"`
	Ownerid int `gorm:"AUTO_INCREMENT;not null"`
}

type Session struct {
	ID       int    `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Userid   string `gorm : "not_null"`
	username string `gorm : "not_null"`
}
