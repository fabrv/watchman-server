package controllers

import (
	"strings"
	"time"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

// GetTimeLogs returns all TimeLogs
// @Summary Get all TimeLogs
// @Description Get all TimeLogs
// @Tags TimeLog
// @Accept json
// @Produce json
// @Success 200 {array} models.TimeLogResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /time-logs [get]
func GetTimeLogs(c *fiber.Ctx) error {
	db := database.DBConn
	var timeLogs []models.TimeLog

	limit := c.Query("limit")
	offset := c.Query("offset")
	userIds := strings.Split(c.Query("user_ids"), ",")

	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}

	query := db.Limit(limit).Offset(offset)

	if userIds[0] != "" {
		query = query.Where("user_id IN (?)", userIds)
	}

	query.Find(&timeLogs)

	var timeLogsResponse []models.TimeLogResponse
	for _, timeLog := range timeLogs {
		timeLogsResponse = append(timeLogsResponse, models.TimeLogResponse{
			ID:          timeLog.ID,
			UserId:      timeLog.UserID,
			LogTypeId:   timeLog.LogTypeID,
			ProjectId:   timeLog.ProjectID,
			TeamId:      timeLog.TeamID,
			StartTime:   timeLog.StartTime,
			EndTime:     timeLog.EndTime,
			Description: timeLog.Description,
		})
	}

	return c.JSON(timeLogsResponse)
}

// GetTimeLog returns a TimeLog
// @Summary Get a TimeLog
// @Description Get a TimeLog
// @Tags TimeLog
// @Accept json
// @Produce json
// @Param id path string true "TimeLog ID"
// @Success 200 {object} models.TimeLogResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /time-logs/{id} [get]
func GetTimeLog(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var timeLog models.TimeLog
	db.Preload("LogType").Preload("Project").Preload("Team").Preload("User").First(&timeLog, id)
	return c.JSON(timeLog)
}

// CreateTimeLog creates a new TimeLog
// @Summary Create a new TimeLog
// @Description Create a new TimeLog
// @Tags TimeLog
// @Accept json
// @Produce json
// @Param time_log body models.TimeLogPayload true "TimeLog"
// @Success 200 {object} models.TimeLogResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /time-logs [post]
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

// UpdateTimeLog updates a TimeLog
// @Summary Update a TimeLog
// @Description Update a TimeLog
// @Tags TimeLog
// @Accept json
// @Produce json
// @Param id path string true "TimeLog ID"
// @Param time_log body models.TimeLogPayload true "TimeLog"
// @Success 200 {object} models.MessageResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /time-logs/{id} [put]
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

// DeleteTimeLog deletes a TimeLog
// @Summary Delete a TimeLog
// @Description Delete a TimeLog
// @Tags TimeLog
// @Accept json
// @Produce json
// @Param id path string true "TimeLog ID"
// @Success 200 {object} models.MessageResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /time-logs/{id} [delete]
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

// FinishTimeLog finishes a TimeLog
// @Summary Finish a TimeLog
// @Description Sets the end time of a TimeLog to the current time
// @Tags TimeLog
// @Accept json
// @Produce json
// @Param id path string true "TimeLog ID"
// @Success 200 {object} models.MessageResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /time-logs/{id}/finish [post]
func FinishTimeLog(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var timeLog models.TimeLog
	db.First(&timeLog, id)
	if timeLog.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "TimeLog not found",
		})
	}
	timeLog.EndTime = time.Now()
	db.Save(&timeLog)
	return c.JSON(fiber.Map{
		"message": "TimeLog finished",
	})
}
