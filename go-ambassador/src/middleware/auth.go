package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		s := c.Get("Authorization")
		if s == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT"})

		}

		token := strings.TrimPrefix(s, "Bearer ")
		claims, err := validateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"message": "unauthenticated"})
		}
		payload := claims.(*jwt.RegisteredClaims)

		c.Locals("token", payload)
		c.Locals("sub", payload.Subject)
		return c.Next()
	}
}

func validateToken(token string) (jwt.Claims, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("MySignature"), nil
	})
	if err != nil || !t.Valid {
		return nil, err
	}

	return t.Claims, nil
}
