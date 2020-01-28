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
	"golang.org/x/crypto/bcrypt"
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
	RecentUser     *entity.Person
}

func NewUserHandler(theService *service.UserService, session *session.Cookiehandler) *UserHandler {
	return &UserHandler{userservice: theService, Sessionservice: session}
}

//  The registration back ed for handling the registration of user asic User
func (user_Admin *UserHandler) RegisterUser(writer http.ResponseWriter, request *http.Request) *entity.SystemMessage {
	person := entity.Person{}
	systemmessage := &entity.SystemMessage{
		Errors:    make(map[string]bool),
		Succesful: false,
		Message:   "Invalid Username Or Email",
	}
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

	/// the input imagge from users registration form
	imagedirectory, header, erro := request.FormFile("image")
	if erro == nil {
		defer imagedirectory.Close()

	}
	// if imagedirectory
	email := strings.ToLower(request.FormValue("email"))
	if email == "" {
		missingdata = true
	}
	password := strings.ToLower(request.FormValue("password"))
	confirmpassword := strings.ToLower(request.FormValue("confirmpassword"))
	if password == "" && confirmpassword == "" {
		missingdata = true
	}
	if strings.Compare(password, confirmpassword) == 0 {
		systemmessage.Succesful = false
		systemmessage.Message = "Invalid Password!\nConfirm Correctly"
		systemmessage.Errors[entity.PASSWORD_MISMATCH] = true
		if missingdata {
			systemmessage.Errors[entity.MISSING_DATA] = true
			systemmessage.Message = "Fill The Datas Correctly"
		}
		return systemmessage
	}
	if missingdata {
		systemmessage.Errors[entity.MISSING_DATA] = true
	}
	// biography := request.FormValue("biography")
	person.Email = email
	person.Username = username
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		systemmessage.Succesful = false
		return systemmessage
	}
	person.Password = string(hashedPassword)
	// person.Biography = biography
	person.Firstname = firstname
	person.Lastname = lastname
	var newFullNameOfTheImageDirectory string
	var file *os.File
	if erro == nil {
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
	}

	systemmessage = user_Admin.userservice.RegisterUser(&person)

	fmt.Println(systemmessage.Succesful)
	fmt.Println(person.ID)
	if systemmessage.Succesful {
		fmt.Println(person.Username, person.Lastname, "<------------------------")
		systemmessage.Message = "Succesfully Registered "
		systemmessage.Succesful = true
		if erro == nil {
			io.Copy(file, imagedirectory)
		}
		fmt.Println(person.ID)
		session := &entity.Session{
			Userid:   int(person.ID),
			Username: person.Username,
		}
		// Save the Session if the User is Succesfully Registered
		user_Admin.Sessionservice.DeleteSession(writer, request)
		user_Admin.Sessionservice.SaveSession(writer, session)
	}
	return systemmessage
}
func (user_Admi *UserHandler) TemplateRegistrationRequest(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	systemmessage := user_Admi.RegisterUser(writer, request)

	if systemmessage.Succesful {
		http.Redirect(writer, request, "/", 301)
	} else {
		SystemTemplates.ExecuteTemplate(writer, "reg.html", systemmessage)
	}
}
func (user_Admin *UserHandler) TemplateRegisterUser(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	systemmessage := user_Admin.RegisterUser(writer, request)

	if systemmessage.Succesful {

		http.Redirect(writer, request, "/", 301)
	} else {
		// http.Redirect(writer, request, "/user/register/", 301)
		SystemTemplates.ExecuteTemplate(writer, "reg.html", systemmessage)
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

		if user_Admin.RecentUser != nil {
			person = user_Admin.RecentUser
		} else {
			sys := user_Admin.userservice.GetUser(person)
			if !sys.Succesful {
				person = nil
			}
		}
	} else {
		person = nil

	}

	/*The Templating work TO be Done to Render the User information */
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

var message *entity.SystemMessage
var person entity.Person

//CheckUser
func (user_controller *UserHandler) LogInRequest(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) (*entity.SystemMessage, entity.Person) {
	request.ParseForm()
	username := strings.ToLower(request.PostFormValue("name"))
	password := strings.ToLower(request.PostFormValue("password"))
	fmt.Println(username, password)
	person = entity.Person{}
	person.Username = username
	// writer.Header().Add("Content-Type", "application/json")
	// var thebinary []byte
	// var erro error
	if username == "" || password == "" {
		message := entity.SystemMessage{}
		message.Succesful = false
		message.Message = "Please Insert The Message Appropriately"
		return &message, person
	}

	message = user_controller.userservice.CheckUser(&person)
	if message.Succesful {

		fmt.Println(person.Password)
		// Comparing the hashed value of the user password form the Database
		err := bcrypt.CompareHashAndPassword([]byte(person.Password), []byte(password))
		// Comparing the hashed value of the user password form the Database
		if err == bcrypt.ErrMismatchedHashAndPassword {
			message.Succesful = false
			return message, person
		}
		if err == nil {
			fmt.Println("\n\n\n\n\n\nThe Person DOes Exist\n\n\n\n\n\n\n\n ")
			session := &entity.Session{
				Userid:   int(person.ID),
				Username: person.Username,
			}
			user_controller.Sessionservice.DeleteSession(writer, request)
			user_controller.Sessionservice.SaveSession(writer, session)
			// SaveSession(username, password, int(person.ID), writer, request)
		}
		// } else if !message.Succesful {
		// 	http.Redirect(writer, request, "/", 301)
		// }
		// thebinary, erro = json.Marshal(message)
		// if erro != nil {
		// 	panic(erro)
		// }

		// http.Redirect(writer, request, "/user/chat/", 301)
	} else {
		message.Message = "Invalid Username/Password  "
		message.Succesful = false
	}
	return message, person
}

func (userhandler *UserHandler) TemplateLogInPage(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	systemmessage, person := userhandler.LogInRequest(writer, request, nil)

	fmt.Println(systemmessage.Succesful, systemmessage.Message)
	if systemmessage.Succesful {
		userhandler.RecentUser = &person
		http.Redirect(writer, request, "/", 301)
	} else {
		http.Redirect(writer, request, "/user/signin/", 301)
	}
}

func (user_controller *UserHandler) LogInPage(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	_, _, ok := user_controller.Sessionservice.Valid(request)

	if ok {
		http.Redirect(writer, request, "/user/chat/", 301)
	}
	SystemTemplates.ExecuteTemplate(writer, "login.html", nil)
}

// func (user_controller *UserHandler) RedirectToHome(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
// 	SystemTemplates.ExecuteTemplate(writer, "index.html", nil)

// }
func (user_controller *UserHandler) TemplateLogOut(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	fmt.Println("   Log Out Method Called ")
	success := user_controller.Sessionservice.DeleteSession(writer, request)
	systemmessage := &entity.SystemMessage{}
	if !success {
		fmt.Println("Can't Log Out ")

		systemmessage.Succesful = false
		systemmessage.Message = "Can't Log Out "

	} else {
		systemmessage.Succesful = true
		systemmessage.Message = "Succesfully Logged Out "

	}
	jsonmessagelogout, _ := json.Marshal(systemmessage)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(jsonmessagelogout)
}

func (user_controller *UserHandler) ViewProfile(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	id, username, _ := user_controller.Sessionservice.Valid(request)
	if id <= 0 {
		//404 page no
		SystemTemplates.ExecuteTemplate(writer, "four04.html", nil)
		return
	}
	fmt.Println(id)
	person := &entity.Person{Username: username, ID: uint(id)}
	systemMessage := user_controller.userservice.GetUser(person)
	fmt.Println(systemMessage.Message)
	if systemMessage.Succesful {
		fmt.Println(person.Email)
		SystemTemplates.ExecuteTemplate(writer, "edit.html", person)
		fmt.Println(person.ID, person.Firstname)
	}
}
func (user_controller *UserHandler) EditProfile(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) (entity.Person, *entity.SystemMessage) {
	person := entity.Person{}

	id, username, _ := user_controller.Sessionservice.Valid(request)
	if id > 0 {
		person.ID = uint(id)
		user_controller.userservice.CheckUser(&person)

	} else {
		http.Redirect(writer, request, "/", 301)

	}
	systemmessage := &entity.SystemMessage{
		Errors:    make(map[string]bool),
		Succesful: false,
		Message:   "Invalid Username Or Email",
	}
	firstname := request.FormValue("firstname")
	if firstname != "" {
		person.Firstname = firstname
	}
	lastname := request.FormValue("lastname")
	if lastname != "" {
		person.Lastname = lastname
	}
	newUsername := strings.ToLower(request.FormValue("name"))
	if username != "" {
		person.Username = newUsername
	}
	/// the input imagge from users registration form
	imagedirectory, header, erro := request.FormFile("image")
	if erro == nil {
		defer imagedirectory.Close()
	}
	// if imagedirectory
	email := strings.ToLower(request.FormValue("email"))
	if email != "" {
		person.Email = email
	}
	biography := request.FormValue("biography")
	if biography != "" {
		person.Biography = biography
	}
	var newFullNameOfTheImageDirectory string
	var file *os.File
	if erro == nil {
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
	}
	systemmessage = user_controller.userservice.UpdateUser(&person)
	if systemmessage.Succesful {
		systemmessage.Message = "Succesfully Registered "
		systemmessage.Succesful = true
		if erro == nil {
			io.Copy(file, imagedirectory)
		}
		session := &entity.Session{
			Userid:   int(person.ID),
			Username: person.Username,
		}
		// Save the Session if the User is Succesfully Registered
		user_controller.Sessionservice.SaveSession(writer, session)

		person.Password = ""
		return person, systemmessage
	}

	recover()
	systemmessage.Succesful = false
	systemmessage.Message = "Can't update !"
	return person, systemmessage
}
