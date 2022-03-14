package controllers

import (
	"fmt"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	db := database.DBConn
	var users []models.User
	// Preload roles and select only id, name, role_id and name
	db.Preload("Role").Select("id, name, email, role_id").Find(&users)

	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var user models.User
	db.Preload("Role").Select("id, name, email, role_id").First(&user, id)
	return c.JSON(user)
}

func AddUser(c *fiber.Ctx) error {
	var userPayload models.UserPayload
	if err := c.BodyParser(&userPayload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(userPayload)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DBConn
	var role models.Role
	db.First(&role, userPayload.RoleID)

	user := models.User{
		Name:     userPayload.Name,
		Email:    userPayload.Email,
		Password: userPayload.Password,
		Role:     role,
	}

	fmt.Println(role)

	status := db.Create(&user)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})
	}
	return c.JSON(userPayload)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	db := database.DBConn
	db.Model(&user).Where("id = ?", id).Updates(user)
	return c.JSON(fiber.Map{
		"message": "User updated",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var user models.User
	db.First(&user, id)
	if user.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	db.Delete(&user)
	return c.JSON(fiber.Map{
		"message": "User deleted",
	})
}
