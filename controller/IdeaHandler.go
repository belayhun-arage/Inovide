package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	CommentRepo "github.com/Projects/Inovide/Comment/Repository"
	ideaService "github.com/Projects/Inovide/Idea/Service"
	session "github.com/Projects/Inovide/Session"
	UsableFunctions "github.com/Projects/Inovide/Usables"
	entity "github.com/Projects/Inovide/models"
	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
)

var (
	LENGTH_OF_FILE_CHARACTER = 30
)

type IdeaHandler struct {
	ideaservice    *ideaService.IdeaService
	commenthandler *CommentHandler
	userrouter     *UserHandler
	Session        *session.Cookiehandler
	commentrepo    *CommentRepo.CommentRepo
}

func NewIdeaHandler(theService *ideaService.IdeaService,
	commenthandle *CommentHandler, userrouters *UserHandler,
	sessin *session.Cookiehandler) *IdeaHandler {
	return &IdeaHandler{ideaservice: theService,
		commenthandler: commenthandle,
		userrouter:     userrouters}
}
func (ideahandle *IdeaHandler) SetCommentRepo(commentrepo *CommentRepo.CommentRepo) {
	commentrepo = commentrepo
}

func (idea_controller *IdeaHandler) CreateIdeaPage(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	SystemTemplates.ExecuteTemplate(writer, "createIdeaMod.html", nil)
}

func (idea_controller *IdeaHandler) TemplateCreateIdea(writer http.ResponseWriter, request *http.Request, param httprouter.Params) {
	systemmessage, _ := idea_controller.CreateIdea(writer, request, nil)
	if systemmessage.Succesful {
		systemmessage.Succesful = true
		systemmessage.Message = "Idea Created Succesfully"
	}

	// for test only
	http.Redirect(writer, request, "/idea/myideas/", 301)

}

//CreateIdea handler
func (idea_Admin *IdeaHandler) CreateIdea(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) (*entity.SystemMessage, *entity.Idea) {
	request.ParseMultipartForm(3939393939393933939)
	idea := &entity.Idea{}
	ideaTitle := request.FormValue("title")
	description := request.FormValue("description")
	visibiitty := request.FormValue("visibility")
	id, _, _ := idea_Admin.Session.Valid(request)
	if id <= 0 {
		return nil, nil
	}
	fmt.Println(id)
	idea.Ideaownerid = id
	idea.Title = ideaTitle
	idea.Description = description
	if visibiitty == "" {
		visibiitty = "pu"
	}
	idea.Visibility = visibiitty
	resourceArray := [5]string{}
	//files := make([]multipart.File, 3)

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
			fmt.Println(resourceArray)
			fmt.Println(newFullNameOfTheFileDirectory)
			counter++
		}
	}
	sliceofResourcesName := resourceArray[0:]
	fmt.Println(sliceofResourcesName[1], sliceofResourcesName[0], sliceofResourcesName[3], sliceofResourcesName[2])
	idea.Resources = pq.StringArray(sliceofResourcesName)
	message := idea_Admin.ideaservice.CreateIdea(idea)
	if message.Succesful {
		for index, filemultipart := range fhs {
			file, err := filemultipart.Open()
			if err == nil {
				newfileCreated, err := os.Create(resourceArray[index])
				if err != nil {
					continue
				}

				fmt.Println("Created")
				io.Copy(newfileCreated, file)

				newfileCreated.Close()
				file.Close()
			} else {
				message.Message = "Error While Creating File "
				message.Succesful = false
				return message, idea

			}

		}
		message.Message = "Idea Succesfuly Created "
		message.Succesful = true
		return message, idea
	}
	message.Message = "Idea Not Created "
	message.Succesful = false
	return message, idea
}

func (idea_Admin *IdeaHandler) GetIdea(writer http.ResponseWriter, request *http.Request) *entity.Idea {
	request.ParseForm()
	idea := &entity.Idea{}

	id, err := strconv.Atoi(request.FormValue("ideaid"))
	if err != nil {
		fmt.Println("Erorr inside GetIdea Application")
		return nil
	}
	idea.Id = id
	systemMessage := idea_Admin.ideaservice.GetIdea(idea, id)
	fmt.Println(idea.Description, idea.Id)
	if systemMessage.Succesful {
		return idea
	}
	return nil
}
func (ideahandler *IdeaHandler) ApiGetIdea(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	ideanull := &entity.Idea{Title: "InValid Idea", Description: "This Idea Is Sent "}
	systemmessage := &entity.SystemMessage{}
	idea := ideahandler.GetIdea(writer, request)

	messageAndIdea := map[string]interface{}{}
	if idea == nil {
		messageAndIdea["Idea"] = ideanull
		systemmessage.Message = "Can't Load the Idea In The Specified Id"
		systemmessage.Succesful = false
		messageAndIdea["Message"] = systemmessage
	} else {
		systemmessage.Message = "Here Is The Idea "
		systemmessage.Succesful = true
		messageAndIdea["Message"] = systemmessage
		messageAndIdea["Idea"] = idea
	}
	jsonidea, _ := json.Marshal(messageAndIdea)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(jsonidea)
}

