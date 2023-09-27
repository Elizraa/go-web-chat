package controllers

import middlewares "github.com/Elizraa/go-web-chat/api/middleware"

// initializeRoutes declares all routes used in the application
func (s *Server) initializeRoutes() {
	// s.Router.Use(middlewares.APICallMiddleware)
	// Use the middleware for all routes
	s.Router.Use(middlewares.RequestResponseMiddleware)
	// s.Router.Use(middlewares.MiddlewareJSON)

	s.Router.HandleFunc("/", s.Home).Methods("GET")
	SetupProductRoutes(s)
	SetupAuthRoutes(s)
}

func SetupProductRoutes(s *Server) {
	// Product Routes
	productRouter := s.Router.PathPrefix("/product").Subrouter()
	productRouter.HandleFunc("", s.FindAllProducts).Methods("GET")
	productRouter.HandleFunc("", middlewares.RequireTokenAuthentication(s.CreateProduct)).Methods("POST")
	productRouter.HandleFunc("/own", middlewares.RequireTokenAuthentication(s.GetProductByUser)).Methods("GET")
	productRouter.HandleFunc("/{id:[0-9]+}", s.GetProduct).Methods("GET")
	productRouter.HandleFunc("/{id:[0-9]+}", s.UpdateProduct).Methods("PUT")
	productRouter.HandleFunc("/{id:[0-9]+}", s.DeleteProduct).Methods("DELETE")
}

func SetupAuthRoutes(s *Server) {
	authRouter := s.Router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/regist", s.RegisterUser).Methods("POST")
	authRouter.HandleFunc("", s.Login).Methods("POST")
}
