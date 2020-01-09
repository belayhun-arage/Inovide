package handler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	ideaService "github.com/Projects/Inovide/Idea/Service"
	UsableFunctions "github.com//Projects/Inovide/Usables"
	entity "github.com/Projects/Inovide/models"
)

var (
	LENGTH_OF_FILE_CHARACTER = 30
)

var TemplateIdea = template.Must(template.ParseFiles("templates/createidea.html"))

type IdeaHandler struct {
	ideaservice *ideaService.IdeaService
}

func NewIdeaHandler(theService *ideaService.IdeaService) *IdeaHandler {
	return &IdeaHandler{ideaservice: theService}
}

func (idea_controller *IdeaHandler) CreateIdeaGetHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Creatinrg idea")
	errr := TemplateIdea.Execute(writer, nil)
	if errr != nil {
		fmt.Println("Error with  executing createidea.html file")
		panic(errr)
	}

}

//CreateIdea handler
func (idea_Admin *IdeaHandler) CreateIdeaPostHandler(writer http.ResponseWriter, r *http.Request) {
	idea := entity.Idea{}

	ideaTitle := r.FormValue("title")
	description := r.FormValue("description")
	filedirectory, header, erro := r.FormFile("filename")

	if erro != nil {
		fmt.Println(erro)
	}
	defer filedirectory.Close()

	idea.Title = ideaTitle
	idea.Description = description
	var newFullNameOfTheFileDirectory string
	var file *os.File
	if header.Filename != "" {

		stringSliceOfNameOfFile := strings.Split(header.Filename, ".")

		fileExtension := stringSliceOfNameOfFile[len(stringSliceOfNameOfFile)-1]

		randomStringForSavingTheFile := UsableFunctions.GenerateRandomString(LENGTH_OF_FILE_CHARACTER, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

		newFullNameOfTheFileDirectory = "public/img/UsersImage/" + randomStringForSavingTheFile + "." + fileExtension

		file, errorCreatingFile := os.Create(newFullNameOfTheFileDirectory)

		if errorCreatingFile != nil {
			fmt.Println("Error While Creating the File directory ", errorCreatingFile)
		}
		defer file.Close()
	}
	message := idea_Admin.ideaservice.CreateIdea(&idea)
	if message.Succesful {
		io.Copy(file, filedirectory)
		//SaveSession(username, password, writer, request)
	}

	writer.Write([]byte(message.Message))
}
