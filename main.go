package main

import (
	"html/template"
	"net/http"

	registrationController "github.com/Samuael/Projects/Inovide/controller"
	"github.com/gorilla/mux"
)

var tpl *template.Template

func init() {
	//tpl = template.Must(template.ParseGlob("C:/Users/user/go%/bin/src/gitlab.com/Mekdii/Projects/templates/*.html"))
	//tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	// fs := http.FileServer(http.Dir("./public"))
	// http.Handle("/public/", http.StripPrefix("/public/", fs))
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.HandleFunc("/register/", registrationController.RegistrationRequest).Methods("GET")
	router.HandleFunc("/register/", registrationController.RegisterUser).Methods("POST")
	router.HandleFunc("/signin/", registrationController.SignUser).Methods("GET")
	router.HandleFunc("/signin/", registrationController.SignUser).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
