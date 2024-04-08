package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kevinjuliow/jwtauth-gofiber/database"
	"github.com/kevinjuliow/jwtauth-gofiber/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func Signup(ctx *fiber.Ctx) error {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body request

	if err := ctx.BodyParser(&body); err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
		return err
	}

	if body.Username == "" || body.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Signup Credentials")
	}

	//Hash password to db
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if errHash != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errHash,
		})
		return errHash
	}
	user := models.Users{
		Username: body.Username,
		Password: string(hashedPassword),
	}

	////Save to database
	result := database.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	//generate jwt tokens
	t, exp, err := createJWTToken(&user)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"token": t,
		"exp":   exp,
	})
}

func Login(ctx *fiber.Ctx) error {
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

	if body.Username == "" || body.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Login Credentials")
	}

	//find the username in the db
	var user models.Users
	result := database.DB.Where("username = ?", body.Username).First(&user)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}

	// Compare the provided password with the stored password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		// Password does not match
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid password")
	}

	//create jwt token
	t, exp, err := createJWTToken(&user)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"token": t,
		"exp":   exp,
	})
}

func createJWTToken(user *models.Users) (string, int64, error) {
	expired := time.Now().Add(time.Hour * 1).Unix()
	// Create the Claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"admin":   true,
		"exp":     expired,
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", 0, err
	}
	return t, expired, nil
}
