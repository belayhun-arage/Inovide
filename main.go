package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	//tpl = template.Must(template.ParseGlob("C:/Users/user/go%/bin/src/gitlab.com/Mekdii/Projects/templates/*.html"))
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", start)
	http.HandleFunc("/PostJop", post)
	http.HandleFunc("/RaiseCapital", raise)
	http.HandleFunc("/GrowStartUP", grow)
	http.ListenAndServe(":8080", nil)

}
func start(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "StartUp.html", nil)

}
func post(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "PostJop.html", nil)

}
func raise(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "RaiseCapital.html", nil)
}
func grow(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "GrowStartUp.html", nil)
}
