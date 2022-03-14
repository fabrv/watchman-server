package controllers

import (
	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetTimeLogs(c *fiber.Ctx) error {
	db := database.DBConn
	var timeLogs []models.TimeLog

	limit := c.Query("limit")
	offset := c.Query("offset")

	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

	query := db.Limit(limit).Offset(offset)
	query.Preload("LogType").Preload("Project").Preload("Team").Find(&timeLogs)

	return c.JSON(timeLogs)
}

func GetTimeLog(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var timeLog models.TimeLog
	db.Preload("LogType").Preload("Project").Preload("Team").Preload("User").First(&timeLog, id)
	return c.JSON(timeLog)
}

func AddTimeLog(c *fiber.Ctx) error {
	var timeLogPayload models.TimeLogPayload
	if err := c.BodyParser(&timeLogPayload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(timeLogPayload)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	timeLog := models.TimeLog{
		UserID:      timeLogPayload.UserID,
		ProjectID:   timeLogPayload.ProjectID,
		TeamID:      timeLogPayload.TeamID,
		LogTypeID:   timeLogPayload.LogTypeID,
		StartTime:   timeLogPayload.StartTime,
		EndTime:     timeLogPayload.EndTime,
		Description: timeLogPayload.Description,
	}

	db := database.DBConn
	status := db.Create(&timeLog)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})
	}
	return c.JSON(timeLogPayload)
}

func UpdateTimeLog(c *fiber.Ctx) error {
	id := c.Params("id")
	var timeLog models.TimeLog
	if err := c.BodyParser(&timeLog); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	db := database.DBConn
	db.Model(&timeLog).Where("id = ?", id).Updates(timeLog)
	return c.JSON(fiber.Map{
		"message": "TimeLog updated",
	})
}

func DeleteTimeLog(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var timeLog models.TimeLog
	db.First(&timeLog, id)
	if timeLog.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "TimeLog not found",
		})
	}
	db.Delete(&timeLog)
	return c.JSON(fiber.Map{
		"message": "TimeLog deleted",
	})
}
