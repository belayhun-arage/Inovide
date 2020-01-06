package main

import (
	ChatRepository "github.com/Samuael/Projects/Inovide/Chat/Repository"
	ChatService "github.com/Samuael/Projects/Inovide/Chat/Service"
	config "github.com/Samuael/Projects/Inovide/DB"
	repository "github.com/Samuael/Projects/Inovide/User/Repository"
	service "github.com/Samuael/Projects/Inovide/User/Service"
	handler "github.com/Samuael/Projects/Inovide/controller"
	entity "github.com/Samuael/Projects/Inovide/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
)

var tpl *template.Template
var TemplateGroupUser = template.Must(template.ParseFiles("templates/reg.html", "templates/home.html", "templates/footer.html", "templates/login.html"))

var db *gorm.DB
var errors error
var userRepository *repository.UserRepo
var userservice *service.UserService
var userrouter *handler.UserHandler

func init() {
	handler.SetSystemTemplate(TemplateGroupUser)
	/*    Initializing Users  Structure        <<Begin>>   */
	db, _ = config.InitializPostgres()
	userRepository = repository.NewUserRepo(db)
	userservice = service.NewUserService(userRepository)
	userrouter = handler.NewUserHandler(userservice)
	/*    Initializing Users  Structure        <<End>>   */

	/*Initializing the Chat and Related Resources */

	initChatComponents()
}

var TheHub *entity.Hub
var chatrouter *handler.ChatHandler
var TheChatRepository *ChatRepository.ChatRepository
var TheChatService *ChatService.ChatService

/*This Method Will Initialize the Chay Component and Starts the Distributor Hub Of The Chat */
func initChatComponents() {
	TheHub = entity.NewHub()
	go TheHub.Run()
	TheChatRepository = ChatRepository.NewChatRepository(db)
	TheChatService = ChatService.NewChatService(TheChatRepository)
	chatrouter = handler.NewChatHandler(TheHub, TheChatService, userservice)

}
func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.HandleFunc("/register/", userrouter.RegistrationPage).Methods("GET")
	router.HandleFunc("/register/", userrouter.RegisterUser).Methods("POST")
	router.HandleFunc("/signin/", userrouter.LogInPage).Methods("GET")
	router.HandleFunc("/signin/", userrouter.LogInRequest).Methods("POST")

	router.HandleFunc("/ws", chatrouter.ChatPage).Methods("GET")
	router.HandleFunc("/Chat/", chatrouter.HandleChat)
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
