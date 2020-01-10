package main

import (
	"html/template"
	"net/http"

	ChatRepository "github.com/Projects/Inovide/Chat/Repository"
	ChatService "github.com/Projects/Inovide/Chat/Service"
	config "github.com/Projects/Inovide/DB"
	IdeaRepository "github.com/Projects/Inovide/Idea/Repository"
	ideaService "github.com/Projects/Inovide/Idea/Service"
	repository "github.com/Projects/Inovide/User/Repository"
	service "github.com/Projects/Inovide/User/Service"
	handler "github.com/Projects/Inovide/controller"
	entity "github.com/Projects/Inovide/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	// "html/template"
	// "net/http"
)

var tpl *template.Template
var TemplateGroupUser = template.Must(template.ParseFiles("templates/reg.html", "templates/edit.html", "templates/createIdea.html", "templates/home.html", "templates/footer.html", "templates/login.html"))
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
	router.HandleFunc("/user/chat/", userrouter.RedirectToHome).Methods("GET")

	// router.HandleFunc("/user/" , )
	// router.HandleFunc("/idea/comment/", idearouter.SaveComment).Methods("POST")
	router.HandleFunc("/ws", chatrouter.ChatPage).Methods("GET")
	router.HandleFunc("/Chat/", chatrouter.HandleChat)
	router.HandleFunc("/", ServeHome)
	http.ListenAndServe(":8080", nil)
}

func ServeHome(writer http.ResponseWriter, request *http.Request) {

	TemplateGroupUser.ExecuteTemplate(writer, "edit.html", nil)
}
