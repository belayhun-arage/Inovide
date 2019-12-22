package registrationController

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	UsableFunctions "github.com/Samuael/Projects/Inovide/controller/Usables"
	usermodel "github.com/Samuael/Projects/Inovide/models"
)

var RegistrationTemplates = template.Must(template.ParseFiles("templates/reg.html", "templates/footer.html", "templates/login.html"))

func RegisterUser(w http.ResponseWriter, request *http.Request) {

	person := usermodel.Person{}

	firstname := request.FormValue("firstname")
	lasetname := request.FormValue("lastname")
	username := request.FormValue("name")
	imagedirectory, header, erro := request.FormFile("image")

	if erro != nil {
		//fmt.Println(erro)
	}
	defer imagedirectory.Close()

	email := request.FormValue("email")
	password := request.FormValue("password")
	biography := request.FormValue("biography")

	person.Email = email
	person.Username = username

	message := person.FindUser(person.Username, password)

	if message.Succesful {
		RegistrationRequestRedirect(w, request, &message)
		return
	}

	var newFullNameOfTheImageDirectory string

	//-------------------Image Saving ---------------------------

	if header.Filename != "" {

		stringSliceOfNameOfImage := strings.Split(header.Filename, ".")

		imageExtension := stringSliceOfNameOfImage[len(stringSliceOfNameOfImage)-1]

		randomStringForSavingTheImage := UsableFunctions.GenerateRandomString(20, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

		newFullNameOfTheImageDirectory = "assets/img/UsersImage/" + randomStringForSavingTheImage + "." + imageExtension

		file, errorCreatingFile := os.Create(newFullNameOfTheImageDirectory)

		if errorCreatingFile != nil {
			fmt.Println("ErrorWhile Creating the Image ", errorCreatingFile)
		}
		defer file.Close()
		person.ImageDirectory = newFullNameOfTheImageDirectory
		io.Copy(file, imagedirectory)
	}

	person.Password = password
	person.Biography = biography
	person.Firstname = firstname
	person.Lastname = lasetname

	//---------------------------------------
	_ = person.RegisterUser()

	mes, erro := json.Marshal(person)

	if erro == nil {
		w.Header().Add("Content-Type", "application/json")
		w.Write(mes)
	}
	// fmt.Println(messageBody)
	// w.Write([]byte("<h1> Succesfully Added </h1> "))
}

func RegistrationRequest(writer http.ResponseWriter, request *http.Request) {
	RegistrationTemplates.ExecuteTemplate(writer, "reg.html", nil)
}

func RegistrationRequestRedirect(writer http.ResponseWriter, request *http.Request, message *usermodel.Message) {
	RegistrationTemplates.ExecuteTemplate(writer, "reg.html", message)
}

func SignUser(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {

		RegistrationTemplates.ExecuteTemplate(writer, "login.html", nil)
	} else {
		username := request.FormValue("name")
		password := request.FormValue("password")

		person := usermodel.Person{Username: username, Password: password}

		message := person.FindUser(username, password)

		if message.Succesful {
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookiename := http.Cookie{Name: "name", Value: person.Username, Expires: expiration}
			cookiepassword := http.Cookie{Name: "password", Value: person.Password, Expires: expiration}
			http.SetCookie(writer, &cookiename)
			http.SetCookie(writer, &cookiepassword)
		}

		writer.Header().Set("Content-Type", "application/json")
		binary, erro := json.Marshal(person)

		if erro != nil {
			return
		}

		writer.Write(binary)
	}

}
