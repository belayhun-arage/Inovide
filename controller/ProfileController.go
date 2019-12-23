package controllerHandlers

import (
	"fmt"
	usermodel "github.com/Samuael/Projects/Inovide/models"
	"html/template"
	"net/http"
)

///Foor Now I Am Going to Show The Profile Of te use r using the cookie Data Of the
// from the user  Header

var SignUpPageTemplate = template.Must(template.ParseFiles("templates/profile.html", "templates/sideNavigationBar.html"))

func ShowProfile(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("IInside The Show Profile Handler")
	username := request.FormValue("name")
	password := request.FormValue("password")

	person := usermodel.Person{}

	person.Username = username
	person.Password = password
	// This Two Parameters Are Optional I Have Added Them Just Incase They are needed
	message := person.FindUser(person.Username, person.Password)

	if message.Succesful {
		SignUpPageTemplate.ExecuteTemplate(writer, "profile.html", person)
	} else {
		RegistrationRequestRedirect(writer, request, &message)
	}

}
