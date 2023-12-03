package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elizraa/chitchat/config"
	"github.com/elizraa/chitchat/data"
	"github.com/elizraa/chitchat/handler"
	"github.com/joho/godotenv"
)

func main() {

	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			log.Printf("Error getting env, not coming through %v", err)
		} else {
			fmt.Println("env values loaded")
		}
	} else if os.IsNotExist(err) {
		fmt.Println(".env file does not exist. Continuing without loading environment variables.")
	} else {
		log.Printf("Error checking .env file: %v", err)
	}

	address := handler.Config.Address
	// If port env var is set, PaaS platform (Heroku) is being used
	if port, ok := os.LookupEnv("PORT"); ok {
		address = "0.0.0.0:" + port
	}

	config.InitDB(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	// Migrate the schema
	config.DB.AutoMigrate(&data.ChatRoomDB{}) // Assuming ChatRoomDB is your database model

	config.ConnectMongoDatabase()

	// starting up the server
	server := &http.Server{
		Addr:           address,
		Handler:        handler.Mux,
		ReadTimeout:    time.Duration(handler.Config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(handler.Config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("ChitChat", version(), "started at", server.Addr)
	if _, exist := os.LookupEnv("PORT"); exist {
		// TLS is already enabled on Heroku PaaS platform
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Error starting server", err.Error())
		}
	} else {
		if err := server.ListenAndServeTLS("gencert/cert.pem", "gencert/key.pem"); err != nil {
			// If TLS fails e.g. because certs are missing on CI test env, we will fallback to regular HTTP
			if err := server.ListenAndServe(); err != nil {
				fmt.Println("Error starting server", err.Error())
			}
		}
	}

}

// version
func version() string {
	return "0.4"
}
