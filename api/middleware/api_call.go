package middlewares

import (
	"context"
	"net/http"

	"github.com/Elizraa/go-web-chat/api/core/requests"
	"github.com/Elizraa/go-web-chat/api/core/responses"
)

func RequestResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new request object
		myRequest := requests.NewMyRequest(r.URL.Path, nil, r.URL.Query(), &r.Header, r.RemoteAddr)

		// Create a new response object
		myResponse := responses.NewMyResponse(myRequest.CtxID, "", "", nil)

		// Pass the request and response objects to the next handler
		ctx := r.Context()
		ctx = context.WithValue(ctx, "myRequest", myRequest)
		ctx = context.WithValue(ctx, "myResponse", myResponse)

		// Call the next handler with the updated context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
