package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	UsableFunctions "github.com/Projects/Inovide/Usables"
	service "github.com/Projects/Inovide/User/Service"
	entity "github.com/Projects/Inovide/models"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

var indexTemplate = template.Must(template.ParseFiles("templates/login.html", "templates/headerr.html"))

var SystemTemplates *template.Template

/*                Session Related Datas                                 <<Begin>>     */
/*                Session Related Datas                                 <<Begin>>     */
/*                Session Related Datas                                 <<Begin>>     */
/*                Session Related Datas                                 <<Begin>>     */

const (
	SESSION_USER_NAME = "username"
	SESSION_PASSWORD  = "password"
	SESSION_ID        = "id"
)

var store = sessions.NewCookieStore([]byte("The Top Secter ot the daya ")) // The place where to save the session Cookies on

func SaveSession(username string, password string, id int, writer http.ResponseWriter, request *http.Request) {

	session, erro := store.Get(request, "session")
	if erro != nil {
		fmt.Println("Error While Reading the session  ")
		return
	}
	session.Values[SESSION_USER_NAME] = username // session.Values is the Map Insid the sessions Package to Hold Messages
	session.Values[SESSION_PASSWORD] = password  // session.Values is the Map Insid the sessions Package to Hold Messages
	session.Values[SESSION_ID] = id
	session.Save(request, writer) // writing the session to the ResposnseWriter
}

//  --  Use This to read The Session From The Sessio Store
func ReadSession(request *http.Request) (string, string, bool) {

	sessional, err := store.Get(request, "session")
	if err != nil {
		return "", "", false
	}

	usernam, ok := (sessional.Values[SESSION_USER_NAME])
	if !ok {
		return "", "", ok
	}
	username := usernam.(string)
	password := sessional.Values[SESSION_PASSWORD].(string)

	// id, err := strconv.Atoi(sessional.Values[SESSION_ID].(string))

	if err != nil {
		return "", "", false
	}
	return username, password, ok
}

func DeleteSession(request *http.Request) bool {
	sessional, err := store.Get(request, "session")
	if err != nil {
		return false
	}
	sessional.Options.MaxAge = -1
	return true
}

/*                Session Related Datas                                 <<Begin>>     */
/*                Session Related Datas                                 <<Begin>>     */

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
	userservice *service.UserService
}

func NewUserHandler(theService *service.UserService) *UserHandler {
	return &UserHandler{userservice: theService}
}

func (user_Admin *UserHandler) RegisterUser(writer http.ResponseWriter, request *http.Request) *entity.SystemMessage {
	person := entity.Person{}

	firstname := request.FormValue("firstname")
	lastname := request.FormValue("lastname")
	username := request.FormValue("name")
	imagedirectory, header, erro := request.FormFile("image")
	if erro != nil {
		return nil
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
		SaveSession(username, password, int(person.ID), writer, request)

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

	var person = &entity.Person{}
	// if present {

	// 	person.Username = username
	// 	person.Password = password
	// 	// person.ID= uint(id)
	// 	sys := user_Admin.userservice.GetUser(person)

	// 	if !sys.Succesful {
	// 		person = nil
	// 	}
	// } else {
	// 	person = nil
	// }
	fmt.Println("Inside Me ...")

	indexTemplate.ExecuteTemplate(writer, "index.html", person)
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
	person.Password = password
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
		SaveSession(username, password, int(person.ID), writer, request)
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
	_, _, present := ReadSession(request)
	if present {
		http.Redirect(writer, request, "/user/chat/", 301)
	}
	SystemTemplates.ExecuteTemplate(writer, "login.html", nil)
}

func (user_controller *UserHandler) RedirectToHome(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	SystemTemplates.ExecuteTemplate(writer, "index.html", nil)

}
func (user_controller *UserHandler) LogOut(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	_ = DeleteSession(request)
	http.Redirect(writer, request, "/", 301)
}

func (user_controller *UserHandler) View(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	username, password, present := ReadSession(request)
	if !present {
		//404 page no
		SystemTemplates.ExecuteTemplate(writer, "four04.html", nil)
		return
	}
	person := &entity.Person{Username: username, Password: password}
	systemMessage := user_controller.userservice.GetUser(person)
	fmt.Println(systemMessage.Message)
	if systemMessage.Succesful {
		fmt.Println(person.Email)
		SystemTemplates.ExecuteTemplate(writer, "edit.html", person)
	}
}
