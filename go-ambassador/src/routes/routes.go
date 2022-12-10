package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/ambassador/src/controllers"
	"github.com/sing3demons/ambassador/src/database"
	"github.com/sing3demons/ambassador/src/middleware"
)

func Setup(app *fiber.App) {
	db := database.DB
	api := app.Group("api")

	admin := api.Group("admin")
	authController := controllers.NewAuthApplication(db)
	ambassadorController := controllers.NewUser(db)
	productController := controllers.NewProduct(db)
	admin.Post("register", authController.Register)
	admin.Post("login", authController.Login)
	{
		admin.Use(middleware.Protected())
		admin.Get("user", authController.GetUser)
		admin.Put("user", authController.UpdateInfo)
		admin.Patch("user/password", authController.UpdatePassword)
		admin.Get("ambassador", ambassadorController.Ambassador)
		admin.Get("products", productController.Products)
		admin.Get("product/:id", productController.GetProduct)
		admin.Put("product/:id", productController.UpdateProduct)
		admin.Delete("product/:id", productController.DeleteProduct)
	}

}
