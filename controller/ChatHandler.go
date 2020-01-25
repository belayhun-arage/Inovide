package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	// "time"
	"crypto/rand"
	"encoding/base64"

	ChatService "github.com/Projects/Inovide/Chat/Service"
	service "github.com/Projects/Inovide/User/Service"
	entity "github.com/Projects/Inovide/models"

	// "github.com/gorilla/websocket"
	session "github.com/Projects/Inovide/Session"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
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
	TheChatService *ChatService.ChatService
	TheUserService *service.UserService
	TheHub         *entity.Hub
	session        *session.Cookiehandler
}

/*Passing the UserService and The Chat Service Returning A ChatHandler interface */
func NewChatHandler(TheHuba *entity.Hub, chatServices *ChatService.ChatService, userService *service.UserService) *ChatHandler {
	return &ChatHandler{TheChatService: chatServices, TheUserService: userService, TheHub: TheHuba}
}

/*    Main Chat Handler Instantiation             << End >>             */
func (chathandler *ChatHandler) HandleChat(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	person := &entity.Person{}
	// username, password, present := ReadSession(request)

	id, username, ok := chathandler.session.Valid(request)
	if !ok {
		return
	}
	person.Username = username
	person.ID = uint(id)
	systemMessage := chathandler.TheUserService.CheckUser(person)
	fmt.Println(person.Username, person.Email, person.ID, "_______----------->> Samuael")
	if !systemMessage.Succesful {
		// 404 Page Not Found Template Here
	}
	// chathandler.CreateWS(response, request, person)
}

//  Upgrading and Starting the Web Socket for the Incomming  Request in the header and Starting A web Socket Connectio With it
var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }
func (chathandler *ChatHandler) CreateWS(conn *websocket.Conn) {
	person := &entity.Person{}
	request := conn.Request() //  getting the request from the web socket COnnection
	id, username, present := chathandler.session.Valid(request)
	if !present {
		return
	}
	person.Username = username
	person.ID = uint(id)
	systemMessage := chathandler.TheUserService.CheckUser(person)
	if systemMessage.Succesful {

	}

	ClientId := chathandler.getClientId(person)
	client := entity.NewClient(chathandler.TheHub, conn, ClientId)
	client.TheDistributor.Register <- client
	go client.WritePump()
	go client.ReadPump()
}

// func (chathandler *ChatHandler) CeateWS(w http.ResponseWriter, r *http.Request, person *entity.Person) {
// 	wsKey, _ := generateKey()
// 	r.Header.Add("Connection", "Upgrade")
// 	r.Header.Add("Upgrade", "websocket")
// 	r.Header.Add("Sec-WebSocket-Version", "13")
// 	r.Header.Add("Sec-WebSocket-Key", wsKey)
// 	log.Printf("ws key '%v' ----  ", wsKey)
// 	// conn, err := upgrader.Upgrade(w, r, nil)

// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	ClientId := chathandler.getClientId(person)
// 	client := entity.NewClient(chathandler.TheHub, conn, ClientId)
// 	client.TheDistributor.Register <- client
// 	go client.WritePump()
// 	go client.ReadPump()
// }
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

func (cathandler *ChatHandler) LoadChatWith(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	r.ParseForm()
	stringid := r.FormValue("alieid")
	id, username, present := cathandler.session.Valid(r)
	fmt.Println(username, id, present)
	if !present {
		return
	}

	// cathandler.TheChatService.

	systemmessage := &entity.SystemMessage{}

	alieid, err := strconv.Atoi(stringid)
	if err != nil {
		systemmessage.Succesful = false
		systemmessage.Message = "Can't Load The Messaege Because of the the Id Intered is not valid " + strconv.Itoa(alieid)
		SystemTemplates.ExecuteTemplate(w, "four04.html", nil)
	}
}
func (chathandler *ChatHandler) ChatPage(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
