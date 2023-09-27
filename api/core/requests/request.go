// In api/requests/my_request.go
package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type MyRequest struct {
	FullPath   string
	Payload    []byte
	Param      url.Values
	Header     *http.Header
	RemoteAddr string
	// Files      []string
	CtxID string
}

func NewMyRequest(fullPath string, payload []byte, param url.Values, h *http.Header, remoteAddr string) *MyRequest {

	req := MyRequest{
		FullPath:   fullPath,
		Payload:    payload,
		Param:      param,
		Header:     h,
		RemoteAddr: remoteAddr,
		// Files:      files,
		CtxID: GenerateAPICallID(),
	}

	// Convert the request dictionary to a string for logging
	// reqStr := fmt.Sprintf("%v", req.ToDict())

	// Create a logger instance with context for API calls
	apiCallLogger := log.With().
		Str("REQUEST", req.ToJSONString()). // Use the string representation here
		Str("ID", req.CtxID).
		Logger()

	// Log API call information using the logger instance
	apiCallLogger.Info().Msg("API_CALL_REQUEST")
	return &req

}

// ToJSONString converts the request data to a JSON string.
func (r *MyRequest) ToJSONString() string {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// Conversion methods
func (r *MyRequest) ToDict() map[string]interface{} {
	return map[string]interface{}{
		"full_path":   r.FullPath,
		"payload":     r.Payload,
		"param":       r.Param,
		"header":      r.Header,
		"remote_addr": r.RemoteAddr,
		// "files":       r.Files,
		"ctx_id": r.CtxID,
	}
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
