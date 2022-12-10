package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/ambassador/src/models"
	"gorm.io/gorm"
)

type Product struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{db: db}
}

type ResponseProduct struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
}

func (h *Product) Products(c *fiber.Ctx) error {
	products := []models.Product{}
	h.db.Find(&products)

	response := []ResponseProduct{}
	for _, product := range products {

		response = append(response, ResponseProduct{
			ID:          product.ID,
			Title:       product.Title,
			Description: product.Description,
			Image:       product.Image,
			Price:       product.Price,
		})
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Product) CreateProducts(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	h.db.Create(&product)

	return c.JSON(product)
}

func (h *Product) GetProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var product models.Product
	product.ID = uint(id)
	h.db.First(&product, id)

	response := ResponseProduct{
		ID:          product.ID,
		Title:       product.Title,
		Description: product.Description,
		Image:       product.Image,
		Price:       product.Price,
	}
	return c.JSON(response)
}

func (h *Product) UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}
	product.ID = uint(id)

	if err := c.BodyParser(&product); err != nil {
		return err
	}
	h.db.Model(&product).Updates(&product)

	return c.JSON(product)
}

func (h *Product) DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}
	product.ID = uint(id)

	h.db.Delete(&product, id)

	return c.SendStatus(fiber.StatusOK)
}
