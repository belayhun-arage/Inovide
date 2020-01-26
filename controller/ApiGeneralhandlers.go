package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	session "github.com/Projects/Inovide/Session"
	entity "github.com/Projects/Inovide/models"
	"github.com/julienschmidt/httprouter"
)

type ApiUserHandler struct {
	userhandler *UserHandler
	session     *session.Cookiehandler
}

func NewApiController(userhandlers *UserHandler, sess *session.Cookiehandler) *ApiUserHandler {
	return &ApiUserHandler{
		userhandler: userhandlers,
		session:     sess,
	}
}

func (apiuserhandler *ApiUserHandler) LogIn(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

}

func (ideahandler *IdeaHandler) SearchResult(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	person := &entity.Person{}

	ideas := &[]entity.Idea{}
	jsonNullgeneralideasearchresult, _ := json.Marshal(ideas)

	writer.Header().Add("Content-Type", "application/json")
	searchtitle := request.FormValue("text")
	if searchtitle != "" {
		if searchtitle == "" {
			writer.Write(jsonNullgeneralideasearchresult)
		}
		id, username, _ := ideahandler.Session.Valid(request)
		if id > 0 {
			person.Paid = 0
		} else {
			person.Username = username
			person.ID = uint(id)
			systemmessage := ideahandler.userrouter.userservice.GetUser(person)
			if !systemmessage.Succesful {
				person.Paid = 0
			}
		}
		systemmessage := ideahandler.ideaservice.SearchResult(searchtitle, person, ideas)

		if !systemmessage.Succesful {
			writer.Write(jsonNullgeneralideasearchresult)
		}
	}
	newjson, erro := json.Marshal(ideas)
	if erro != nil {
		writer.Write(jsonNullgeneralideasearchresult)
	}

	writer.Write(newjson)
}

func (userhandler *UserHandler) FollowUser(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	// username, password, present := ReadSession(request)
	id, username, present := userhandler.Sessionservice.Valid(request)

	systemMessage := &entity.SystemMessage{}
	systemmessagejson, _ := json.Marshal(systemMessage)
	person := &entity.Person{}
	_, err := strconv.Atoi(request.FormValue("followingid"))
	if err != nil {
		writer.Write(systemmessagejson)
	}

	request.Header.Add("Content-Type", "application/json")
	if !present {
		writer.Write(systemmessagejson)
	}
	person.Username = username
	person.ID = uint(id)
	systemmessage := userhandler.userservice.GetUser(person)
	if !systemmessage.Succesful {
		writer.Write(systemmessagejson)
	}

	if person.ID != 0 {
		writer.Write(systemmessagejson)
	}

	systemmessagejson, err = json.Marshal(person)
	if err != nil {
		writer.Write(systemmessagejson)
	}
	writer.Write(systemmessagejson)

}

func (apiuserhandler *ApiUserHandler) ApiRegisterUser(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	systemmessage := apiuserhandler.userhandler.RegisterUser(writer, request)
	jsonmessage, _ := json.Marshal(systemmessage)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(jsonmessage)
}

func (apiuserhandler *ApiUserHandler) ApiSignin(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

	systemmessage, _ := apiuserhandler.userhandler.LogInRequest(writer, request, nil)
	jsontoken, _ := json.Marshal(systemmessage)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(jsontoken)
}

// func
