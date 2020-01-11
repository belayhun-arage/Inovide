package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	ideaService "github.com/Projects/Inovide/Idea/Service"
	UsableFunctions "github.com/Projects/Inovide/Usables"
	entity "github.com/Projects/Inovide/models"

	"github.com/lib/pq"
)

var (
	LENGTH_OF_FILE_CHARACTER = 30
)

type IdeaHandler struct {
	ideaservice    *ideaService.IdeaService
	commenthandler *CommentHandler
	userrouter     *UserHandler
}

func NewIdeaHandler(theService *ideaService.IdeaService, commenthandle *CommentHandler, userrouters *UserHandler) *IdeaHandler {
	return &IdeaHandler{ideaservice: theService, commenthandler: commenthandle, userrouter: userrouters}
}
func (idea_controller *IdeaHandler) CreateIdeaPage(writer http.ResponseWriter, request *http.Request) {
	SystemTemplates.ExecuteTemplate(writer, "createIdea.html", nil)
}

//CreateIdea handler
func (idea_Admin *IdeaHandler) CreateIdea(writer http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(3939393939393933939)
	idea := entity.Idea{}

	ideaTitle := request.FormValue("title")
	description := request.FormValue("description")
	visibiitty := request.FormValue("visibility")
	idea.Ideaownerid = 0
	idea.Title = ideaTitle
	idea.Description = description
	idea.Visibility = visibiitty

	resourceArray := []string{}
	files := make([]multipart.File, 3)

	request.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
	fhs := request.MultipartForm.File["files"]
	counter := 0
	for _, fh := range fhs {

		filedirectory, erro := fh.Open()
		if erro != nil {
			//fmt.Println(erro)
			break
		}
		fmt.Println("Looping ")
		defer filedirectory.Close()

		var newFullNameOfTheFileDirectory string
		// var file *os.File
		if fh.Filename != "" {

			fmt.Println("Fetching Dir")
			stringSliceOfNameOfFile := strings.Split(fh.Filename, ".")

			fileExtension := stringSliceOfNameOfFile[len(stringSliceOfNameOfFile)-1]

			randomStringForSavingTheFile := UsableFunctions.GenerateRandomString(LENGTH_OF_FILE_CHARACTER, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

			newFullNameOfTheFileDirectory = "public/img/IdeaResource/" + randomStringForSavingTheFile + "." + fileExtension

			fil, err := os.Create(newFullNameOfTheFileDirectory)
			if err != nil {
				fmt.Println("CAn't Save the Image ")
				continue
			}
			defer fil.Close()
			io.Copy(fil, filedirectory)

			resourceArray[counter] = newFullNameOfTheFileDirectory

			counter++
		}

	}

	fmt.Println(resourceArray)
	idea.Resources = pq.StringArray(resourceArray)
	message := idea_Admin.ideaservice.CreateIdea(&idea)
	if message.Succesful {
		for index, filemultipart := range files {
			filecreate, err := os.Create(resourceArray[index])
			if err != nil {
				continue
			}

			fmt.Println("Created")
			io.Copy(filecreate, filemultipart)
			filecreate.Close()

		}
	}

	writer.Write([]byte(message.Message))
}

func (idea_Admin *IdeaHandler) GetIdea(writer http.ResponseWriter, request *http.Request) *entity.Idea {
	request.ParseForm()

	Id := request.FormValue("ideadId")
	id, err := strconv.Atoi(Id)
	if err != nil {
		return nil
	}

	idea := &entity.Idea{}
	idea, systemMessage := idea_Admin.ideaservice.GetIdea(idea, id)
	fmt.Println(idea.Description, idea.Id)
	if systemMessage.Succesful {
		return idea
	}
	return nil
}
func (idea_Admin *IdeaHandler) TemplateGetIdea(writer http.ResponseWriter, request *http.Request) {

	theidea := idea_Admin.GetIdea(writer, request)
	if theidea == nil {
		SystemTemplates.ExecuteTemplate(writer, "four04.html", nil)
	}
	SystemTemplates.ExecuteTemplate(writer, "", theidea)
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

func (idea_Admin *IdeaHandler) GetDetailIdea(writer http.ResponseWriter, request *http.Request) *entity.IdeaPersonComments {
	ideapersoncomment := &entity.IdeaPersonComments{}
	idea := idea_Admin.GetIdea(writer, request)
	comments, ok := idea_Admin.commenthandler.GetComments(idea.Id)
	user := idea_Admin.userrouter.UserById(idea.Ideaownerid)

	if !ok || idea == nil || comments == nil || user == nil {
		ideapersoncomment.Succesful = false
		return ideapersoncomment
	}
	ideapersoncomment.Succesful = true

	commentwithPersons := idea_Admin.commenthandler.GetCommentWithPerson(comments)

	if commentwithPersons == nil {
		return nil
	}

	ideapersoncomment.Succesful = true
	ideapersoncomment.CommentAndPerson = *commentwithPersons
	ideapersoncomment.Idea = *idea
	return ideapersoncomment

}

func (idea_Admin *IdeaHandler) TemplateGetDetailIdea(writer http.ResponseWriter, request *http.Request) {

	listed := idea_Admin.GetDetailIdea(writer, request)

	SystemTemplates.ExecuteTemplate(writer, "", listed)
}

// func (idea_Admin *IdeaHandler) SaveComment(writer http.ResponseWriter, request *http.Request) {
// 	request.ParseForm()
// 	fmt.Println("Inside The Handler ")
// 	ideaid, err := strconv.Atoi(request.FormValue("ideaid"))

// 	if err != nil {
// 		return
// 	}
// 	commentorid, err := strconv.Atoi(request.FormValue("commentorid"))
// 	if err != nil {
// 		return
// 	}
// 	commentdata := request.FormValue("commentdata")

// 	if commentdata == "" {
// 		return
// 	}
// comment := &entity.Comment{Ideaid: ideaid, Commentorid: commentorid, Commentdata: commentdata}
// //systemmessage := idea_Admin.ideaservice.SaveCommentIdea(comment)

// mapping := map[string]interface{}{"systemmessage": systemmessage,
// 	"comment": comment}
// jsonbinary, err := json.Marshal(mapping) //  can be sent for the user through  Template(The Map) or api (The Json )
// if err != nil {
// 	return
// }
// writer.Write(jsonbinary)
// }
