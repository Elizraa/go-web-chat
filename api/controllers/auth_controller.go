package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Elizraa/go-web-chat/api/core/responses"
	"github.com/Elizraa/go-web-chat/api/models"
	"github.com/Elizraa/go-web-chat/api/utils"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	defer func() {
		if r := recover(); r != nil {
			myResponse.WriteToResponse(w, http.StatusInternalServerError, r)
		}
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("err", err.Error())

		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = user.Validate("")
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	// Validate user input (e.g., check for required fields)

	// Hash the user's password before saving it to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	user.Password = string(hashedPassword)

	userCreated, err := user.CreateUser(server.DB)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	myResponse.WriteToResponse(w, http.StatusCreated, userCreated)
}

// LoginResponse represents the JSON response for successful login
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	defer func() {
		if r := recover(); r != nil {
			myResponse.WriteToResponse(w, http.StatusInternalServerError, r)
		}
	}()
	user := models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("err", err.Error())

		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Authenticate the user (e.g., validate username and password)
	// Replace this with your own authentication logic
	// Here, we assume you have a function `authenticateUser` that returns a user if authentication is successful
	authenticatedUser, err := user.AuthenticateUser(server.DB, user.Email, user.Password)
	if err != nil {
		err := utils.FormatError("Invald email or password")
		myResponse.WriteToResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Set claims (user information) in the token
	claims["user_id"] = authenticatedUser.ID
	claims["email"] = authenticatedUser.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time (adjust as needed)

	// Sign the token with your secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		err := utils.FormatError("Could not generate token")
		myResponse.WriteToResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Respond with the access token
	response := LoginResponse{
		AccessToken: tokenString,
	}

	myResponse.WriteToResponse(w, http.StatusOK, response)
}
