package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	ideaService "github.com/Samuael/Projects/Inovide/Idea/Service"
	UsableFunctions "github.com/Samuael/Projects/Inovide/Usables"
	entity "github.com/Samuael/Projects/Inovide/models"
)

var (
	LENGTH_OF_FILE_CHARACTER = 30
)

type IdeaHandler struct {
	ideaservice *ideaService.IdeaService
}

func NewIdeaHandler(theService *ideaService.IdeaService) *IdeaHandler {
	return &IdeaHandler{ideaservice: theService}
}

func (idea_controller *IdeaHandler) CreateIdeaPage(writer http.ResponseWriter, request *http.Request) {
	SystemTemplates.ExecuteTemplate(writer, "createIdea.html", nil)
}

//CreateIdea handler
func (idea_Admin *IdeaHandler) CreateIdea(writer http.ResponseWriter, request *http.Request) {
	idea := entity.Idea{}

	ideaTitle := request.FormValue("title")
	description := request.FormValue("description")
	filedirectory, header, erro := request.FormFile("filename")
	visibiitty := request.FormValue("visibility")

	if erro != nil {
		//fmt.Println(erro)
	}
	defer filedirectory.Close()

	idea.Ideaownerid = 0
	idea.Title = ideaTitle
	idea.Description = description
	idea.Visibility = visibiitty
	var newFullNameOfTheFileDirectory string
	var file *os.File
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
	message := idea_Admin.ideaservice.CreateIdea(&idea)
	if message.Succesful {
		io.Copy(file, filedirectory)
		//SaveSession(username, password, writer, request)
	}

	writer.Write([]byte(message.Message))
}

func (idea_Admin *IdeaHandler) GetIdea(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	Id := request.FormValue("ideadId")
	id, err := strconv.Atoi(Id)
	if err != nil {
		return
	}

	idea := &entity.Idea{}
	idea, systemMessage := idea_Admin.ideaservice.GetIdea(idea, id)
	fmt.Println(idea.Description, idea.Id)
	if systemMessage.Succesful {
		json, _ := json.Marshal(idea)
		fmt.Println(string(json))
		writer.Header().Add("Content-type", "application/json")
		writer.Write(json)
	}
}

func (idea_Admin *IdeaHandler) DeleteIdea(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	Id := request.FormValue("ideadId")
	id, err := strconv.Atoi(Id)
	if err != nil {
		return
	}
	systemmessage := idea_Admin.ideaservice.DeleteIdea(id)

	jsonMessage, _ := json.Marshal(systemmessage)

	writer.Write(jsonMessage)
}

func (idea_Admin *IdeaHandler) UpdateIdea(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()

	id := request.FormValue("id")
	title := request.FormValue("title")
	description := request.FormValue("description")
	// file, header, err := request.FormFile("file")
	visibility := request.FormValue("visibility")

	var mapps map[string]string
	if id != "" {
		mapps["id"] = id
	}
	if title != "" {

		mapps["title"] = title
	}
	if description != "" {

		mapps["description"] = title
	}
	// if file != "" {

	// 	mapps["resources"] = title
	// }
	if visibility != "" {
		mapps["visibility"] = visibility
	}
}

func (idea_Admin *IdeaHandler) VoteIdea(writer http.ResponseWriter, request *http.Request) {
	ideaid, err := strconv.Atoi(request.FormValue("id"))
	if err != nil {
		return
	}
	voterid, err := strconv.Atoi(request.FormValue("voterid"))
	if err != nil {
		return
	}
	systemmessage := idea_Admin.ideaservice.VoteIdea(ideaid, voterid)

	jsonbinary, err := json.Marshal(systemmessage)

	if err != nil {
		return
	}
	writer.Write(jsonbinary)
}
