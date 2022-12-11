package controllers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sing3demons/ambassador/src/models"
	"gorm.io/gorm"
)

type authController struct {
	db *gorm.DB
}

func NewAuthApplication(db *gorm.DB) *authController {
	return &authController{db: db}
}

func (h *authController) Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if data["password"] != data["password_confirm"] {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "password do no match",
		})
	}
	// password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		IsAmbassador: strings.Contains(c.Path(), "/api/ambassador"),
	}

	if err := user.EncryptPassword(data["password"]); err != nil {
		log.Println(err)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.db.Create(&user).Error; err != nil {
		log.Println()
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(user)
}

func (h *authController) Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user models.User
	// h.db.Where("email = ?", data["email"]).Find(&user)
	// if user.ID == 0 {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"message": "Invalid Credentials",
	// 	})
	// }
	if err := h.db.First(&user, "email = ?", data["email"]).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
	// 	log.Println(err.Error())
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"message": "Invalid Credentials",
	// 	})
	// }

	if err := user.CheckPassword(data["password"]); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}
	sub := strconv.Itoa(int(user.ID))

	// payload := jwt.StandardClaims{
	// 	Subject:   sub,
	// 	ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
	// }

	payload := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload).SignedString([]byte("MySignature"))
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 2),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func (h *authController) User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	t, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("MySignature"), nil
	})
	if err != nil || !t.Valid {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "unauthenticated"})
	}

	payload := t.Claims

	return c.Status(fiber.StatusOK).JSON(payload)
}

func (h *authController) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func (h *authController) findUserById(c *fiber.Ctx) (models.User, error) {
	sub := c.Locals("sub")
	var user models.User
	err := h.db.Where("id = ?", sub).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, err
	}
	return user, nil
}

func (h *authController) GetUser(c *fiber.Ctx) error {
	// sub := c.Locals("sub")
	// var user models.User
	// h.db.Where("id = ?", sub).First(&user)
	user, _ := h.findUserById(c)
	response := fiber.Map{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *authController) UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	sub := c.Locals("sub").(string)
	id, _ := strconv.Atoi(sub)

	user := models.User{
		ID:        uint(id),
		FirstName: data["first_name"],
		LastName:  data["last_name"],
	}

	h.db.Model(&user).Updates(&user)

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *authController) UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if data["password"] != data["password_confirm"] {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "password do no match",
		})
	}

	sub := c.Locals("sub").(string)
	id, _ := strconv.Atoi(sub)

	user := models.User{
		ID: uint(id),
	}

	user.EncryptPassword(data["password"])

	h.db.Model(&user).Updates(&user)

	return c.Status(fiber.StatusOK).JSON(user)
}
