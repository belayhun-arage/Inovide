package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	chatService "github.com/Samuael/Projects/Inovide/Chat/Service"
	service "github.com/Samuael/Projects/Inovide/User/Service"

	"io"

	"github.com/gorilla/websocket"

	// "time"

	"crypto/rand"
	"encoding/base64"

	entity "github.com/Samuael/Projects/Inovide/models"
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
	username, password, present := ReadSession(request)
	if !present {
		return
	}
	person.Username = username
	person.Password = password
	systemMessage := chathandler.TheUserService.CheckUser(person)
	if !systemMessage.Succesful {
		// 404 Page Not Found Template Here
	}
	chathandler.CreateWS(response, request)
}

//  Upgrading and Starting the Web Socket for the Incomming  Request in the header and Starting A web Socket Connectio With it
var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (chathandler *ChatHandler) CreateWS(w http.ResponseWriter, r *http.Request) {
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
	client := entity.NewClient(chathandler.TheHub, conn, 0)

	client.TheDistributor.Register <- client
	fmt.Println("Messagagagagagag")

	go client.WritePump()
	go client.ReadPump()
}

func (chathandler *ChatHandler) ChatPage(w http.ResponseWriter, r *http.Request) {

	SystemTemplates.ExecuteTemplate(w, "home.html", nil)

}
