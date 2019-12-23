package main

import (
	"html/template"
	"net/http"

	Controller "github.com/Samuael/Projects/Inovide/controller"

	"github.com/gorilla/mux"
)

var tpl *template.Template



func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	router.HandleFunc("/register/", Controller.RegistrationRequest).Methods("GET")
	router.HandleFunc("/register/", Controller.RegisterUser).Methods("POST")
	router.HandleFunc("/signin/", Controller.SignUser).Methods("GET")
	router.HandleFunc("/signin/", Controller.SignUser).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
