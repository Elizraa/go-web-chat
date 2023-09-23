package controllers

import middlewares "github.com/DLzer/go-product-api/api/middleware"

// initializeRoutes declares all routes used in the application
func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	SetupProductRoutes(s)
	SetupAuthRoutes(s)
}

func SetupProductRoutes(s *Server) {
	// Product Routes
	productRouter := s.Router.PathPrefix("/product").Subrouter()
	productRouter.HandleFunc("", middlewares.SetMiddlewareJSON(s.FindAllProducts)).Methods("GET")
	productRouter.HandleFunc("", middlewares.RequireTokenAuthentication((s.CreateProduct))).Methods("POST")
	productRouter.HandleFunc("/{id:[0-9]+}", middlewares.SetMiddlewareJSON(s.GetProduct)).Methods("GET")
	productRouter.HandleFunc("/{id:[0-9]+}", middlewares.SetMiddlewareJSON(s.UpdateProduct)).Methods("PUT")
	productRouter.HandleFunc("/{id:[0-9]+}", middlewares.SetMiddlewareJSON(s.DeleteProduct)).Methods("DELETE")
}

func SetupAuthRoutes(s *Server) {
	authRouter := s.Router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/regist", middlewares.SetMiddlewareJSON(s.RegisterUser)).Methods("POST")
	authRouter.HandleFunc("", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
}
