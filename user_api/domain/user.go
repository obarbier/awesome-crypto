package domain

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           string `json:"id" bson:"_id"`
	FirstName    string `json:"firstName" bson:"firstName"`
	LastName     string `json:"lastName" bson:"lastName"`
	UserId       string `json:"userId" bson:"userId"`
	PasswordHash string `json:"passwordHash" bson:"passwordHash"`
}

// IsCorrectPassword checks if the provided password is correct or not
func (user *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
