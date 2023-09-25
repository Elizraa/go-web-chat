package controllers

import (
	"net/http"

	"github.com/Elizraa/go-web-chat/api/responses"
)

// Home response
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}
