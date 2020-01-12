package api_handler

import (
	handler "github.com/Projects/Inovide/controller"
	"net/http"
)

type ApiUserHandler struct {
	userhandler *handler.UserHandler
}

func (apiuserhandler *ApiUserHandler) LogIn(writer http.ResponseWriter, request *http.Request) {

}
