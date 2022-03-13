package controllers

import (
	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetTeams(c *fiber.Ctx) error {
	db := database.DBConn
	var teams []models.Team
	db.Find(&teams)
	return c.JSON(teams)
}

func GetTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var team models.Team
	db.First(&team, id)
	return c.JSON(team)
}

func AddTeam(c *fiber.Ctx) error {
	var team models.Team
	if err := c.BodyParser(&team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(team)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DBConn
	status := db.Create(&team)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})
	}
	return c.JSON(team)
}

func UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	var team models.Team
	if err := c.BodyParser(&team); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	db := database.DBConn
	db.Model(&team).Where("id = ?", id).Updates(team)
	return c.JSON(fiber.Map{
		"message": "Team updated",
	})
}

func DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var team models.Team
	db.First(&team, id)
	if team.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Team not found",
		})
	}
	db.Delete(&team)
	return c.JSON(fiber.Map{
		"message": "Team deleted",
	})
}