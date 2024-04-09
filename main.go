package main

import (
	"errors"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kevinjuliow/jwtauth-gofiber/controller"
	"github.com/kevinjuliow/jwtauth-gofiber/database"
	"github.com/kevinjuliow/jwtauth-gofiber/models"
	"gorm.io/gorm"
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
			Key: []byte("secret"),
		},
	}))
	private.Get("/", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userId := claims["user_id"].(float64)

		//find the user in the database
		var findUser models.Users
		if err := database.DB.Where("id = ?", userId).First(&findUser).Error; err != nil {
			// Handle the case where user is not found or an error occurs
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("User not found")
			}
			// Handle other errors
			return err // You might want to handle other types of errors differently
		}

		return c.JSON(fiber.Map{
			"Message": "Welcome",
			"data":    findUser,
		})
	})
	app.Listen(":8000")
}
