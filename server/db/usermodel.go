package db

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// UserModel is used to store user informatio
type UserModel struct {
	gorm.Model
	Username string
	Email    string
	Password string
	IsActive bool
}

// SetPassword hashes and stores user password
func (u *UserModel) SetPassword(p string) error {
	password := []byte(p)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// MatchPassword matches plaintext password with stored hash password
func (u UserModel) MatchPassword(p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return false
	}
	return true
}
