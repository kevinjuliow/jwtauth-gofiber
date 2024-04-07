package main

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kevinjuliow/jwtauth-gofiber/database"
	"time"
)

func main() {
	database.OpenConnection()
	app := fiber.New()
	app.Use(logger.New())

	app.Post("/login", login)

	// Unprotected route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome Guest!")
	})

	var mySigningKey = []byte("secret")
	// Protected route
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: mySigningKey,
	}))

	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("Welcome User!")
	})

	app.Listen(":3000")
}

func login(ctx *fiber.Ctx) error {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body request

	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
		return err
	}

	if body.Username != "john" && body.Password != "doe" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"errpr": "Bad Credentials",
		})
		return err
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.JSON(fiber.Map{"token": t})
}
