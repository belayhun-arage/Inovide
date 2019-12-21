package main

import (
	registrationController "github.com/Samuael/Projects/Inovide/controller"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	//tpl = template.Must(template.ParseGlob("C:/Users/user/go%/bin/src/gitlab.com/Mekdii/Projects/templates/*.html"))
	//tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {

	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("public"))
	router.Handle("/assets/", http.StripPrefix("/assets/", fs))
	router.HandleFunc("/register/", registrationController.RegistrationRequest).Methods("GET")
	router.HandleFunc("/register/", registrationController.RegisterUser).Methods("POST")

	http.ListenAndServe(":8080", router)

}

// func start(w http.ResponseWriter, r *http.Request) {
// 	tpl.ExecuteTemplate(w, "StartUp.html", nil)

// }
// func post(w http.ResponseWriter, r *http.Request) {
// 	tpl.ExecuteTemplate(w, "PostJop.html", nil)

// }
// func raise(w http.ResponseWriter, r *http.Request) {
// 	tpl.ExecuteTemplate(w, "RaiseCapital.html", nil)
// }
// func grow(w http.ResponseWriter, r *http.Request) {
// 	tpl.ExecuteTemplate(w, "GrowStartUp.html", nil)
// }
