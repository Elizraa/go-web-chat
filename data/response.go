package data

import (
	"encoding/json"
	"net/http"
	"time"
	// "github.com/rs/zerolog/log"
)

type MyResponse struct {
	Status                 string      `json:"status"`
	MessageID              string      `json:"message_id"`
	MessageAction          string      `json:"message_action"`
	MessageDesc            string      `json:"message_desc"`
	MessageData            interface{} `json:"message_data"`
	MessageRequestDatetime string      `json:"message_request_datetime"`
}

func NewMyResponse(ctxID, status, desc string, data interface{}) *MyResponse {
	now := time.Now().Format("2006-01-02 15:04:05")

	return &MyResponse{
		Status:                 status,
		MessageID:              ctxID,
		MessageAction:          status,
		MessageDesc:            desc,
		MessageData:            data,
		MessageRequestDatetime: now,
	}
}

func (r *MyResponse) ToDict() map[string]interface{} {
	return map[string]interface{}{
		"status":                   r.Status,
		"message_id":               r.MessageID,
		"message_action":           r.MessageAction,
		"message_desc":             r.MessageDesc,
		"message_data":             r.MessageData,
		"message_request_datetime": r.MessageRequestDatetime,
	}
}

// ToJSONString converts the request data to a JSON string.
func (r *MyResponse) ToJSONString() string {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// WriteToResponse writes the response to an http.ResponseWriter with the given status code and data.
func (r *MyResponse) WriteToResponse(w http.ResponseWriter, statusCode int, responseData interface{}) {
	// Set the Content-Type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code
	w.WriteHeader(statusCode)

	if responseData != nil {
		r.MessageData = responseData
	}

	if statusCode >= 200 && statusCode < 300 {
		r.MessageDesc = "Request successfully processed"
	} else {
		r.MessageDesc = "Failed request, please try again later"
	}

	// Construct a new response dictionary
	response := map[string]interface{}{
		"message_id":               r.MessageID,
		"message_desc":             r.MessageDesc,
		"message_data":             r.MessageData,
		"message_request_datetime": r.MessageRequestDatetime,
	}

	// Encode and send the response to the ResponseWriter
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Handle encoding errors if needed
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// // Create a logger instance with context for API calls
	// apiCallLogger := log.With().
	// 	Str("RESPONSE", r.ToJSONString()). // Use the string representation here
	// 	Str("ID", r.MessageID).
	// 	Logger()

	// // Log API call information using the logger instance
	// apiCallLogger.Info().Msg("API_CALL_RESPONSE")
}
