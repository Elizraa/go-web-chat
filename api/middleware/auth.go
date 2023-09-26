// middlewares/authmiddleware.go
package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Elizraa/go-web-chat/api/core/responses"
	"github.com/Elizraa/go-web-chat/api/utils" // Import utility functions for handling errors

	"github.com/dgrijalva/jwt-go"
)

func RequireTokenAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		// Check if the tokenString starts with "Bearer "
		if strings.HasPrefix(tokenString, "Bearer ") {
			// Remove the "Bearer " prefix
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		if tokenString == "" {
			err := utils.FormatError("Authorization token is required")
			responses.ERROR(w, http.StatusUnauthorized, err)
			return
		}
		fmt.Println(tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// You should provide the secret key used to sign the tokens here
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, err)
			return
		}
		if !token.Valid {
			err := utils.FormatError("Invald or expired token")
			responses.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		// You can access the user claims here if needed
		claims, _ := token.Claims.(jwt.MapClaims)

		// Store the user information in the request context for use in your handlers
		r = r.WithContext(context.WithValue(r.Context(), "user", claims))
		next(w, r)
	}
}
