package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/ambassador/src/models"
	"gorm.io/gorm"
)

type userController struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *userController {
	return &userController{db: db}
}

func (h *userController) Ambassador(c *fiber.Ctx) error {
	var users []models.User

	h.db.Where("is_ambassador = true").Find(&users)

	type Response struct {
		ID        uint   `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	response := []Response{}

	for _, user := range users {
		result := Response{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}
		response = append(response, result)
	}

	return c.Status(200).JSON(response)
}

// func (h *userController) Ambassador(c *fiber.Ctx) error {
// 	return c.Status(200).JSON(fiber.Map{})
// }
