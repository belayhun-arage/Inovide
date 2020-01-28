package main

import (
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	AdminRepository "github.com/Projects/Inovide/Admin/Repository"

	ChatRepository "github.com/Projects/Inovide/Chat/Repository"
	ChatService "github.com/Projects/Inovide/Chat/Service"
	CommentRepo "github.com/Projects/Inovide/Comment/Repository"
	CommentService "github.com/Projects/Inovide/Comment/Service"
	config "github.com/Projects/Inovide/DB"
	"github.com/Projects/Inovide/Idea"
	IdeaRepository "github.com/Projects/Inovide/Idea/Repository"
	ideaService "github.com/Projects/Inovide/Idea/Service"
	session "github.com/Projects/Inovide/Session"
	SessionRepo "github.com/Projects/Inovide/Session/Repository"
	repository "github.com/Projects/Inovide/User/Repository"
	service "github.com/Projects/Inovide/User/Service"
	handler "github.com/Projects/Inovide/controller"
	entity "github.com/Projects/Inovide/models"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template
var TemplateGroupUser = template.Must(template.ParseGlob("templates/*.html"))
var db *gorm.DB
var errors error
var userRepository *repository.UserRepo //1
var userservice *service.UserService    //2
var userrouter *handler.UserHandler     //3
var ideaRepository Idea.IdeaRepository  // *IdeaRepository.IdeaRepo
var ideaservice *ideaService.IdeaService
var idearouter *handler.IdeaHandler
var sessionHandler *session.Cookiehandler
var sessionrepo *SessionRepo.SessionRepository
var TheHub *entity.Hub
var chatrouter *handler.ChatHandler
var TheChatRepository *ChatRepository.ChatRepository
var TheChatService *ChatService.ChatService
var commentrouter *handler.CommentHandler
var commentservice *CommentService.CommentService
var commentrepo *CommentRepo.CommentRepo
var apicontroller *handler.ApiUserHandler
var apiideahandler *handler.ApiIdeaHandler
var adminhandler *handler.AdminHandler
var adminrepo *AdminRepository.AdminRepo

func init() {
	handler.SetSystemTemplate(TemplateGroupUser)
	/*    Initializing Users  Structure        <<Begin>>   */
	db, _ = config.InitializPostgres()
	sessionrepo = SessionRepo.NewSessionRepo(db)
	sessionHandler = session.NewCookieHandler(sessionrepo)
	userRepository = repository.NewUserRepo(db)
	userservice = service.NewUserService(userRepository)
	userrouter = handler.NewUserHandler(userservice, sessionHandler)
	/*    Initializing Users  Structure        <<End>>   */
	/*Initializing the Chat and Related Resources */
	ideaRepository = IdeaRepository.NewIdeaRepo(db)
	ideaservice = ideaService.NewIdeaService(ideaRepository)
	idearouter = handler.NewIdeaHandler(ideaservice, commentrouter, userrouter, sessionHandler)
	initChatComponents()
	initCommentComponent()
	initApiRelatedData()
	initAdmin()
	idearouter.SetCommentRepo(commentrepo)
}

func initApiRelatedData() {
	apicontroller = &handler.ApiUserHandler{}
}

func initAdmin() {
	adminrepo = AdminRepository.NewAdminRepo(db)
	adminhandler = handler.NewAdminHandler(adminrepo, userrouter, sessionHandler)
}

func initAprIdeaHandler() {
	apiideahandler = handler.NewApiIdeaHandler(userrouter, idearouter, commentrouter, sessionHandler)
}

/*This Method Will Initialize the Chay Component and Starts the Distributor Hub Of The Chat */
func initChatComponents() {
	TheHub = entity.NewHub()
	go TheHub.Run()
	TheChatRepository = ChatRepository.NewChatRepository(db)
	TheChatService = ChatService.NewChatService(TheChatRepository, TheHub)
	chatrouter = handler.NewChatHandler(TheHub, TheChatService, userservice)
	go chatrouter.MessageCoordinator()
}

func initCommentComponent() {
	commentrepo = CommentRepo.NewCommentRepo(db)
	commentservice = CommentService.NewCommentService(commentrepo)
	commentrouter = handler.NewCommentHandler(commentservice, userrouter, sessionHandler)
}

func main() {

	router := httprouter.New() //.StrictSlash(true)

	http.Handle("/", router)
	// router.GET()
	//router.ServeFiles("/src/*filepath", http.Dir("/var/www"))

	//router.ServeFiles("/public/*filepath", http.Dir("/C:/Users/user/go%/bin/src/github.com/Projects/Inovide/public/"))
	//router.ServeFiles("/public/*filepath", http.Dir("public"))
	// router.ServeFiles("/public/*filepath", http.Dir("public"))
	//	router.ServeFiles("/public/*filepath", http.Dir("public"))
	// http.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	// router.PathPrefix("/public/").Handler(http.FileServer(http.Dir("/public/")))

	// router.NotFound = http.FileServer(http.Dir("public"))

	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	router.GET("/", userrouter.ServeHome)

	//router.GET("/", userrouter.RegistrationPage)
	router.GET("/user/register/", userrouter.RegistrationPage)
	router.GET("/user/signin/", userrouter.LogInPage)
	router.POST("/user/register/", userrouter.TemplateRegistrationRequest)
	router.GET("/user/profile/", userrouter.ViewProfile)

	router.GET("/api/v1/admin/", userrouter.SearchUsers)
	router.DELETE("/api/admin/users/:userid", adminhandler.AdminDeleteIdea)
	router.DELETE("/api/admin/ideas/:ideaid", adminhandler.AdminDeleteIdea)
	router.GET("/admin/analysis/", adminhandler.AnalysisPage)
	router.POST("/admin/admins/", adminhandler.CreateAdmin)

	router.POST("/api/v1/user/register/", apicontroller.ApiRegisterUser)
	router.POST("/user/signin/", userrouter.TemplateLogInPage)
	router.PUT("/api/v1/user/update/", userrouter.ApiEditeProfile)
	router.GET("/logout/", userrouter.TemplateLogOut)
	router.GET("/idea/search/", idearouter.SearchResult)
	router.POST("/idea/new/", idearouter.TemplateCreateIdea)
	router.GET("/idea/new/", idearouter.CreateIdeaPage)
	router.PATCH("/idea/votes/", idearouter.VoteIdea) // Partials Modification of the idea meaning i am Changing the Number of Vote of The Idea
	router.PUT("/ideas/", idearouter.UpdateIdea)
	router.DELETE("/ideas/", idearouter.DeleteIdea)
	router.GET("/ideas/", idearouter.TemplateGetDetailIdea)
	router.GET("/default/idea/", idearouter.ApiGetIdea)
	router.GET("/idea/commentandpersons/", idearouter.GetCommentWithPerson)
	//_-----------------Commenting related ---------
	router.POST("/idea/comments/", commentrouter.CommentOnIdea)
	router.DELETE("/comments/comment/", commentrouter.DeleteComment)
	router.GET("/idea/myideas/", idearouter.ApiMyIdeas)
	// router.POST("/api/v1/user/signin/", apicontroller.ApiSignin)
	/***************************************Idea Related*********************************/
	// router.GET("/user/ideas/new/", idearouter.CreateIdeaPage)
	router.POST("/api/v1/user/ideas/", apiideahandler.CreateIdea)
	// Filtered  _________------------------------------------_____________________----------

	// router.POST("/user/register/", userrouter.TemplateRegisterUser)
	// router.POST("/idea/update/", idearouter.UpdateIdea)
	// // router.GET("/user/chat/", userrouter.RedirectToHome)
	// router.GET("/chat/ws", chatrouter.ChatPage)
	// router.GET("/private/user/Chat/", chatrouter.HandleChat)

	router.GET("/a/", TellMyId)

	// router.POST("/v1/user/register/", userrouter.TemplateRegisterUser)
	// // router.POST("/v1/user/signin/", userrouter.LogInRequest)
	// router.POST("/v1/idea/get/", idearouter.TemplateGetIdea)
	router.GET("/v1/idea/delete/", idearouter.DeleteIdea)
	// router.GET("/v1/idea/search/", idearouter.SearchResult)
	// router.POST("/v1/user/FollowUser/", userrouter.FollowUser)
	// /*The Comemnt Handler Related Api  */
	// router.POST("/v1/Comment/Create/", commentrouter.APICreateComment)
	// router.GET("/v1/Comment/GetComments/", commentrouter.ApiGetCommentListed)

	/**************************************  Chat Lists *******************************************************/

	router.GET("/chat/connection/", chatrouter.HandleChat)
	router.GET("/chat/ChatPage", chatrouter.ChatPage)
	// router.GET("/chat/friend/:friendid" , chatrouter.)

	router.GET("/chat/friend/", chatrouter.RecentFriends)
	router.POST("/chat/friend/", chatrouter.ConnectFriend)
	router.GET("/chat/message/:friendid", chatrouter.LoadMessages)
	// router.DELETE("/chat/message/" , )

	/**********************************************************************************************************/
	http.ListenAndServe(":8080", nil)
}
func DirectoryListener(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func TellMyId(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	cookie := http.Cookie{
		Name:    "session",
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
		Value:   "",
	}

	http.SetCookie(writer, &cookie)
	in, _, _ := sessionHandler.Valid(request)

	//fmt.Println(in)

	writer.Write([]byte(strconv.Itoa(in)))

}
