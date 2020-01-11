package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	chatService "github.com/Projects/Inovide/Chat/Service"
	service "github.com/Projects/Inovide/User/Service"
	"github.com/gorilla/websocket"
	"io"

	// "time"
	"crypto/rand"
	"encoding/base64"
	entity "github.com/Projects/Inovide/models"
)

/*    Main Chat Handler Instantiation                << Begin >>           */
/*
This is the Handler package Function can be accesed in Userhandler.go and ChatHandler.go Class and
We Will Be Using this method in the Main Method to distribute the template in the Main method and the templates are Created Once in the
Main.go file and used in any of the Handlers of the System (<<Handlers in the controller Directory >>)
*/
func SetSystemTemplate(temple *template.Template) {
	SystemTemplates = temple
}
func generateKey() (string, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p), nil
}

type ChatHandler struct {
	TheChatService *chatService.ChatService
	TheUserService *service.UserService
	TheHub         *entity.Hub
}

/*Passing the UserService and The Chat Service Returning A ChatHandler interface */
func NewChatHandler(TheHuba *entity.Hub, chatServices *chatService.ChatService, userService *service.UserService) *ChatHandler {
	return &ChatHandler{TheChatService: chatServices, TheUserService: userService, TheHub: TheHuba}
}

/*    Main Chat Handler Instantiation             << End >>             */
func (chathandler *ChatHandler) HandleChat(response http.ResponseWriter, request *http.Request) {
	person := &entity.Person{}
	username, password, id, present := ReadSession(request)
	if !present {
		return
	}
	person.Username = username
	person.Password = password
	person.ID = uint(id)
	systemMessage := chathandler.TheUserService.CheckUser(person)
	fmt.Println(person.Username, person.Email, person.ID, "_______----------->> Samuael")
	if !systemMessage.Succesful {
		// 404 Page Not Found Template Here
	}
	chathandler.CreateWS(response, request, person)
}

//  Upgrading and Starting the Web Socket for the Incomming  Request in the header and Starting A web Socket Connectio With it
var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (chathandler *ChatHandler) CreateWS(w http.ResponseWriter, r *http.Request, person *entity.Person) {
	wsKey, _ := generateKey()
	r.Header.Add("Connection", "Upgrade")
	r.Header.Add("Upgrade", "websocket")
	r.Header.Add("Sec-WebSocket-Version", "13")
	r.Header.Add("Sec-WebSocket-Key", wsKey)
	log.Printf("ws key '%v' ----  ", wsKey)
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	ClientId := chathandler.getClientId(person)
	client := entity.NewClient(chathandler.TheHub, conn, ClientId)
	client.TheDistributor.Register <- client
	go client.WritePump()
	go client.ReadPump()
}
func (chathandler *ChatHandler) MessageCoordinator() {
	for {
		select {
		case newMessage := <-chathandler.TheHub.Messages:
			if newMessage == nil {
				continue
			}
			SavedMessage := chathandler.SaveMesage(newMessage)
			if SavedMessage != nil {
				chathandler.TheHub.Message <- SavedMessage
			}
			// default:
			// 	close(chathandler.TheHub.Messages)
			// delete(chathandler.TheHub.Messages)
		}
	}
}
func (chathandler *ChatHandler) ChatPage(w http.ResponseWriter, r *http.Request) {
	SystemTemplates.ExecuteTemplate(w, "home.html", nil)
}
func (chathandler *ChatHandler) SaveMesage(message *entity.Message) *entity.Message {
	TheMessage := chathandler.TheChatService.CreateMessage(message)
	fmt.Println(TheMessage.Message, " Is The Message Found From The CliendService CLass and This i The Data")
	if TheMessage.Succesful {
		return message
	}
	return nil
}
func (chathandler *ChatHandler) getClientId(person *entity.Person) int {
	theServiceMesasge, id := chathandler.TheChatService.GetId(person)
	if theServiceMesasge.Succesful {
		return int(id)
	}
	return -1
}
