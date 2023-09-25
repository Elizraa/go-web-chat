package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Elizraa/go-web-chat/api/models"
	"github.com/Elizraa/go-web-chat/api/responses"
	"github.com/Elizraa/go-web-chat/api/utils"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// // Parse the JSON request body into the user struct
	// if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }

	fmt.Println("======================================")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("err", err)

		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = user.Validate("")
	if err != nil {
		fmt.Println("aaaaaaaaaaa")
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	fmt.Println(user)
	fmt.Println("======================================")

	// Validate user input (e.g., check for required fields)

	// Hash the user's password before saving it to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user.Password = string(hashedPassword)

	userCreated, err := user.CreateUser(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)

}

// LoginResponse represents the JSON response for successful login
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("err", err)

		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Authenticate the user (e.g., validate username and password)
	// Replace this with your own authentication logic
	// Here, we assume you have a function `authenticateUser` that returns a user if authentication is successful
	authenticatedUser, err := user.AuthenticateUser(server.DB, user.Email, user.Password)
	if err != nil {
		err := utils.FormatError("Invald email or password")
		responses.ERROR(w, http.StatusUnauthorized, err)
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
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	// Respond with the access token
	response := LoginResponse{
		AccessToken: tokenString,
	}

	responses.JSON(w, http.StatusOK, response)
}
