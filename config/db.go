package config

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

var DB *gorm.DB

func InitDB(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var DBURL string
	switch Dbdriver {
	case "mysql":
		DBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	case "postgres":
		DBURL = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", DbHost, DbPort, DbUser, DbName, DbPassword)
	default:
		log.Fatal("Invalid database driver")
		return
	}

	db, err := gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", Dbdriver)
	}

	// Assign the database connection to the global variable
	DB = db
}
