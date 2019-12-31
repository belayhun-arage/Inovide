package main

import (
	"html/template"
	"net/http"

	config "github.com/Samuael/Projects/Inovide/DB"
	repository "github.com/Samuael/Projects/Inovide/User/Repository"
	service "github.com/Samuael/Projects/Inovide/User/Service"
	userhandler "github.com/Samuael/Projects/Inovide/controller"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var tpl *template.Template
var db *gorm.DB
var errors error
var userRepository *repository.UserRepo
var userservice *service.UserService
var userrouter *userhandler.UserHandler

func init() {

<<<<<<< HEAD
=======
	db, errors = config.InitializPostgres()

	if errors != nil {
		panic(errors)
	}
	userRepository = repository.NewUserRepo(db)
	userservice = service.NewUserService(userRepository)
	userrouter = userhandler.NewUserHandler(userservice)

}

>>>>>>> origin/master
func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.HandleFunc("/register/", userrouter.RegistrationPage).Methods("GET")
	router.HandleFunc("/register/", userrouter.RegisterUser).Methods("POST")
	router.HandleFunc("/signin/", userrouter.LogInPage).Methods("GET")
	router.HandleFunc("/signin/", userrouter.LogInRequest).Methods("POST")
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