func (idea_Admin *IdeaHandler) TemplateGetIdea(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	theidea := idea_Admin.GetIdea(writer, request)
	if theidea == nil {
		SystemTemplates.ExecuteTemplate(writer, "four04.html", nil)
	}
	SystemTemplates.ExecuteTemplate(writer, "viewidea.layout", theidea)
}
func (idea_Admin *IdeaHandler) DeleteIdea(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	request.ParseForm()
	idea := &entity.Idea{}
	systemmessage := &entity.SystemMessage{}
	jsonsystemmessage, _ := json.Marshal(systemmessage)
	ideaowner, _, _ := idea_Admin.Session.Valid(request)
	writer.Header().Add("Content-Type", "application/json")

	if ideaowner <= 0 {
		//  invalid Request Notification
		systemmessage.Message = "Invalid User Pleas Log In First "
		systemmessage.Succesful = false
		jsonsystemmessage, _ = json.Marshal(systemmessage)
		writer.Write(jsonsystemmessage)
	}
	idea.Ideaownerid = ideaowner

	id, erro := strconv.Atoi(request.FormValue("id"))
	if erro != nil || id <= 0 {
		systemmessage.Message = "Invalid Idea Id!"
		systemmessage.Succesful = false
		jsonsystemmessage, _ = json.Marshal(systemmessage)
		writer.Write(jsonsystemmessage)
		// Notifying the User That The Idea id is not Valid
	}
	idea.Id = id
	systemmessage = idea_Admin.ideaservice.DeleteIdea(idea)
	jsonMessage, _ := json.Marshal(systemmessage)
	writer.Write(jsonMessage)
}

func (idea_Admin *IdeaHandler) UpdateIdea(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	request.ParseForm()
	writer.Header().Add("Content-Type", "application/json")

	idea := &entity.Idea{}
	systemmessage := &entity.SystemMessage{}
	jsonsystemmessage, _ := json.Marshal(systemmessage)

	ideaowner, _, _ := idea_Admin.Session.Valid(request)
	if ideaowner > 0 {
		idea.Ideaownerid = ideaowner
	} else {

		systemmessage.Message = "Sign Up First Invalid User "
		jsonsystemmessage, _ := json.Marshal(systemmessage)

		writer.Write(jsonsystemmessage)

		// Invalid  Request  Response
	}
	id, err := strconv.Atoi(request.FormValue("id"))
	if err != nil || id <= 0 {

		systemmessage.Message = "Invalid Idea Id "
		jsonsystemmessage, _ := json.Marshal(systemmessage)

		writer.Write(jsonsystemmessage)

		// Response For Invalid Id
	}
	idea.Id = id
	title := request.FormValue("title")
	if title != "" {
		idea.Title = title
	}
	description := request.FormValue("description")
	if description != "" {
		idea.Description = description
	}
	// file, header, err := request.FormFile("file")
	visibility := request.FormValue("visibility")
	if visibility == "pu" || visibility == "pr" || visibility == "pv" {
		idea.Visibility = visibility
	}
	systemmessage = idea_Admin.ideaservice.UpdateIdea(idea)
	jsonsystemmessage, _ = json.Marshal(systemmessage)

	writer.Write(jsonsystemmessage)

}

func (idea_Admin *IdeaHandler) VoteIdea(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ideaid, err := strconv.Atoi(request.FormValue("id"))
	writer.Header().Add("Content-Type", "application/json")
	systemmessage := &entity.SystemMessage{}
	systemmessage.Message = "You Have To Log In To Vote An Idea"
	systemmessage.Succesful = false
	if err != nil {
		return
	}
	voterid, _, _ := idea_Admin.Session.Valid(request)
	if voterid <= 0 {
		jsonbinary, _ := json.Marshal(systemmessage)
		writer.Write(jsonbinary)
	}
	if err != nil {
		return
	}
	systemmessage = idea_Admin.ideaservice.VoteIdea(ideaid, voterid)
	jsonbinary, err := json.Marshal(systemmessage)
	if err != nil {
		return
	}
	writer.Write(jsonbinary)
}

