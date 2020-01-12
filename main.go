package main

import (
	"html/template"
	"net/http"

	ChatRepository "github.com/Projects/Inovide/Chat/Repository"
	ChatService "github.com/Projects/Inovide/Chat/Service"
	CommentRepo "github.com/Projects/Inovide/Comment/Repository"
	CommentService "github.com/Projects/Inovide/Comment/Service"
	config "github.com/Projects/Inovide/DB"
	IdeaRepository "github.com/Projects/Inovide/Idea/Repository"
	ideaService "github.com/Projects/Inovide/Idea/Service"
	repository "github.com/Projects/Inovide/User/Repository"
	service "github.com/Projects/Inovide/User/Service"
	handler "github.com/Projects/Inovide/controller"
	entity "github.com/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template
var TemplateGroupUser = template.Must(template.ParseGlob("templates/*.*"))
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
	initCommentComponent()
	ideaRepository = IdeaRepository.NewIdeaRepo(db)
	ideaservice = ideaService.NewIdeaService(ideaRepository)
	idearouter = handler.NewIdeaHandler(ideaservice, commentrouter, userrouter)
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

var commentrouter *handler.CommentHandler
var commentservice *CommentService.CommentService
var commentrepo *CommentRepo.CommentRepo

func initCommentComponent() {
	commentrepo = CommentRepo.NewCommentRepo(db)
	commentservice = CommentService.NewCommentService(commentrepo)
	commentrouter = handler.NewCommentHandler(commentservice, userrouter)
}

func main() {

	router := httprouter.New() //.StrictSlash(true)
	http.Handle("/", router)
	// router.GET()

	router.ServeFiles("/public/*filepath", http.Dir("/home/samuael/WorkSpace/src/github.com/Projects/Inovide/public/"))
	router.GET("/", userrouter.ServeHome)
	// http.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	// router.PathPrefix("/public/").Handler(http.FileServer(http.Dir("/public/")))
	// router.NotFound = http.FileServer(http.Dir("public"))
	router.GET("/user/register/", userrouter.RegistrationPage)
	router.POST("/user/register/", userrouter.TemplateRegisterUser)
	router.GET("/user/signin/", userrouter.LogInPage)
	router.POST("/user/signin/", userrouter.LogInRequest)
	router.POST("/idea/update/", idearouter.UpdateIdea)
	router.POST("/idea/vote/", idearouter.VoteIdea)
	router.GET("/user/chat/", userrouter.RedirectToHome)
	router.GET("/chat/ws", chatrouter.ChatPage)
	router.GET("/private/user/Chat/", chatrouter.HandleChat)
	router.GET("/idea/create/", idearouter.CreateIdeaPage)

	router.GET("/v1/logout/", userrouter.LogOut)
	router.POST("/v1/user/register/", userrouter.TemplateRegisterUser)
	router.POST("/v1/user/signin/", userrouter.LogInRequest)
	router.POST("/v1/idea/create/", idearouter.CreateIdea)
	router.POST("/v1/idea/get/", idearouter.TemplateGetIdea)
	router.GET("/v1/idea/delete/", idearouter.DeleteIdea)
	router.GET("/v1/idea/search/", idearouter.SearchResult)
	router.POST("/v1/user/FollowUser/", userrouter.FollowUser)
	/*The Comemnt Handler Related Api  */
	router.POST("/v1/Comment/Create/", commentrouter.APICreateComment)
	router.GET("/v1/Comment/GetComments/", commentrouter.ApiGetCommentListed)
	http.ListenAndServe(":8080", nil)
}
