package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/ambassador/src/models"
	"gorm.io/gorm"
)

type Order struct {
	db *gorm.DB
}

func NewOrder(db *gorm.DB) *Order {
	return &Order{db}
}

func (h *Order) Order(c *fiber.Ctx) error {
	var orders []models.Order

	h.db.Preload("OrderItems").Find(&orders)
	for i, order := range orders {
		orders[i].Name = order.FullName()
		orders[i].Total = order.GetTotal()
	}
	return c.JSON(orders)

}
