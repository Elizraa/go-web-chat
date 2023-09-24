package middlewares

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

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

type APICallInfo struct {
	ID       string `json:"id"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

// Middleware to capture API call information
func APICallMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique API call ID
		apiCallID := GenerateAPICallID()

		// Wrap the original ResponseWriter to capture the response
		rw := &responseWriterWithCapture{ResponseWriter: w}

		// Add the api_call_id to the request context
		ctx := context.WithValue(r.Context(), "api_call_id", apiCallID)

		// Create a logger instance with context for API calls
		apiCallLogger := log.With().
			Str("REQUEST_PATH", r.URL.Path).
			Str("ID", apiCallID).Logger()

		// Log API call information using the logger instance
		apiCallLogger.Info().Msg("API_CALL_REQUEST")

		next.ServeHTTP(rw, r.WithContext(ctx))

		// Log the API call information, including the response
		logAPICallInfo(apiCallID, r, rw)
	})
}

// Custom ResponseWriter to capture the response
type responseWriterWithCapture struct {
	http.ResponseWriter
	capturedResponse bytes.Buffer
}

func (rw *responseWriterWithCapture) Write(b []byte) (int, error) {
	// Capture the response data
	n, err := rw.capturedResponse.Write(b)
	if err != nil {
		return n, err
	}
	return rw.ResponseWriter.Write(b)
}

// Logging function for API call information
func logAPICallInfo(apiCallID string, r *http.Request, rw *responseWriterWithCapture) {
	// Log or store the API call information as needed
	// You can access the captured request and response data here
	apiCallInfo := APICallInfo{
		ID:       apiCallID,
		Request:  r.URL.Path,
		Response: rw.capturedResponse.String(),
	}

	// Create a logger instance with context for API calls
	apiCallLogger := log.With().
		Str("RESPONSE", apiCallInfo.Response).
		Str("ID", apiCallInfo.ID).Logger()

	// Log API call information using the logger instance
	apiCallLogger.Info().Msg("API_CALL_RESPONSE")

}
