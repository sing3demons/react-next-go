package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID           uint
	FirstName    string
	LastName     string
	Email        string `gorm:"unique"`
	Password     string
	IsAmbassador bool
}

func (u *User) EncryptPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
