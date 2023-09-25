package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Elizraa/go-web-chat/api/controllers"
	"github.com/Elizraa/go-web-chat/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var productInstance = models.Product{}

// Prepare environment variables, and DB connection
func testMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

// Connect to the database
func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

// Drop the table and re-migrate
func refreshProductTable() error {
	err := server.DB.DropTableIfExists(&models.Product{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Product{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed product table")
	return nil
}

// Seed a single product into the database
func seedOneProduct() (models.Product, error) {

	refreshProductTable()

	var product = models.Product{
		Name:  "Test Product",
		Price: 22.05,
	}

	err := server.DB.Model(&models.Product{}).Create(&product).Error
	if err != nil {
		log.Fatalf("cannot seed scans table: %v", err)
	}
	return product, nil
}