// TODO This will be done after the Comment Application of the System Is Done
func (idea_Admin *IdeaHandler) GetDetailIdea(writer http.ResponseWriter, request *http.Request) *entity.GeneralIdeaPersonComments {
	idea := idea_Admin.GetIdea(writer, request)
	commentsList := []entity.Comment{}
	commentwithpersons := []entity.CommentWithPerson{}
	ideapersoncomment := &entity.GeneralIdeaPersonComments{}
	commentsList = idea_Admin.commentrepo.GetCommentsa(idea.Id)
	user := idea_Admin.userrouter.UserById(idea.Ideaownerid)
	IdeaPersonComments := &entity.IdeaPersonComments{}
	IdeaPersonComments.Idea = *idea
	IdeaPersonComments.Succesful = true
	if idea == nil || user == nil {
		ideapersoncomment.IdeaOwner = user
		ideapersoncomment.Succesful = false
		return ideapersoncomment
	}
	ideapersoncomment.IdeaOwner = user
	ideapersoncomment.Succesful = true
	for index, comment := range commentsList {

		usr := idea_Admin.userrouter.UserById(comment.Commentorid)
		newcommentwithperson := entity.CommentWithPerson{}
		newcommentwithperson.Person = *usr
		newcommentwithperson.Comment = &comment
		commentwithpersons[index] = newcommentwithperson
	}
	IdeaPersonComments.CommentAndPerson = commentwithpersons

	ideapersoncomment.Succesful = true
	ideapersoncomment.IdeaRelatedData = IdeaPersonComments
	IdeaPersonComments.Idea = ideapersoncomment.IdeaRelatedData.Idea

	return ideapersoncomment
}

func (idea_Admin *IdeaHandler) GetDetailIdeas(writer http.ResponseWriter, request *http.Request) entity.FullIdeaView {

	FullIdeaView := entity.FullIdeaView{}
	person := entity.Person{}
	Idea := entity.Idea{}
	// comments := []entity.Comment{}
	// commentwithperson := []entity.CommentWithPerson{}
	ideaid, err := strconv.Atoi(request.FormValue("ideaid"))
	if err != nil {
		fmt.Println("The Idea id Take  ")

		return FullIdeaView
	}
	Idea.Id = ideaid
	systemmessage := idea_Admin.ideaservice.GetIdea(&Idea, ideaid)
	if !systemmessage.Succesful {
		fmt.Println("The Idea Ckeck ")
		return FullIdeaView
	}

	person.ID = uint(Idea.Ideaownerid)
	systemmessage = idea_Admin.userrouter.userservice.GetUser(&person)
	if !systemmessage.Succesful {
		fmt.Println("The useckache  Ckeck ", person.ID)

		return FullIdeaView
	}
	// comments = idea_Admin.commentrepo.GetCommentsa(ideaid)
	// if len(comments) > 0 {
	// 	for idx, comment := range comments {
	// 		newperson := entity.Person{ID: uint(comment.Commentorid)}
	// 		idea_Admin.userrouter.userservice.CheckUser(&newperson)
	// 		commentWithPerson := entity.CommentWithPerson{
	// 			Person:  &newperson,
	// 			Comment: &comment,
	// 		}

	// 		commentwithperson[idx] = commentWithPerson
	// 	}

	// }
	FullIdeaView.Success = true
	FullIdeaView.Person = person
	FullIdeaView.Idea = Idea
	// FullIdeaView.CommentWithPerson = commentwithperson
	return FullIdeaView

}

func (ideahandler *IdeaHandler) GetCommentWithPerson(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	writer.Header().Add("Content-Type", "application/json")
	ideaid, err := strconv.Atoi(request.FormValue("ideaid"))
	commentwithperson := []entity.CommentWithPerson{}
	jsoncommentwithperson, _ := json.Marshal(commentwithperson)
	if err != nil {
		writer.Write(jsoncommentwithperson)
	}
	commentwithperson = ideahandler.userrouter.userservice.Userrepo.GetCommentWithPerson(ideaid)
	jsoncommentwithperson, _ = json.Marshal(commentwithperson)
	writer.Write(jsoncommentwithperson)

}

func (idea_Admin *IdeaHandler) TemplateGetDetailIdea(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	listed := idea_Admin.GetDetailIdeas(writer, request)

	fmt.Println(listed.Success)
	SystemTemplates.ExecuteTemplate(writer, "bviewidea.html", listed)
}

func (ideahandler *IdeaHandler) ApiMyIdeas(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	ideas := &[]entity.Idea{}
	id, _, _ := ideahandler.Session.Valid(request)
	if id <= 0 {
		SystemTemplates.ExecuteTemplate(writer, "youridea.html", nil)
		return
	}
	systemmessage := ideahandler.ideaservice.MyIdeas(id, ideas)
	if systemmessage.Succesful {
		SystemTemplates.ExecuteTemplate(writer, "youridea.html", ideas)
		return
	} else {
		SystemTemplates.ExecuteTemplate(writer, "youridea.html", nil)
		return
	}
}
