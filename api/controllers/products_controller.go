package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Elizraa/go-web-chat/api/core/responses"
	"github.com/Elizraa/go-web-chat/api/models"
	"github.com/gorilla/mux"
)

// Create a single product
func (server *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err)
	}
	product := models.Product{}
	err = json.Unmarshal(body, &product)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = product.Validate("")
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err)
		return
	}
	fmt.Println(product)

	productCreated, err := product.CreateProduct(server.DB)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, productCreated.ID))
	myResponse.WriteToResponse(w, http.StatusCreated, productCreated)
}

func (server *Server) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	body, err := ioutil.ReadAll(r.Body)
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)

	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err)
	}
	product := models.Product{}
	err = json.Unmarshal(body, &product)

	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = product.Validate("")
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	productUpdated, err := product.UpdateProduct(server.DB, uint32(uid))
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, productUpdated.ID))
	myResponse.WriteToResponse(w, http.StatusCreated, productUpdated)

}

// Get a single product
func (server *Server) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Get the request and response objects from the context
	// myRequest := r.Context().Value("myRequest").(*requests.MyRequest)
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)

	vars := mux.Vars(r)
	// Retrieve the api_call_id from the request context
	// apiCallID := r.Context().Value("api_call_id").(string)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusBadRequest, err)
		return
	}
	product := models.Product{}
	productGotten, err := product.GetProduct(server.DB, uint32(uid))
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusBadRequest, err)
		return
	}
	myResponse.WriteToResponse(w, http.StatusOK, productGotten)
}

// Get all products
func (server *Server) FindAllProducts(w http.ResponseWriter, r *http.Request) {
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)

	product := models.Product{}

	products, err := product.FindAllProducts(server.DB)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusInternalServerError, err)
		return
	}
	myResponse.WriteToResponse(w, http.StatusOK, products)

}

// Delete a single product
func (server *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	myResponse := r.Context().Value("myResponse").(*responses.MyResponse)

	vars := mux.Vars(r)

	product := models.Product{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusBadRequest, err)
		return
	}

	_, err = product.DeleteProduct(server.DB, uint32(uid))
	if err != nil {
		myResponse.WriteToResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	myResponse.WriteToResponse(w, http.StatusNoContent, "")
}
