package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint
	FirstName    string
	LastName     string
	Email        string `gorm:"unique"`
	Password     string
	IsAmbassador bool
	Revenue      *float64 `json:"revenue,omitempty" gorm:"-"`
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

type Admin User

func (admin *Admin) CalculateRevenue(db *gorm.DB) {
	var orders []Order

	db.Preload("OrderItems").Find(&orders, &Order{
		UserId:   admin.ID,
		Complete: true,
	})

	var revenue float64 = 0

	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AdminRevenue
		}
	}

	admin.Revenue = &revenue
}

type Ambassador User

func (ambassador *Ambassador) CalculateRevenue(db *gorm.DB) {
	var orders []Order

	db.Preload("OrderItems").Find(&orders, &Order{
		UserId:   ambassador.ID,
		Complete: true,
	})

	var revenue float64 = 0.0

	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AmbassadorRevenue
		}
	}

	ambassador.Revenue = &revenue
}
