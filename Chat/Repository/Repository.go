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

// func (chre *ChatRepository) GetFriends(friends []*entity.Person, id int) int64 {

// 	ids := []int{}

// 	rowsaffectd := chre.db.Debug().Table("")

// }

func (chre *ChatRepository) SaveAlies(alie1, alie2 int) int64 {
	alies := &entity.Alie{}
	alies.Userid = alie1
	alies.Alieid = alie2
	affectedRows := chre.db.Debug().Table("alies").Create(alies).RowsAffected
	defer recover()
	return affectedRows
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

func (chatrepo *ChatRepository) GetMessages(yourid, friendid int, messages []*entity.Message) int64 {
	rowsaffected := chatrepo.db.Debug().Table("message").Where("(senderid=? and recieverid =?) or (senderid=? and recieverid=? )", yourid, friendid, friendid, yourid).Find(messages).RowsAffected
	defer recover()
	return rowsaffected
}
func (chatrepo *ChatRepository) GetFriends(person []*entity.Person, id int) int64 {
	alies := []entity.Alie{}
	rowsaffectd := chatrepo.db.Debug().Table("alies").Where("userid=? or alieid=? ", id, id).Find(&alies).RowsAffected
	if rowsaffectd <= 0 {
		return -1
	}
	if len(alies) > 1 {
		for index, alie := range alies {
			newPerson := &entity.Person{}
			switch alie.Userid {
			case int(id):
				{
					newPerson.ID = uint(alie.Alieid)
				}
			default:
				{
					newPerson.ID = uint(alie.Userid)
				}
			}
			rows := chatrepo.db.Table("users").Where("id=?", newPerson.ID).Find(newPerson).RowsAffected
			if rows <= 0 {
				continue
			}
			person[index] = newPerson
		}
	}
	return rowsaffectd
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
