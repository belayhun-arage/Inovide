package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	CommentService "github.com/Projects/Inovide/Comment/Service"
	entity "github.com/Projects/Inovide/models"
	"github.com/julienschmidt/httprouter"
)

type CommentHandler struct {
	CommentService *CommentService.CommentService
	Userrouter     *UserHandler
}

func NewCommentHandler(commentService *CommentService.CommentService, userrouter *UserHandler) *CommentHandler {
	return &CommentHandler{CommentService: commentService, Userrouter: userrouter}
}

func (commenthandler *CommentHandler) APICreateComment(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	commentwithperson := commenthandler.CreateComment(writer, request)

	jsonized, err := json.Marshal(commentwithperson)

	jsonfailde, _ := json.Marshal(&entity.CommentWithPerson{})

	if err != nil {
		writer.Write(jsonfailde)
	}
	writer.Write(jsonized)
}

func (commentHandler *CommentHandler) CreateComment(writer http.ResponseWriter, request *http.Request) *entity.CommentWithPerson {
	request.ParseForm()
	commentWithThePerson := &entity.CommentWithPerson{}

	person := &entity.Person{}
	comment := &entity.Comment{}

	ideaid, err := strconv.Atoi(request.FormValue("ideaid"))

	if err != nil {
		return commentWithThePerson
	}
	commentdata := request.FormValue("commentdata")

	username, password, present := ReadSession(request)
	if !present {
		return commentWithThePerson
	}

	person.Username = username
	person.Password = password
	systemmessageforuser := commentHandler.Userrouter.userservice.GetUser(person)
	if !systemmessageforuser.Succesful {
		return nil
	}

	if commentdata == "" {
		return commentWithThePerson
	}
	comment.Ideaid = ideaid
	comment.Commentdata = commentdata
	comment.Commentorid = int(person.ID)

	systemMessage := commentHandler.CommentService.CreateComment(comment)

	if systemMessage.Succesful {
		commentWithThePerson.Person = person
		commentWithThePerson.Comment = comment
		commentWithThePerson.Succesfull = true
		return commentWithThePerson
	}
	commentWithThePerson.Succesfull = false
	commentWithThePerson.Person = person
	commentWithThePerson.Comment = comment
	return commentWithThePerson
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
