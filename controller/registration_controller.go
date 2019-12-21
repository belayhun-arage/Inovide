package registrationController

import (
	"fmt"
	"html/template"
	"net/http"
)

var RegistrationTemplates = template.Must(template.ParseFiles("templates/reg.html"))

func RegisterUser(w http.ResponseWriter, request *http.Request) {

}

func RegistrationRequest(w http.ResponseWriter, request *http.Request) {

	fmt.Println("I AM Called ")
	RegistrationTemplates.ExecuteTemplate(w, "templates/reg.html", nil)
}
