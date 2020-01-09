package main

import (
	"html/template"
	"net/http"

	ChatRepository "github.com/Samuael/Projects/Inovide/Chat/Repository"
	ChatService "github.com/Samuael/Projects/Inovide/Chat/Service"
	config "github.com/Samuael/Projects/Inovide/DB"
	IdeaRepository "github.com/Samuael/Projects/Inovide/Idea/Repository"
	ideaService "github.com/Samuael/Projects/Inovide/Idea/Service"
	repository "github.com/Samuael/Projects/Inovide/User/Repository"
	service "github.com/Samuael/Projects/Inovide/User/Service"
	handler "github.com/Samuael/Projects/Inovide/controller"
	entity "github.com/Samuael/Projects/Inovide/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var tpl *template.Template
var TemplateGroupUser = template.Must(template.ParseFiles("templates/reg.html", "templates/edit.html", "templates/createIdea.html", "templates/index.html", "templates/footer.html", "templates/login.html", "templates/sidebarr.html", "templates/four04.html"))
var db *gorm.DB
var errors error
var userRepository *repository.UserRepo
var userservice *service.UserService
var userrouter *handler.UserHandler
var ideaRepository *IdeaRepository.IdeaRepo
var ideaservice *ideaService.IdeaService
var idearouter *handler.IdeaHandler

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

	ideaRepository = IdeaRepository.NewIdeaRepo(db)
	ideaservice = ideaService.NewIdeaService(ideaRepository)
	idearouter = handler.NewIdeaHandler(ideaservice)

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
	TheChatService = ChatService.NewChatService(TheChatRepository, TheHub)
	chatrouter = handler.NewChatHandler(TheHub, TheChatService, userservice)
	go chatrouter.MessageCoordinator()
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	http.Handle("/", router)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.HandleFunc("/register/", userrouter.RegistrationPage).Methods("GET")
	router.HandleFunc("/register/", userrouter.RegisterUser).Methods("POST")
	router.HandleFunc("/signin/", userrouter.LogInPage).Methods("GET")
	router.HandleFunc("/signin/", userrouter.LogInRequest).Methods("POST")
	router.HandleFunc("/idea/create/", idearouter.CreateIdeaPage).Methods("GET")
	router.HandleFunc("/idea/create/", idearouter.CreateIdea).Methods("POST")
	router.HandleFunc("/idea/get/", idearouter.GetIdea).Methods("POST")
	router.HandleFunc("/idea/delete/", idearouter.DeleteIdea).Methods("POST")
	router.HandleFunc("/idea/update/", idearouter.UpdateIdea).Methods("POST")
	router.HandleFunc("/idea/vote/", idearouter.VoteIdea).Methods("POST")
	router.HandleFunc("/ws", chatrouter.ChatPage).Methods("GET")
	router.HandleFunc("/Chat/", chatrouter.HandleChat)
	router.HandleFunc("/", ServeHome)
	http.ListenAndServe(":8080", nil)
}

func ServeHome(writer http.ResponseWriter, request *http.Request) {

	TemplateGroupUser.ExecuteTemplate(writer, "edit.html", nil)
}
