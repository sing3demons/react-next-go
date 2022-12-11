package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/ambassador/src/models"
	"gorm.io/gorm"
)

type Link struct {
	db *gorm.DB
}

func NewLink(db *gorm.DB) *Link {
	return &Link{db}
}

func (h *Link) Link(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var links []models.Link

	h.db.Where("user_id = ?", id).Find(&links)
	for i, link := range links {
		var orders []models.Order
		h.db.Where("code = ? and complete = true", link.Code).Find(&orders)
		links[i].Orders = orders

	}

	return c.Status(fiber.StatusOK).JSON(links)
}
