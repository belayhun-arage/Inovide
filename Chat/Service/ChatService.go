package ChatService

import (
	//	"bytes"
	// "flag"
	// "log"
	ChatRepository "github.com/Samuael/Projects/Inovide/Chat/Repository"
	entity "github.com/Samuael/Projects/Inovide/models"
)

/*      Instantiation of Chat        */

type ChatService struct {
	chatRepo *ChatRepository.ChatRepository
}

func NewChatService(chatRepository *ChatRepository.ChatRepository) *ChatService {
	return &ChatService{chatRepo: chatRepository}
}

func (chatService *ChatService) CreateMessage(message *entity.Message) *entity.SystemMessage {
	if message.SenderId == 0 || message.RecieverId == 0 || (message.MessageData == "" && message.MessageResource == nil) {
		return &entity.SystemMessage{
			Message:   "Can't Save The Message InValid Message Body ",
			Succesful: false}
	}

	err := chatService.chatRepo.SaveChat(message)
	if err != nil {
		return &entity.SystemMessage{
			Message:   "Error While Saving The Message ",
			Succesful: false,
		}
	}
	return &entity.SystemMessage{
		Message:   "Ok",
		Succesful: true,
	}
}
