package api

import (
	"fmt"
	"log"
	"os"

	"github.com/Elizraa/go-web-chat/api/controllers"
	"github.com/Elizraa/go-web-chat/api/seed"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {
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
	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	seed.Load(server.DB)

	server.Run(":" + os.Getenv("PORT"))
}
