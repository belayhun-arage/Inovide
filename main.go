package main

import (
	registrationController "github.com/Samuael/Projects/Inovide/controller"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	//tpl = template.Must(template.ParseGlob("C:/Users/user/go%/bin/src/gitlab.com/Mekdii/Projects/templates/*.html"))
	//tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {

	router := mux.NewRouter()
	http.Handle("/", router)
	fs := http.FileServer(http.Dir(""))
	router.Handle("/public/", fs)
	router.HandleFunc("/register/", registrationController.RegistrationRequest).Methods("GET")
	router.HandleFunc("/register/", registrationController.RegisterUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
