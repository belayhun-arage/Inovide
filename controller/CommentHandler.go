package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	ideaid, err := strconv.Atoi(params.ByName("ideaid"))

	if err != nil {
		return commentWithThePerson, systemmessage
	}
	commentdata := request.FormValue("commentdata")
	id, username, present := commentHandler.Sessionservice.Valid(request)
	if !present {
		return commentWithThePerson, systemmessage
	}
	person.Username = username
	person.ID = uint(id)
	systemmessageforuser := commentHandler.Userrouter.userservice.GetUser(person)
	if !systemmessageforuser.Succesful {
		return nil, systemmessage
	}
	if commentdata == "" {
		return commentWithThePerson, systemmessage
	}
	comment.Ideaid = ideaid
	comment.Commentdata = commentdata
	comment.Commentorid = int(person.ID)
	systemMessage := commentHandler.CommentService.CreateComment(comment)

	if systemMessage.Succesful {
		commentWithThePerson.Person = person
		commentWithThePerson.Comment = comment
		commentWithThePerson.Succesfull = true
		systemMessage.Succesful = true
		return commentWithThePerson, systemmessage
	}
	commentWithThePerson.Succesfull = false
	commentWithThePerson.Person = person
	commentWithThePerson.Comment = comment
	systemMessage.Succesful = false
	return commentWithThePerson, systemmessage
}
func (commentHandler *CommentHandler) GetCommentWithPerson(comments *[]entity.Comment) *[]entity.CommentWithPerson {
	commentWithThePersons := []entity.CommentWithPerson{}
	for index, comment := range *comments {
		person := commentHandler.Userrouter.UserById(comment.Commentorid)
		if person == nil {
			return nil
		}
		commentwithperson := entity.CommentWithPerson{}
		commentwithperson.Person = person
		commentwithperson.Comment = &comment
		commentwithperson.Succesfull = true
		commentWithThePersons[index] = commentwithperson
	}
	return &commentWithThePersons
}
func (commentHandler *CommentHandler) ApiGetCommentListed(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	request.ParseForm()
	theideaid := request.FormValue("ideaid")
	ideaid, err := strconv.Atoi(theideaid)

	comments := &[]entity.Comment{}

	thejson, _ := json.Marshal(comments)

	if err != nil {
		request.Header.Add("Content-Type", "application/json")
		writer.Write(thejson)
	}

	listtofcomment, succesfull := commentHandler.GetComments(ideaid)
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
func (commentHandler *CommentHandler) GetComments(ideaid int) (*[]entity.Comment, bool) {
	coments := &[]entity.Comment{}
	systemmessage := commentHandler.CommentService.GetComments(coments, ideaid)
	if systemmessage.Succesful {
		return coments, true
	}
	return coments, false
}
