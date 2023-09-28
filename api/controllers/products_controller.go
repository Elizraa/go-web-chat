package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Elizraa/go-web-chat/api/core/responses"
	"github.com/Elizraa/go-web-chat/api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// Create a single product
func (server *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	defer func() {
		if r := recover(); r != nil {
			myResponse.WriteToResponse(w, http.StatusInternalServerError, r)
		}
	}()
	body, err := ioutil.ReadAll(r.Body)
	userClaims, _ := r.Context().Value("user").(jwt.MapClaims)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	product := models.Product{}
	err = json.Unmarshal(body, &product)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = product.Validate("")
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	product.UserID = uint32(userClaims["user_id"].(float64))

	productCreated, err := product.CreateProduct(server.DB)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, productCreated.ID))
	myResponse.WriteToResponse(w, http.StatusCreated, productCreated)
}

func (server *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	defer func() {
		if r := recover(); r != nil {
			myResponse.WriteToResponse(w, http.StatusInternalServerError, r)
		}
	}()
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
	}
	product := models.Product{}
	err = json.Unmarshal(body, &product)

	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = product.Validate("")
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	userClaims, _ := r.Context().Value("user").(jwt.MapClaims)
	userID := uint32(userClaims["user_id"].(float64))

	productGotten, err := product.GetProduct(server.DB, uint32(uid))
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if productGotten.UserID != userID {
		myResponse.WriteToResponse(w, http.StatusUnauthorized, "Not the owner of the product")
		return
	}
	productUpdated, err := product.UpdateProduct(server.DB, uint32(uid))
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, productUpdated.ID))
	myResponse.WriteToResponse(w, http.StatusCreated, productUpdated)

}

// Get a single product
func (server *Server) GetProduct(w http.ResponseWriter, r *http.Request) {
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	defer func() {
		if r := recover(); r != nil {
			myResponse.WriteToResponse(w, http.StatusInternalServerError, r)
		}
	}()
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	product := models.Product{}
	productGotten, err := product.GetProduct(server.DB, uint32(uid))
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	myResponse.WriteToResponse(w, http.StatusOK, productGotten)
}

// Get all products
func (server *Server) FindAllProducts(w http.ResponseWriter, r *http.Request) {
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	defer func() {
		if r := recover(); r != nil {
			myResponse.WriteToResponse(w, http.StatusInternalServerError, r)
		}
	}()
	product := models.Product{}

	products, err := product.FindAllProducts(server.DB)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	myResponse.WriteToResponse(w, http.StatusOK, products)

}

// Delete a single product
func (server *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	defer func() {
		if r := recover(); r != nil {
			myResponse.WriteToResponse(w, http.StatusInternalServerError, r)
		}
	}()
	vars := mux.Vars(r)

	product := models.Product{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = product.DeleteProduct(server.DB, uint32(uid))
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	myResponse.WriteToResponse(w, http.StatusNoContent, "")
}

// Get all product by user
func (server *Server) GetProductByUser(w http.ResponseWriter, r *http.Request) {
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	defer func() {
		if r := recover(); r != nil {
			myResponse.WriteToResponse(w, http.StatusInternalServerError, r)
		}
	}()

	userClaims, _ := r.Context().Value("user").(jwt.MapClaims)
	userID := uint32(userClaims["user_id"].(float64))

	product := models.Product{}
	products_result, err := product.GetProductByUser(server.DB, userID)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	myResponse.WriteToResponse(w, http.StatusOK, products_result)
}
