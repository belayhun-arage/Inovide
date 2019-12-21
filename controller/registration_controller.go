package registrationController

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	UsableFunctions "github.com/Samuael/Projects/Inovide/controller/Usables"
	usermodel "github.com/Samuael/Projects/Inovide/models"
)

var RegistrationTemplates = template.Must(template.ParseFiles("templates/reg.html"))

type Person struct {
	Id             int    `json:"id,omitempty"  bson:"id,omitempty"`
	Firstname      string `json:"firstname,omitempty"  bson:"firstname,omitempty"`
	Lastname       string `json:"lastname,omitempty"  bson:"lastname,omitempty"`
	Username       string `json:"name,omitempty"  bson:"name,omitempty"`
	Password       string `json:"password,omitempty"  bson:"password,omitempty"`
	Email          string `json:"email,omitempty"  bson:"email,omitempty"`
	Biography      string `json:"biography,omitempty"  bson:"biography,omitempty"`
	Followers      int    `json:"followers,omitempty"  bson:"followers,omitempty"`
	Ideas          int    `json:"idea,omitempty"   bson:"idea,omitempty"`
	ImageDirectory string `json:"imagedirectory,omitempty"   bson:"imagedirectory,omitempty"`
}

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

	var newFullNameOfTheImageDirectory string

	//-------------------Image Saving ---------------------------
	if header.Filename != "" {

		stringSliceOfNameOfImage := strings.Split(header.Filename, ".")

		imageExtension := stringSliceOfNameOfImage[len(stringSliceOfNameOfImage)-1]

		randomStringForSavingTheImage := UsableFunctions.GenerateRandomString(20, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

		newFullNameOfTheImageDirectory = "/public/img/UsersImage/" + randomStringForSavingTheImage + "." + imageExtension

		file, errorCreatingFile := os.Create(newFullNameOfTheImageDirectory)
		if errorCreatingFile != nil {
			fmt.Println("ErrorWhile Creating the Image ")
		}
		defer file.Close()
		person.ImageDirectory = newFullNameOfTheImageDirectory
		io.Copy(file, imagedirectory)
	}

	person.Username = username
	person.Password = password
	person.Email = email
	person.Biography = biography
	person.Firstname = firstname
	person.Lastname = lasetname

	//---------------------------------------

	message := person.RegisterUser()
	fmt.Println(message)
	w.Write([]byte("<h1> Succesfully Added </h1> "))

}

func RegistrationRequest(writer http.ResponseWriter, request *http.Request) {
	RegistrationTemplates.Execute(writer, "templates/reg.html")
}
