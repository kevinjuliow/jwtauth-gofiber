package main

import (
	jwtware "github.com/gofiber/contrib/jwt"
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
	private.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: "HS256",
			Key:    []byte("secret"),
		},
	}))
	private.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome User Protected Routes!")
	})
	app.Listen(":8000")
}
