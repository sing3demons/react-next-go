package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/ambassador/src/controllers"
	"github.com/sing3demons/ambassador/src/database"
	"github.com/sing3demons/ambassador/src/middleware"
)

func Setup(app *fiber.App) {
	db := database.GetDB()
	api := app.Group("api")

	admin := api.Group("admin")
	authController := controllers.NewAuthApplication(db)

	admin.Post("register", authController.Register)
	admin.Post("login", authController.Login)
	{
		admin.Use(middleware.Protected())
		admin.Get("user", authController.GetUser)
		admin.Put("user", authController.UpdateInfo)
		admin.Patch("user/password", authController.UpdatePassword)

	}

	linkController := controllers.NewLink(db)
	{
		admin.Use(middleware.IsAuthenticated)
		admin.Get("user/:id/links", linkController.Link)
	}

	orderController := controllers.NewOrder(db)
	{
		admin.Use(middleware.IsAuthenticated)
		admin.Get("orders", orderController.Order)
	}

	ambassadorController := controllers.NewUser(db)
	ambassador := admin.Group("ambassador")
	{
		ambassador.Use(middleware.Protected())
		ambassador.Get("", ambassadorController.Ambassador)
	}

	productController := controllers.NewProduct(db)
	productGroup := admin.Group("products")
	{
		productGroup.Use(middleware.Protected())
		productGroup.Get("", productController.Products)
		productGroup.Get("/:id", productController.GetProduct)
		productGroup.Put("/:id", productController.UpdateProduct)
		productGroup.Delete("/:id", productController.DeleteProduct)
	}

	ambassadorGroup := api.Group("ambassador")
	ambassadorGroup.Post("register", authController.Register)
	ambassadorGroup.Post("login", authController.Login)
	{
		ambassadorGroup.Use(middleware.IsAuthenticated)
		ambassadorGroup.Get("user", authController.User)

	}
}
