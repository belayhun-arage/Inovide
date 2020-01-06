package ChatRepository

import (
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
	return nil
}
func (chre *ChatRepository) GetId(person *entity.Person) error {
	err := chre.db.Table("message").Model(&entity.Person).Debug().Find(person).Error
	if err != nil {
		return err
	}
	return nil
}
