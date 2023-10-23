package controllers

import (
	"net/http"

	"github.com/Elizraa/go-web-chat/api/core/responses"
)

// Home response
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	// Get the request and response objects from the context
	// myRequest := r.Context().Value("myRequest").(*requests.MyRequest)

	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	defer func() {
		if r := recover(); r != nil {
			myResponse.WriteToResponse(w, http.StatusInternalServerError, r)
		}
	}()
	myResponse.WriteToResponse(w, http.StatusOK, "Welcome To This Awesome API")
}
