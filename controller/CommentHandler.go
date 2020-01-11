package handler

import (
	"net/http"
	"strconv"

	CommentService "github.com/Projects/Inovide/Comment/Service"
	entity "github.com/Projects/Inovide/models"
)

type CommentHandler struct {
	CommentService *CommentService.CommentService
	Userrouter     *UserHandler
}

func NewCommentHandler(commentService *CommentService.CommentService, userrouter *UserHandler) *CommentHandler {
	return &CommentHandler{CommentService: commentService, Userrouter: userrouter}
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

	username, password, id, present := ReadSession(request)
	if !present {
		return commentWithThePerson
	}

	person.Username = username
	person.ID = uint(id)
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

func (commentHandler *CommentHandler) GetComments(ideaid int) (*[]entity.Comment, bool) {
	coments := &[]entity.Comment{}
	systemmessage := commentHandler.CommentService.GetComments(coments, ideaid)
	if systemmessage.Succesful {
		return coments, true
	}
	return coments, false
}
