package main

import (
	"html/template"
	"net/http"

	config "github.com/Samuael/Projects/Inovide/DB"
	repository "github.com/Samuael/Projects/Inovide/User/Repository"
	service "github.com/Samuael/Projects/Inovide/User/Service"
	handler "github.com/Samuael/Projects/Inovide/controller"

	IdeaRepository "github.com/Samuael/Projects/Inovide/Idea/Repository"
	ideaService "github.com/Samuael/Projects/Inovide/Idea/Service"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var tpl *template.Template
var db *gorm.DB
var errors error
var userRepository *repository.UserRepo
var userservice *service.UserService
var userrouter *userhandler.UserHandler

var ideaRepository *IdeaRepository.IdeaRepo
var ideaservice *ideaService.IdeaService
var idearouter *handler.IdeaHandler

func init() {

	db, errors = config.InitializPostgres()

	if errors != nil {
		panic(errors)
	}
	userRepository = repository.NewUserRepo(db)
	userservice = service.NewUserService(userRepository)
	userrouter = userhandler.NewUserHandler(userservice)

	ideaRepository = IdeaRepository.NewIdeaRepo(db)
	ideaservice = ideaService.NewIdeaService(ideaRepository)
	idearouter = handler.NewIdeaHandler(ideaservice)

}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.HandleFunc("/register/", userrouter.RegistrationPage).Methods("GET")
	router.HandleFunc("/register/", userrouter.RegisterUser).Methods("POST")
	router.HandleFunc("/signin/", userrouter.LogInPage).Methods("GET")
	router.HandleFunc("/signin/", userrouter.LogInRequest).Methods("POST")
	router.HandleFunc("/create-idea/", idearouter.CreateIdeaPage).Methods("GET")
	router.HandleFunc("/create-idea/", idearouter.CreateIdea).Methods("POST")
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
