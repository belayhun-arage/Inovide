package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	UsableFunctions "github.com/BELAY-hun/Projects/Inovide/Usables"
	service "github.com/BELAY-hun/Projects/Inovide/User/Service"
	entity "github.com/BELAY-hun/Projects/Inovide/models"
)

var (
	LENGTH_OF_FILE_CHARACTER = 30
)

var TemplateIdea = template.Must(template.ParseFiles("templates/createidea.html"))

type IdeaHandler struct {
	ideaservice *service.IdeaService
}

func NewIdeaHandler(theService *service.IdeaService) *IdeaHandler {
	return &IdeaHandler{ideaservice: theService}
}

func (idea_controller *IdeaHandler) CreateIdeaPage(writer http.ResponseWriter, request *http.Request) {
	TemplateIdea.Execute(writer, nil)
}

//CreateIdea handler
func (idea_controller *IdeaHandler) CreateIdea(writer http.ResponseWriter, request *http.Request) {
	idea := entity.Idea{}

	ideaTitle := request.FormValue("title")
	description := request.FormValue("description")
	filedirectory, header, erro := request.FormFile("filename")

	if erro != nil {
		//fmt.Println(erro)
	}
	defer filedirectory.Close()

	idea.Title = ideaTitle
	idea.Description = description
	var newFullNameOfTheFileDirectory string
	//var file *os.File
	if header.Filename != "" {

		stringSliceOfNameOfFile := strings.Split(header.Filename, ".")

		fileExtension := stringSliceOfNameOfFile[len(stringSliceOfNameOfFile)-1]

		randomStringForSavingTheFile := UsableFunctions.GenerateRandomString(LENGTH_OF_FILE_CHARACTER, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

		newFullNameOfTheFileDirectory = "public/img/UsersImage/" + randomStringForSavingTheFile + "." + fileExtension

		file, errorCreatingFile := os.Create(newFullNameOfTheFileDirectory)

		if errorCreatingFile != nil {
			fmt.Println("Error While Creating the Image ", errorCreatingFile)
		}
		defer file.Close()
	}
}
