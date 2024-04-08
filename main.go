package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kevinjuliow/jwtauth-gofiber/controller"
	"github.com/kevinjuliow/jwtauth-gofiber/database"
)

func main() {
	database.OpenConnection()
	app := fiber.New()
	app.Use(logger.New())

	app.Post("/signup", controller.Signup)
	app.Post("/login", controller.Login)

	// Unprotected route
	public := app.Group("/public")
	public.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome Guest!")
	})
	//protected route
	private := app.Group("/private")
	private.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome User!")
	})

	app.Listen(":8000")
}
