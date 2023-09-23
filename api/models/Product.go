package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Product Interface
type Product struct {
	ID    uint32  `gorm:"primary_key;auto_increment" json:"id"`
	Name  string  `gorm:"size:255;not null;" json:"name"`
	Price float64 `json:"price"`
	BaseModel
}

// GetProduct will return a single product given a valid ID
func (p *Product) GetProduct(db *gorm.DB, uid uint32) (*Product, error) {
	err := db.Debug().Model(Product{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Product{}, errors.New("Product Not Found")
	}

	return p, err
}

// updateProduct will update an individual product row given a valid ID
func (p *Product) UpdateProduct(db *gorm.DB, uid uint32) (*Product, error) {

	db = db.Debug().Model(&Product{}).Where("id = ?", uid).Take(&Product{}).UpdateColumns(
		map[string]interface{}{
			"Name":  p.Name,
			"Price": p.Price,
		},
	)
	if db.Error != nil {
		return &Product{}, db.Error
	}
	// This is to display the updated product
	err := db.Debug().Model(&Product{}).Where("id = ?", uid).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}
	return p, nil
}

// deleteProduct will remove a product row given a valid ID
func (p *Product) DeleteProduct(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&Product{}).Where("id = ?", uid).Take(&Product{}).Delete(&Product{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// createProduct will save a new product given the correct payload
func (p *Product) CreateProduct(db *gorm.DB) (*Product, error) {
	err := db.Debug().Create(&p).Error
	if err != nil {
		return &Product{}, err
	}
	return p, nil
}

// Return all scans, limiting by 100
func (p *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
	products := []Product{}
	err := db.Debug().Model(&Product{}).Limit(100).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}
	return &products, err
}

// Validate a product
func (p *Product) Validate(action string) error {
	switch action {
	case "normal":
		if p.Name == "" {
			return errors.New("required name")
		}
		return nil
	default:
		if p.Name == "" {
			return errors.New("required name")
		}
		return nil
	}
}
