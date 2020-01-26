package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	session "github.com/Projects/Inovide/Session"
	UsableFunctions "github.com/Projects/Inovide/Usables"
	service "github.com/Projects/Inovide/User/Service"
	entity "github.com/Projects/Inovide/models"
	"github.com/julienschmidt/httprouter"
)

var indexTemplate = template.Must(template.ParseFiles("templates/login.html", "templates/headerr.html"))

var SystemTemplates *template.Template
var (
	LENGTH_OF_IMAGE_CHARACTER        = 30
	SITE_HOST                 string = "127.0.0.1"
	SITE_PORT                        = 8080
)

// var templatesOf = template.Must(template.ParseFiles("templates/index.html", "templates/detailPage.html"))

// var client *redis.Client

// PortNumber := strconv.Itoa(SITE_PORT)
// prePath = "http://" + SITE_HOST + ":" + PortNumber ;

type UserHandler struct {
	userservice    *service.UserService
	Sessionservice *session.Cookiehandler
}

func NewUserHandler(theService *service.UserService, session *session.Cookiehandler) *UserHandler {
	return &UserHandler{userservice: theService}
}

//  The registration back ed for handling the registration of user asic User
func (user_Admin *UserHandler) RegisterUser(writer http.ResponseWriter, request *http.Request) *entity.SystemMessage {
	person := entity.Person{}
	systemmessage := &entity.SystemMessage{}
	missingdata := false
	firstname := request.FormValue("firstname")
	if firstname == "" {
		missingdata = true
	}
	lastname := request.FormValue("lastname")
	if lastname == "" {
		missingdata = true
	}
	username := strings.ToLower(request.FormValue("name"))
	if username == "" {
		missingdata = true
	}
	imagedirectory, header, erro := request.FormFile("image")
	if erro != nil {
		return nil
	}
	// if imagedirectory
	defer imagedirectory.Close()
	email := strings.ToLower(request.FormValue("email"))
	if email == "" {
		missingdata = true
	}
	password := strings.ToLower(request.FormValue("password"))
	confirmpassword := strings.ToLower(request.FormValue("confirmpassword"))
	if password == "" && confirmpassword == "" {
		missingdata = true
	}
	if password != confirmpassword {
		systemmessage.Succesful = false
		systemmessage.Errors[entity.PASSWORD_MISMATCH] = true
		if missingdata {
			systemmessage.Errors[entity.MISSING_DATA] = true
		}
		return systemmessage
	}
	if missingdata {
		systemmessage.Errors[entity.MISSING_DATA] = true
	}
	// biography := request.FormValue("biography")
	person.Email = email
	person.Username = username
	person.Password = password
	// person.Biography = biography
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
		session := &entity.Session{
			Userid:   int(person.ID),
			Username: person.Username,
		}
		// Save the Session if the User is Succesfully Registered
		user_Admin.Sessionservice.SaveSession(writer, session)
	}
	return message

}

func (user_Admin *UserHandler) TemplateRegisterUser(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	systemmessage := user_Admin.RegisterUser(writer, request)

	if systemmessage.Succesful {
		http.Redirect(writer, request, "/", 301)
	} else {
		http.Redirect(writer, request, "/user/register/", 301)
	}
}
func (user_Admin *UserHandler) ServeHome(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	// username, password, present := ReadSession(request)

	id, name, ok := user_Admin.Sessionservice.Valid(request)

	var person = &entity.Person{}
	if ok {

		person.ID = uint(id)
		person.Username = name
		// person.ID= uint(id)
		sys := user_Admin.userservice.GetUser(person)

		if !sys.Succesful {
			person = nil
		}
	} else {
		person = nil

	}
	SystemTemplates.ExecuteTemplate(writer, "Home2.html", person)
}

func (user_Admin *UserHandler) UserById(id int) *entity.Person {

	person := &entity.Person{}
	person.ID = uint(id)
	systemmessage := user_Admin.userservice.GetUserById(person)
	if systemmessage.Succesful {
		return person
	}
	return nil
}

//returns thhe page for the Registration page Template
func (user_controller *UserHandler) RegistrationPage(w http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	fmt.Println("Inside Me ")
	SystemTemplates.ExecuteTemplate(w, "reg.html", nil)
}

//CheckUser
func (user_controller *UserHandler) LogInRequest(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	request.ParseForm()

	username := request.PostFormValue("name")
	password := request.PostFormValue("password")
	fmt.Println(username, password)
	person := entity.Person{}
	person.Username = username
	// writer.Header().Add("Content-Type", "application/json")

	var thebinary []byte
	var erro error
	if username == "" || password == "" {
		message := entity.SystemMessage{}
		message.Succesful = false
		message.Message = "Please Insert The Message Appropriately"

		thebinary, erro = json.Marshal(message)

		if erro != nil {
			panic(erro)
		}

		writer.Write(thebinary)
		user_controller.View(writer, request, nil)
	}

	message := user_controller.userservice.CheckUser(&person)
	if message.Succesful {
		if person.Password == password {

			session := &entity.Session{
				Userid:   int(person.ID),
				Username: person.Username,
			}
			user_controller.Sessionservice.SaveSession(writer, session)
			// SaveSession(username, password, int(person.ID), writer, request)
		}
	} else if !message.Succesful {
		http.Redirect(writer, request, "/signin/", 301)
	}
	thebinary, erro = json.Marshal(message)
	if erro != nil {
		panic(erro)
	}

	http.Redirect(writer, request, "/user/chat/", 301)
}

func (user_controller *UserHandler) LogInPage(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	_, _, ok := user_controller.Sessionservice.Valid(request)

	if ok {
		http.Redirect(writer, request, "/user/chat/", 301)
	}
	SystemTemplates.ExecuteTemplate(writer, "login.html", nil)
}

func (user_controller *UserHandler) RedirectToHome(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	SystemTemplates.ExecuteTemplate(writer, "index.html", nil)

}
func (user_controller *UserHandler) LogOut(writer http.ResponseWriter, request *http.Request, param httprouter.Params, _ httprouter.Params) {
	user_controller.Sessionservice.DeleteSession(writer, request)
	http.Redirect(writer, request, "/", 301)
}

func (user_controller *UserHandler) View(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	id, username, present := user_controller.Sessionservice.Valid(request)
	if !present {
		//404 page no
		SystemTemplates.ExecuteTemplate(writer, "four04.html", nil)
		return
	}
	person := &entity.Person{Username: username, ID: uint(id)}
	systemMessage := user_controller.userservice.GetUser(person)
	fmt.Println(systemMessage.Message)
	if systemMessage.Succesful {
		fmt.Println(person.Email)
		SystemTemplates.ExecuteTemplate(writer, "edit.html", person)
	}
}
