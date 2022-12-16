package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	models "zendx.io/P2P-Drive/models"
)

var userLogin models.LoginRequest

// -------------------------- Login User --------------------------\\

func UserLogin(c *fiber.Ctx) error {

	userLogin.Username = c.Query("Username")
	userLogin.UserPassword = c.Query("Password")

	fmt.Print(userLogin.Username)
	fmt.Print(userLogin.UserPassword)

	//return c.JSON(user)

	Database := Connection()

	token := Database.Login(&userLogin)

	if token == "Incorrect Password" {
		return c.JSON("ERROR")

	} else {
		return c.SendString(token)
	}
}
