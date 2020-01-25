package ChatRepository

import (
	"fmt"

	entity "github.com/Projects/Inovide/models"
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

func (chatrepo *ChatRepository) DeleteChat(message *entity.Message) []error {

	errors := chatrepo.db.Table("massage").Debug().Delete(message).GetErrors()
	return errors
}
func (chatrepo *ChatRepository) DeleteContact(user1, user2 int) []error {

	errors := chatrepo.db.Debug().Table("alies").Delete(&entity.Alie{Userid: user1, Alieid: user2}).Delete(&entity.Alie{Userid: user2, Alieid: user1}).GetErrors()

	return errors
}

func (chatrepo *ChatRepository) UpdateChat(message *entity.Message) []error {
	errors := chatrepo.db.Table("message").Debug().Save(message).GetErrors()
	return errors
}
func (chatrepo *ChatRepository) LoadAlies(person *entity.Person) ([]*entity.Person, []error) {
	alies := []entity.Alie{}

	clients := []*entity.Person{}
	errors := chatrepo.db.Debug().Table("alies").Where("userid=? or alieid=? ", person.ID, person.ID).Find(&alies).GetErrors()
	if errors != nil {
		return nil, errors
	}
	if len(alies) > 1 {
		for index, alie := range alies {
			newPerson := &entity.Person{}

			switch alie.Userid {
			case int(person.ID):
				{
					newPerson.ID = uint(alie.Alieid)
				}

			default:
				{
					newPerson.ID = uint(alie.Userid)
				}
			}
			errors = chatrepo.db.Table("users").Find(newPerson).GetErrors()
			if errors != nil {
				continue
			}
			clients[index] = newPerson

		}

	}
	return clients, errors
}

func (chatrepo *ChatRepository) LoadMessages(alie1, alie2 int) ([]*entity.Message, []error) {

	messages := []*entity.Message{}
	errors := chatrepo.db.Table("message").Debug().Where("senderid =? or recieverid=? and senderid=? or recieverid=? ", alie1, alie1, alie2, alie2).Find(messages).GetErrors()
	if errors != nil {
		return nil, errors
	}
	return messages, errors
}

//this function returns a single message pointer taking the id of the message as an argument
func (chatrepo *ChatRepository) LoadMessage(id int) (*entity.Message, []error) {
	message := &entity.Message{}
	errors := chatrepo.db.Debug().Table("message").Where(&entity.Person{}, id).Find(message).GetErrors()
	if errors != nil {
		return nil, errors
	}
	return message, nil
}

func (chatrepo *ChatRepository) DeleteMessage(message *entity.Message) []error {
	errors := chatrepo.db.Debug().Table("message").Delete(message).GetErrors()
	return errors
}

// func (chatrepo *ChatRepository)
