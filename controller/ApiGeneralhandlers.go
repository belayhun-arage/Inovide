package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	entity "github.com/Projects/Inovide/models"
	"github.com/julienschmidt/httprouter"
)

type ApiUserHandler struct {
	userhandler *UserHandler
}

func (apiuserhandler *ApiUserHandler) LogIn(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

}

func (ideahandler *IdeaHandler) SearchResult(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	person := &entity.Person{}

	ideas := &[]entity.Idea{}
	jsonNullgeneralideasearchresult, _ := json.Marshal(ideas)

	searchtitle := request.FormValue("text")
	request.Header.Add("Content-Type", "application/json")

	if searchtitle == "" {
		writer.Write(jsonNullgeneralideasearchresult)
	}

	username, password, present := ReadSession(request)
	if !present {
		person.Paid = 0
	} else {
		person.Username = username
		person.Password = password
		systemmessage := ideahandler.userrouter.userservice.GetUser(person)
		if !systemmessage.Succesful {
			person.Paid = 0
		}
	}
	systemmessage := ideahandler.ideaservice.SearchResult(searchtitle, person, ideas)

	if !systemmessage.Succesful {
		writer.Write(jsonNullgeneralideasearchresult)
	}
	newjson, erro := json.Marshal(ideas)
	if erro != nil {
		writer.Write(jsonNullgeneralideasearchresult)
	}
	writer.Write(newjson)
}

func (userhandler *UserHandler) FollowUser(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	username, password, present := ReadSession(request)

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
	person.Password = password
	systemmessage := userhandler.userservice.GetUser(person)
	if !systemmessage.Succesful {
		writer.Write(systemmessagejson)
	}

	if person.ID != 0 {
		writer.Write(systemmessagejson)
	}

	systemmessagejson, err := json.Marshal(person)
	if err != nil {
		writer.Write(systemmessagejson)
	}
	writer.Write(systemmessagejson)

}
