package modeltests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver

	"github.com/DLzer/go-product-api/api/models"
	"gopkg.in/go-playground/assert.v1"
)

// Testing returning all scans
func TestFindAllProducts(t *testing.T) {

	err := refreshProductTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedOneProduct()
	if err != nil {
		log.Fatal(err)
	}

	products, err := productInstance.FindAllProducts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the scans: %v\n", err)
		return
	}
	assert.Equal(t, len(*products), 1)
}

// Test saving a single product
func TestSaveProduct(t *testing.T) {

	err := refreshProductTable()
	if err != nil {
		log.Fatal(err)
	}

	newProduct := models.Product{
		Name:  "Test Product",
		Price: 22.05,
	}

	savedProduct, err := newProduct.CreateProduct(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the products: %v\n", err)
		return
	}

	assert.Equal(t, newProduct.ID, savedProduct.ID)
}

// Test getting a single product by ID
func TestGetProductByID(t *testing.T) {

	err := refreshProductTable()
	if err != nil {
		log.Fatal(err)
	}

	product, err := seedOneProduct()
	if err != nil {
		log.Fatalf("cannot seed products table: %v", err)
	}
	foundScan, err := productInstance.GetProduct(server.DB, product.ID)
	if err != nil {
		t.Errorf("this is the error getting one product: %v\n", err)
		return
	}
	assert.Equal(t, foundScan.ID, product.ID)
}

// Test deleting a scan by ID
func TestDeleteScan(t *testing.T) {

	err := refreshProductTable()
	if err != nil {
		log.Fatal(err)
	}

	product, err := seedOneProduct()

	if err != nil {
		log.Fatalf("Cannot seed product: %v\n", err)
	}

	isDeleted, err := productInstance.DeleteProduct(server.DB, product.ID)
	if err != nil {
		t.Errorf("this is the error deleting the product: %v\n", err)
		return
	}

	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
