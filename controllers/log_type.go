package controllers

import (
	"strings"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

// GetLogTypes
// @Summary Get all LogTypes
// @Description Get all LogTypes
// @Tags LogTypes
// @Accept  json
// @Produce  json
// @Success 200 {array} models.LogTypeResponse
// @Router /log-types [get]
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

// GetLogType
// @Summary Get one logType
// @Description Get one logType
// @Tags LogTypes
// @Accept  json
// @Produce  json
// @Param id path string true "LogType ID"
// @Success 200 {object} models.LogTypeResponse
// @Router /log-types/{id} [get]
func GetLogType(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var logType models.LogType
	db.First(&logType, id)
	return c.JSON(logType)
}

// CreateLogType
// @Summary Create a new logType
// @Description Create a new logType
// @Tags LogTypes
// @Accept  json
// @Produce  json
// @Param log_type body models.LogTypePayload true "LogType"
// @Success 200 {object} models.LogTypeResponse
// @Router /log-types [post]
func AddLogType(c *fiber.Ctx) error {
	var logType models.LogTypePayload
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
	payLoadModel := models.LogType{
		Name:        logType.Name,
		Description: logType.Description,
	}

	status := db.Create(&payLoadModel)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})

	}
	return c.JSON(models.LogTypeResponse{
		ID:          payLoadModel.ID,
		CreatedAt:   payLoadModel.CreatedAt,
		UpdatedAt:   payLoadModel.UpdatedAt,
		Name:        payLoadModel.Name,
		Description: payLoadModel.Description,
	})
}

// UpdateLogType
// @Summary Update a logType
// @Description Update a logType
// @Tags LogTypes
// @Accept  json
// @Produce  json
// @Param id path string true "LogType ID"
// @Param log_type body models.LogTypePayload true "LogType"
// @Success 200 {object} models.LogTypeResponse
// @Router /log-types/{id} [put]
func UpdateLogType(c *fiber.Ctx) error {
	id := c.Params("id")
	var logType models.LogTypePayload
	if err := c.BodyParser(&logType); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	db := database.DBConn
	logTypeModel := models.LogType{
		Name:        logType.Name,
		Description: logType.Description,
	}

	db.Model(&logTypeModel).Where("id = ?", id).Updates(logTypeModel)
	return c.JSON(fiber.Map{
		"message": "LogType updated",
	})
}

// Delete LogType
// @Summary Delete a logType
// @Description Delete a logType
// @Tags LogTypes
// @Accept  json
// @Produce  json
// @Param id path string true "LogType ID"
// @Success 200 {object} models.LogTypeResponse
// @Router /log-types/{id} [delete]
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
