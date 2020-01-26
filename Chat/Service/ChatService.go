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
