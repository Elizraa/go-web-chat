package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/elizraa/chitchat/data"
	"github.com/google/uuid"
)

type errHandler func(http.ResponseWriter, *http.Request, string) error

func (fn errHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctxId := GenerateAPICallID()
	if err := fn(w, r, ctxId); err != nil {
		if apierr, ok := err.(*data.APIError); ok {

			w.Header().Set("Content-Type", "application/json")
			apierr.SetMsg()
			Warning("API error:", apierr.Error())
			if apierr.Code == 101 || apierr.Code == 201 {
				notFound(w, r)
			} else if apierr.Code == 102 || apierr.Code == 202 || apierr.Code == 303 || apierr.Code == 105 {
				badRequest(w, r)
			} else if apierr.Code == 104 || apierr.Code == 204 || apierr.Code == 304 || apierr.Code == 401 || apierr.Code == 402 {
				unauthorized(w, r)
			} else if apierr.Code == 403 {
				forbidden(w, r)
			} else {
				badRequest(w, r)
			}
			ReportStatus(w, false, apierr, ctxId)
		} else {
			Danger("Server error", err.Error())
			WriteErrorResponse(w, ctxId, err.Error())
		}
	}
	Info("[REQUEST]", ctxId, r.RequestURI)
}

func WriteErrorResponse(w http.ResponseWriter, ctxId string, errMsg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)
	// Create a map to hold ctxId and response data
	responseDataWithCtxId := map[string]interface{}{
		"API_CALL_ID": ctxId,
		"error":       errMsg,
		"data":        map[string]interface{}{},
		"status":      false,
	}
	// Convert responseDataWithCtxId to JSON
	jsonData, err := json.Marshal(responseDataWithCtxId)
	if err != nil {
		Danger("Error marshaling response data", err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// Write the JSON response to the http.ResponseWriter
	w.Write(jsonData)
}

// WriteResponse writes the ctxId and response data to the http.ResponseWriter
func WriteResponse(w http.ResponseWriter, ctxId string, responseData interface{}) {
	// You can customize this function based on your response format
	w.Header().Set("Content-Type", "application/json")
	// Create a map to hold ctxId and response data
	responseDataWithCtxId := map[string]interface{}{
		"API_CALL_ID": ctxId,
		"data":        responseData,
		"status":      true,
		"error":       map[string]interface{}{},
	}
	// Convert responseDataWithCtxId to JSON
	jsonData, err := json.Marshal(responseDataWithCtxId)
	if err != nil {
		Danger("Error marshaling response data", err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// Write the JSON response to the http.ResponseWriter
	w.Write(jsonData)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	Info("Not found request:", r.RequestURI)
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(401)
	Info("forbidden:", r.RequestURI, r.Body)
}

func forbidden(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	Warning("forbidden:", r.RequestURI, r.Body)
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	Info("Bad request:", r.RequestURI, r.Body)
}

// Convenience function to redirect to the error message page
func errorMessage(writer http.ResponseWriter, request *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(writer, request, strings.Join(url, ""), 302)
}

// GET /err?msg=
// shows the error message page
func handleError(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	fmt.Fprintf(writer, "Error: %s!", vals.Get("msg"))
	Warning(fmt.Sprintf("Error: %s!", vals.Get("msg")))
}

func GenerateAPICallID() string {
	// Generate a UUID
	uuidStr := uuid.New().String()

	// Get the current timestamp epoch in milliseconds
	timestampEpoch := time.Now().UnixNano() / 1e6

	// Extract the first 5 characters from the UUID
	uuidShort := uuidStr[:5]

	// Create the API call ID using the components
	apiCallID := fmt.Sprintf("API_CALL_%d_%s", timestampEpoch, uuidShort)

	return apiCallID
}
