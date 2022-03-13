package controllers

import (
	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetProjects(c *fiber.Ctx) error {
	db := database.DBConn
	var projects []models.Project
	db.Find(&projects)
	return c.JSON(projects)
}

func GetProject(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var project models.Project
	db.First(&project, id)
	return c.JSON(project)
}

func AddProject(c *fiber.Ctx) error {
	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(project)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DBConn
	status := db.Create(&project)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})
	}
	return c.JSON(project)
}

func UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	db := database.DBConn
	db.Model(&project).Where("id = ?", id).Updates(project)
	return c.JSON(fiber.Map{
		"message": "Project updated",
	})
}

func DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var project models.Project
	db.First(&project, id)
	if project.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Project not found",
		})
	}
	db.Delete(&project)
	return c.JSON(fiber.Map{
		"message": "Project deleted",
	})
}
