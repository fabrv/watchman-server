package controllers

import (
	"strings"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetLogTypes(c *fiber.Ctx) error {
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

	var logTypes []models.LogType
	query := db.Limit(limit).Offset(offset)

	if ids[0] != "" {
		query = query.Where("id IN (?)", ids)
	}

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Find(&logTypes)
	return c.JSON(logTypes)
}

func GetLogType(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var logType models.LogType
	db.First(&logType, id)
	return c.JSON(logType)
}

func AddLogType(c *fiber.Ctx) error {
	var logType models.LogType
	if err := c.BodyParser(&logType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(logType)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DBConn
	status := db.Create(&logType)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})
	}
	return c.JSON(logType)
}

func UpdateLogType(c *fiber.Ctx) error {
	id := c.Params("id")
	var logType models.LogType
	if err := c.BodyParser(&logType); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	db := database.DBConn
	db.Model(&logType).Where("id = ?", id).Updates(logType)
	return c.JSON(fiber.Map{
		"message": "LogType updated",
	})
}

func DeleteLogType(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var logType models.LogType
	db.First(&logType, id)
	if logType.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "LogType not found",
		})
	}
	db.Delete(&logType)
	return c.JSON(fiber.Map{
		"message": "LogType deleted",
	})
}
