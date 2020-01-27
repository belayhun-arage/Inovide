package ChatService

import (
	"fmt"

	ChatRepository "github.com/Projects/Inovide/Chat/Repository"
	entity "github.com/Projects/Inovide/models"
)

/*      Instantiation of Chat        */
type ChatService struct {
	chatRepo *ChatRepository.ChatRepository
	TheHub   *entity.Hub
}

func NewChatService(chatRepository *ChatRepository.ChatRepository, theHub *entity.Hub) *ChatService {
	return &ChatService{chatRepo: chatRepository, TheHub: theHub}
}
func (chatService *ChatService) CreateMessage(message *entity.Message) *entity.SystemMessage {
	if message.Senderid == -1 || message.Recieverid == -1 || (message.Messagedata == "" && message.Messageresource == nil) {
		return &entity.SystemMessage{
			Message:   "Can't Save The Message Invalid Message Body ",
			Succesful: false}
	}
	/*  The Service Will start serving the message to The Database Repository starting from here   */
	err := chatService.chatRepo.SaveChat(message) // Saving the Chat to the Database Repository
	if err != nil {
		fmt.Println(err, "----------------------------")
		return &entity.SystemMessage{
			Message:   "Error While Saving The Message ",
			Succesful: false,
		}
	}
	chatService.SendMessage(message)
	return &entity.SystemMessage{
		Message:   "Ok",
		Succesful: true,
	}
}
func (chatService *ChatService) SendMessage(message *entity.Message) {
	chatService.TheHub.Message <- message
}

var unit uint

func (chatService *ChatService) GetId(person *entity.Person) (*entity.SystemMessage, int64) {
	message := &entity.SystemMessage{}
	err := chatService.chatRepo.GetId(person)
	if err != nil {
		return message, -1
	}
	message.Message = "The  User Do Have An Account "
	message.Succesful = true
	return message, int64(person.ID)
}

func (chatservice *ChatService) GetFriends(friends []*entity.Person, id int) *entity.SystemMessage {

	message := &entity.SystemMessage{}
	message.Message = "You Don't Have any FRiends  Please Say Hi for some freinds "
	message.Succesful = false
	rowsaffected := chatservice.chatRepo.GetFriends(friends, id)
	if rowsaffected <= 0 {
		return message
	}
	message.Message = "Alies Fetched  "
	message.Succesful = true
	return message
}

func (chatservice *ChatService) GetMessages(yourid, friendid int, messages []*entity.Message) *entity.SystemMessage {
	systemmessage := &entity.SystemMessage{}
	systemmessage.Succesful = false
	systemmessage.Message = "InValid List Of Users d"
	rowsaffected := chatservice.chatRepo.GetMessages(yourid, friendid, messages)
	if rowsaffected <= 0 {
		return systemmessage
	}
	systemmessage.Message = "Succesful"
	systemmessage.Succesful = true
	return systemmessage
}
func (chatservice *ChatService) SaveAlies(alie1, alie2 int) *entity.SystemMessage {
	systemmessage := &entity.SystemMessage{}
	systemmessage.Message = "Not Succesful "
	systemmessage.Succesful = false
	rowsaffected := chatservice.chatRepo.SaveAlies(alie1, alie2)
	if rowsaffected <= 0 {
		return systemmessage
	}
	systemmessage.Succesful = true
	systemmessage.Message = "The Two Alies Has Registered "
	return systemmessage
}
