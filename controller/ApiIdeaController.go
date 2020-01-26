package handler

import session "github.com/Projects/Inovide/Session"

import "net/http"

import "github.com/julienschmidt/httprouter"

import entity "github.com/Projects/Inovide/models"

import "encoding/json"

type ApiIdeaHandler struct {
	Userrouter     *UserHandler
	ideaHandler    *IdeaHandler
	commenthandler *CommentHandler
	Session        *session.Cookiehandler
}

func NewApiIdeaHandler(userhandler *UserHandler,
	ideahandler *IdeaHandler,
	commenthandler *CommentHandler,
	session *session.Cookiehandler) *ApiIdeaHandler {

	return &ApiIdeaHandler{
		Userrouter:     userhandler,
		ideaHandler:    ideahandler,
		commenthandler: commenthandler,
		Session:        session,
	}
}
func (ideahandler *ApiIdeaHandler) CreateIdea(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	systemmessage := &entity.SystemMessage{}
	idea := &entity.Idea{}
	systemmessage, idea = ideahandler.ideaHandler.CreateIdea(writer, request, nil)
	if systemmessage.Succesful {
		systemmessage.Succesful = true
		systemmessage.Message = "Idea Created Succesfully"
	}
	var tobesent = map[string]interface{}{"SystemMessage": systemmessage, "Idea": idea}
	jsontobesent, _ := json.Marshal(tobesent)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(jsontobesent)
}
func (user_controller *UserHandler) ApiEditeProfile(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	person, systemmessage := user_controller.EditProfile(writer, request, nil)
	finalfile := map[string]interface{}{
		"Person":  person,
		"Message": systemmessage,
	}
	jsonfile, _ := json.Marshal(finalfile)
	writer.Header().Add("Content-Type", "applicaion/json")
	writer.Write(jsonfile)
}
