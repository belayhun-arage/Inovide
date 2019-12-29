package handler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	UsableFunctions "github.com/Samuael/Projects/Inovide/Usables"
	service "github.com/Samuael/Projects/Inovide/User/Service"
	entity "github.com/Samuael/Projects/Inovide/models"
)

var (
	LENGTH_OF_IMAGE_CHARACTER        = 30
	SITE_HOST                 string = "127.0.0.1"
	SITE_PORT                        = 8080
)

// PortNumber := strconv.Itoa(SITE_PORT)
// prePath = "http://" + SITE_HOST + ":" + PortNumber ;

type UserHandler struct {
	userservice *service.UserService
}

var TemplateGroupUser = template.Must(template.ParseFiles("templates/reg.html", "templates/footer.html", "templates/login.html"))

func NewUserHandler(theService *service.UserService) *UserHandler {
	return &UserHandler{userservice: theService}
}

func (user_Admin *UserHandler) RegisterUser(writer http.ResponseWriter, request *http.Request) {
	person := entity.Person{}

	firstname := request.FormValue("firstname")
	lastname := request.FormValue("lastname")
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
	person.Password = password
	person.Biography = biography
	person.Firstname = firstname
	person.Lastname = lastname
	var newFullNameOfTheImageDirectory string
	var file *os.File
	if header.Filename != "" {

		stringSliceOfNameOfImage := strings.Split(header.Filename, ".")

		imageExtension := stringSliceOfNameOfImage[len(stringSliceOfNameOfImage)-1]

		randomStringForSavingTheImage := UsableFunctions.GenerateRandomString(LENGTH_OF_IMAGE_CHARACTER, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

		newFullNameOfTheImageDirectory = "public/img/UsersImage/" + randomStringForSavingTheImage + "." + imageExtension

		file, errorCreatingFile := os.Create(newFullNameOfTheImageDirectory)

		if errorCreatingFile != nil {
			fmt.Println("Error While Creating the Image ", errorCreatingFile)
		}
		defer file.Close()

		person.Imagedir = "/" + newFullNameOfTheImageDirectory
	}

	message := user_Admin.userservice.RegisterUser(&person)
	if message.Succesful {
		io.Copy(file, imagedirectory)
	}

	writer.Write([]byte(message.Message))

	/*

		user_Admin.ShowProfile(writer, request)
	*/
}

func (user_Admin *UserHandler) RegistrationRequestRedirect(writer http.ResponseWriter, request *http.Request) {

}

func (user_Admin *UserHandler) ShowProfile(w http.ResponseWriter, request *http.Request) {

}

func (user_controller *UserHandler) RegistrationPage(w http.ResponseWriter, request *http.Request) {
	fmt.Println("Inside Me ")
	TemplateGroupUser.ExecuteTemplate(w, "reg.html", nil)
}
