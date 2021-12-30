package domain

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           string `json:"id" bson:"_id,omitempty"`
	FirstName    string `json:"firstName" bson:"firstName,omitempty"`
	LastName     string `json:"lastName" bson:"lastName,omitempty"`
	UserId       string `json:"userId" bson:"userId,omitempty"`
	PasswordHash string `json:"-" bson:"passwordHash,omitempty"`
}

// IsCorrectPassword checks if the provided password is correct or not
func (user *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

// encryptPassword perform a bcrypt encryption and return the byte and error
func encryptPassword(password string) ([]byte, error) {
	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		return nil, hashedPasswordErr
	}
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return nil, fmt.Errorf("could not encrypt and verify password")
	}
	return hashedPassword, nil
}
