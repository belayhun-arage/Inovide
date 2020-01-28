package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	AdminRepository "github.com/Projects/Inovide/Admin/Repository"

	session "github.com/Projects/Inovide/Session"
	UsableFunctions "github.com/Projects/Inovide/Usables"
	entity "github.com/Projects/Inovide/models"
	"github.com/julienschmidt/httprouter"
	// "github.com/Projects/Inovide/"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
	AdminService *AdminRepository.AdminRepo
	Adminrepo    *AdminRepository.AdminRepo
	Session      *session.Cookiehandler
	UserHandler  *UserHandler
}

func NewAdminHandler(adminrepo *AdminRepository.AdminRepo, userhandler *UserHandler, session *session.Cookiehandler) *AdminHandler {
	return &AdminHandler{Adminrepo: adminrepo, UserHandler: userhandler, Session: session}
}

func (user_Admin *AdminHandler) CreateAdmin(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	person := entity.Person{}
	systemmessage := &entity.SystemMessage{
		Errors:    make(map[string]bool),
		Succesful: false,
		Message:   "Invalid Username Or Email",
	}
	writer.Header().Add("Content-Type", "application/json")
	jsonsystemmessage, _ := json.Marshal(systemmessage)
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
		writer.Write(jsonsystemmessage)
	}
	if missingdata {
		systemmessage.Errors[entity.MISSING_DATA] = true
	}
	// biography := request.FormValue("biography")
	person.Email = email
	person.Username = username
	isAdmin, _ := user_Admin.Session.Authorize(request)

	if isAdmin {
		person.IsAdmin = true
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		systemmessage.Succesful = false
		writer.Write(jsonsystemmessage)

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

	systemmessage = user_Admin.UserHandler.userservice.RegisterUser(&person)

	fmt.Println(systemmessage.Succesful)
	fmt.Println(person.ID)
	if systemmessage.Succesful {
		fmt.Println(person.Username, person.Lastname, "<------------------------")
		systemmessage.Message = "Succesfully Registered "
		systemmessage.Succesful = true
		if erro == nil {
			io.Copy(file, imagedirectory)
		}
	}
	writer.Write(jsonsystemmessage)

}

func (adminhandler *AdminHandler) AdminDeleteIdea(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	systemmessage := &entity.SystemMessage{}
	systemmessage.Message = "Not Succesful "
	systemmessage.Succesful = false
	jsonsystemmessage, _ := json.Marshal(systemmessage)
	writer.Header().Add("Content-Type", "application/json")
	isAdmin, isUser := adminhandler.Session.Authorize(request)
	if !isAdmin || !isUser {
		systemmessage.Message = "Not Authorized "
		systemmessage.Succesful = false
		jsonsystemmessage, _ = json.Marshal(systemmessage)
		writer.Write(jsonsystemmessage)

	}
	ideaid, err := strconv.Atoi(request.FormValue("ideaid"))
	if err != nil {
		writer.Write(jsonsystemmessage)
	}

	systemmessage = adminhandler.UserHandler.userservice.AdminDeleteIdea(&entity.Idea{Id: int(ideaid)})
	if systemmessage.Succesful {
		systemmessage.Message = "Succesfully Deleted The Idea"
		jsonsystemmessage, _ = json.Marshal(systemmessage)
		writer.Write(jsonsystemmessage)

	} else {
		writer.Write(jsonsystemmessage)
	}
}

func (adminhandler *AdminHandler) AdminDeleteUser(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	systemmessage := &entity.SystemMessage{}
	systemmessage.Message = "Not Succesful "
	systemmessage.Succesful = false
	jsonsystemmessage, _ := json.Marshal(systemmessage)
	writer.Header().Add("Content-Type", "application/json")
	isAdmin, isUser := adminhandler.Session.Authorize(request)
	if !isAdmin || !isUser {
		systemmessage.Message = "Not Authorized "
		systemmessage.Succesful = false
		jsonsystemmessage, _ = json.Marshal(systemmessage)
		writer.Write(jsonsystemmessage)

	}
	userid, err := strconv.Atoi(request.FormValue("userid"))
	if err != nil {
		writer.Write(jsonsystemmessage)
	}

	systemmessage = adminhandler.UserHandler.userservice.AdminDeleteuser(&entity.Person{ID: uint(userid)})
	if systemmessage.Succesful {
		systemmessage.Message = "Succesfully Deleted The User"
		jsonsystemmessage, _ = json.Marshal(systemmessage)
		writer.Write(jsonsystemmessage)
	} else {
		writer.Write(jsonsystemmessage)
	}
}

func (adminhandler *AdminHandler) AnalysisPage(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	analysis := map[string]int64{}
	_, _ = adminhandler.Session.Authorize(request)
	// if !isuser || !isAdmin {
	// 	SystemTemplates.ExecuteTemplate(writer, "four04.html", nil)
	// 	return
	// }
	usrs := adminhandler.Adminrepo.CountUsers()
	analysis["users"] = usrs
	ideas := adminhandler.Adminrepo.CountIdeas()
	analysis["ideas"] = ideas
	admins := adminhandler.Adminrepo.CountAdmins()
	analysis["admins"] = admins
	messages := adminhandler.Adminrepo.CountMessages()
	analysis["messages"] = messages
	activeusers := adminhandler.Adminrepo.CountActiveUsers()
	analysis["activeusers"] = activeusers
	SystemTemplates.ExecuteTemplate(writer, "dashboard.html", analysis)
}

func (adminhandler *AdminHandler) TemplateAdminUser(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	SystemTemplates.ExecuteTemplate(writer, "adminIdea.html", nil)
}

func (user_Admin *AdminHandler) TemplateCreateAdmin(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	SystemTemplates.ExecuteTemplate(writer, "adminuser.html", nil)

}
