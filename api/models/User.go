package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User Interface
type User struct {
	ID       uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name     string `gorm:"size:255;not null;" json:"name"`
	Email    string `gorm:"size:255;unique;not null;" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
	BaseModel
}

// createUser will save a new User given the correct payload
func (u *User) CreateUser(db *gorm.DB) (*User, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// Validate a user
func (u *User) Validate(action string) error {
	switch action {
	case "normal":
		if u.Name == "" {
			return errors.New("required name")
		}
		return nil
	default:
		if u.Name == "" {
			return errors.New("required name")
		}
		if u.Email == "" {
			return errors.New("required email")
		}
		return nil
	}
}

// Replace this function with your own authentication logic
func (u *User) AuthenticateUser(db *gorm.DB, email, password string) (*User, error) {
	// Query the database to find the user by username
	if err := db.Where("email= ?", email).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
		return nil, err
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, err
	}

	return u, nil
}
