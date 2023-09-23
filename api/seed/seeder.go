package seed

import (
	"log"

	"github.com/DLzer/go-product-api/api/models"
	"github.com/jinzhu/gorm"
)

var product = models.Product{
	Name:  "Test Product",
	Price: 22.05,
}

func Load(db *gorm.DB) {
	// err := db.Debug().DropTableIfExists(&models.Product{}, &models.Product{}).Error
	// if err != nil {
	// 	log.Fatalf("cannot drop table: %v", err)
	// }
	// err = db.Debug().AutoMigrate(&models.Product{}, &models.Product{}).Error
	// if err != nil {
	// 	log.Fatalf("cannot migrate table: %v", err)
	// }
	// err = db.Debug().Model(&models.Product{}).Create(&product).Error
	// if err != nil {
	// 	log.Fatalf("cannot seed scan table: %v", err)
	// }

	err := db.Debug().AutoMigrate(&models.User{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

}
