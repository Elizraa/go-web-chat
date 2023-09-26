package controllers

import (
	"fmt"
	"net/http"

	"github.com/Elizraa/go-web-chat/api/core/requests"
	"github.com/Elizraa/go-web-chat/api/core/responses"
)

// Home response
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	// Get the request and response objects from the context
	myRequest := r.Context().Value("myRequest").(*requests.MyRequest)

	fmt.Println(myRequest.CtxID)

	// responses.JSON(w, http.StatusOK, "Welcome To This Awesome API", myRequest.CtxID)
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	myResponse.WriteToResponse(w, http.StatusOK, "Welcome To This Awesome API")
}
