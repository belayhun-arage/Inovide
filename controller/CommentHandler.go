package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	CommentService "github.com/Projects/Inovide/Comment/Service"
	session "github.com/Projects/Inovide/Session"
	entity "github.com/Projects/Inovide/models"
	"github.com/julienschmidt/httprouter"
)

type CommentHandler struct {
	CommentService *CommentService.CommentService
	Userrouter     *UserHandler
	Sessionservice *session.Cookiehandler
}

func NewCommentHandler(commentService *CommentService.CommentService, userrouter *UserHandler, session *session.Cookiehandler) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
		Userrouter:     userrouter,
		Sessionservice: session}
}

func (commenthandler *CommentHandler) APICreateComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	commentwithperson, systemmessage := commenthandler.CreateComment(writer, request, params)

	jsonized, err := json.Marshal(commentwithperson)

	jsonfailde, _ := json.Marshal(&entity.CommentWithPerson{})
	if !systemmessage.Succesful {
		writer.Write(jsonfailde)
	}
	if err != nil {
		writer.Write(jsonfailde)
	}
	writer.Write(jsonized)
}

// this method return an instance of CommentWithPerson for ApiCreateComment function
func (commentHandler *CommentHandler) CreateComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) (*entity.CommentWithPerson, entity.SystemMessage) {
	request.ParseForm()
	commentWithThePerson := &entity.CommentWithPerson{}
	systemmessage := entity.SystemMessage{}
	systemmessage.Succesful = false
	person := &entity.Person{}
	comment := &entity.Comment{}

	// ideaid, err := strconv.Atoi(request.FormValue("ideaid"))
	ideaid, err := strconv.Atoi(request.FormValue("ideaid"))

	if err != nil {
		fmt.Println("Creating Comment ")

		return commentWithThePerson, systemmessage
	}
	commentdata := request.FormValue("commentdata")
	id, username, _ := commentHandler.Sessionservice.Valid(request)
	if id <= 0 {
		fmt.Println("No The Perosn is Not valid ")
		return commentWithThePerson, systemmessage
	}
	person.Username = username
	person.ID = uint(id)
	systemmessageforuser := commentHandler.Userrouter.userservice.GetUser(person)
	if !systemmessageforuser.Succesful {
		fmt.Println("UsrNot Found ")
		commentWithThePerson.Person = *person
		return commentWithThePerson, systemmessage
	}
	if commentdata == "" {
		systemmessage.Message = "The CommentData Invalid "
		systemmessage.Succesful = false
		return commentWithThePerson, systemmessage
	}
	comment.Ideaid = ideaid
	comment.Commentdata = commentdata
	comment.Commentorid = int(person.ID)
	comment.Commentdate = time.Now().String()
	systemMessage := commentHandler.CommentService.CreateComment(comment)

	if systemMessage.Succesful {
		commentWithThePerson.Person = *person
		commentWithThePerson.Comment = comment
		commentWithThePerson.Succesfull = true
		systemMessage.Succesful = true
		return commentWithThePerson, systemmessage
	}
	commentWithThePerson.Succesfull = false
	commentWithThePerson.Person = *person
	commentWithThePerson.Comment = comment
	systemMessage.Succesful = false
	return commentWithThePerson, systemmessage
}
func (commentHandler *CommentHandler) GetCommentWithPersons(commentwithpersons *[]entity.CommentWithPerson, comments *[]entity.Comment) *[]entity.CommentWithPerson {
	commentos := []entity.CommentWithPerson{}
	for index, comment := range *comments {
		person := commentHandler.Userrouter.UserById(comment.Commentorid)
		if person == nil {
			return commentwithpersons
		}
		commentwithperson := entity.CommentWithPerson{}
		commentwithperson.Person = *person
		commentwithperson.Comment = &comment
		commentwithperson.Succesfull = true

		commentos[index] = commentwithperson
	}
	commentwithpersons = &commentos
	return commentwithpersons
}

func (commentHandler *CommentHandler) CommentOnIdea(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	commentwithperson := &entity.CommentWithPerson{}
	commentwithperson.Succesfull = false
	writer.Header().Add("Content-Type", "application/json")

	jsoncommentwithperson, _ := json.Marshal(commentwithperson)
	commentwithperson, systemmessage := commentHandler.CreateComment(writer, request, nil)

	if systemmessage.Succesful {
		jsoncommentwithperson, _ = json.Marshal(commentwithperson)
		commentwithperson.Succesfull = false
		writer.Write(jsoncommentwithperson)

	} else {
		jsoncommentwithperson, _ = json.Marshal(commentwithperson)
		commentwithperson.Succesfull = true
		writer.Write(jsoncommentwithperson)
	}
}

func (commentHandler *CommentHandler) ApiGetCommentListed(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	request.ParseForm()
	theideaid := request.FormValue("ideaid")
	ideaid, err := strconv.Atoi(theideaid)

	listtofcomment := &[]entity.Comment{}

	thejson, _ := json.Marshal(listtofcomment)

	if err != nil {
		request.Header.Add("Content-Type", "application/json")
		writer.Write(thejson)
	}

	succesfull := commentHandler.GetComments(listtofcomment, ideaid)
	if !succesfull {
		request.Header.Add("Content-Type", "application/json")
		writer.Write(thejson)
	}
	request.Header.Add("Content-Type", "application/json")
	thejsonmain, err := json.Marshal(listtofcomment)
	if err != nil {
		request.Header.Add("Content-Type", "application/json")
		writer.Write(thejson)
	}
	request.Header.Add("Content-Type", "application/json")
	writer.Write(thejsonmain)
}
func (commentHandler *CommentHandler) GetComments(comment *[]entity.Comment, ideaid int) bool {
	coments := &[]entity.Comment{}
	systemmessage := commentHandler.CommentService.GetComments(coments, ideaid)
	if systemmessage.Succesful {
		return true
	}
	return false
}
func (commentHandler *CommentHandler) DeleteComment(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, _, _ := commentHandler.Sessionservice.Valid(request)
	comment := &entity.Comment{}
	systemmessage := &entity.SystemMessage{}
	systemmessage.Succesful = false
	systemmessage.Message = "Can't Delete The Message "

	jsonsystemmessage, _ := json.Marshal(systemmessage)
	writer.Header().Add("Content-Type", "application/json")

	if id <= 0 {
		//Informing that the user is not valid
		writer.Write(jsonsystemmessage)
	}
	// ideaid, err := strconv.Atoi(request.FormValue("ideaid"))
	commentid, err := strconv.Atoi(request.FormValue("commentid"))
	if err != nil || commentid <= 0 { //  ideaid <= 0 ||
		writer.Write(jsonsystemmessage)
	}
	comment.Id = commentid
	// comment.Ideaid= ideaid
	comment.Commentorid = id
	systemmessage = commentHandler.CommentService.DeleteComment(comment)
	if systemmessage.Succesful {
		systemmessage.Message = "Succesfully Deleted the messsage"
		systemmessage.Succesful = true
		jsonsystemmessage, _ = json.Marshal(systemmessage)
	}
	writer.Write(jsonsystemmessage)
}
