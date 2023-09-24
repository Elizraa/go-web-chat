package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Response with JSON
func JSON(w http.ResponseWriter, statusCode int, data interface{}, apiCallID ...string) {
	// currenTime := time.Now().Format(time.RFC3339)
	currenTime := time.Now().Format("2006-01-02 15:04:05")
	// Include the apiCallID in the response data
	responseData := map[string]interface{}{
		"message_id":               apiCallID[0],
		"message_request_datetime": currenTime,
		"message_data":             data, // Include the original data
	}

	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(responseData)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// Respond with an error
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
