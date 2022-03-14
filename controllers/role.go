package controllers

import (
	"strings"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetRoles(c *fiber.Ctx) error {
	db := database.DBConn

	limit := c.Query("limit")
	offset := c.Query("offset")
	name := c.Query("name")
	ids := strings.Split(c.Query("ids"), ",")

	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

	var roles []models.Role
	query := db.Limit(limit).Offset(offset)

	if ids[0] != "" {
		query = query.Where("id IN (?)", ids)
	}

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Find(&roles)
	return c.JSON(roles)
}

func GetRole(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var role models.Role
	db.First(&role, id)
	return c.JSON(role)
}

func AddRole(c *fiber.Ctx) error {
	var role models.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(role)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DBConn
	status := db.Create(&role)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})
	}
	return c.JSON(role)
}

func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var role models.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	db := database.DBConn
	db.Model(&role).Where("id = ?", id).Updates(role)
	return c.JSON(fiber.Map{
		"message": "Role updated",
	})
}

func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var role models.Role
	db.First(&role, id)
	if role.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Role not found",
		})
	}
	db.Delete(&role)
	return c.JSON(fiber.Map{
		"message": "Role deleted",
	})
}
