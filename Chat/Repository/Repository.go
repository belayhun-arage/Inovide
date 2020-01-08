package ChatRepository

import (
	"fmt"

	entity "github.com/Samuael/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(database *gorm.DB) *ChatRepository {
	return &ChatRepository{db: database}
}

func (chRe *ChatRepository) SaveChat(chatmessage *entity.Message) error {

	theError := chRe.db.Table("message").Debug().Create(chatmessage).Error

	if theError != nil {
		return theError
	}
	// theError = chRe.db.Table("message").Debug().Where(chatmessage).Find(chatmessage).Error
	// if theError != nil {
	// 	return theError
	// }
	return nil
}
func (chre *ChatRepository) GetId(person *entity.Person) error {
	val := &entity.Person{}
	err := chre.db.Table("users").Model(&entity.Person{}).Where("username =$1 and password=$2", person.Username, person.Password).Find(val).Error
	if err != nil {
		fmt.Println("The Error In The Repository Class is : ", err)
		return err
	}
	fmt.Println(val.ID, val.Email)
	person.ID = val.ID
	return nil
}
