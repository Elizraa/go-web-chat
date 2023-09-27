// middlewares/authmiddleware.go
package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/Elizraa/go-web-chat/api/core/responses"
	// Import utility functions for handling errors
	"github.com/dgrijalva/jwt-go"
)

func RequireTokenAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		myResponse := r.Context().Value("myResponse").(*responses.MyResponse)

		// Check if the tokenString starts with "Bearer "
		if strings.HasPrefix(tokenString, "Bearer ") {
			// Remove the "Bearer " prefix
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		if tokenString == "" {
			myResponse.WriteToResponse(w, http.StatusUnauthorized, "Authorization token is required")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			myResponse.WriteToResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		if !token.Valid {
			myResponse.WriteToResponse(w, http.StatusUnauthorized, "Invald or expired token")
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		// Store the user information in the request context for use in handlers
		r = r.WithContext(context.WithValue(r.Context(), "user", claims))
		next(w, r)
	}
}
