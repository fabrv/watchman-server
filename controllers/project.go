package controllers

import (
	"strings"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

// Get Projects
// @Summary Get All Projects
// @Description Get all projects
// @Tags Projects
// @Produce json
// @Param limit query number false "Limit"
// @Param offset query number false "Offset"
// @Param name query string false "Name"
// @Param ids query string false "IDs"
// @Success 200 {array} models.ProjectResponse
// @Router /projects [get]
func GetProjects(c *fiber.Ctx) error {
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

	var projects []models.Project
	query := db.Limit(limit).Offset(offset)

	if ids[0] != "" {
		query = query.Where("id IN (?)", ids)
	}

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Find(&projects)

	var projectResponses []models.ProjectResponse
	for _, project := range projects {
		projectResponses = append(projectResponses, models.ProjectResponse{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		})
	}

	return c.JSON(projectResponses)
}

// Get Project
// @Summary Get Project
// @Description Get project
// @Tags Projects
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.ProjectResponse
// @Error 404 {object} models.ErrorResponse
// @Router /projects/{id} [get]
func GetProject(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var project models.Project
	db.First(&project, id)

	if project.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	return c.JSON(models.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	})
}

// Create Project
// @Summary Create Project
// @Description Create project
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body models.ProjectPayload true "Project"
// @Success 200 {object} models.ProjectResponse
// @Error 500 {object} models.ErrorResponse
// @Router /projects [post]
func AddProject(c *fiber.Ctx) error {
	var project models.ProjectPayload
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
	projectModel := models.Project{
		Name:        project.Name,
		Description: project.Description,
	}
	status := db.Create(&projectModel)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})
	}
	return c.JSON(project)
}

// Update Project
// @Summary Update Project
// @Description Update project
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param project body models.ProjectPayload true "Project"
// @Success 200 {object} models.ProjectResponse
// @Error 503 {object} models.ErrorResponse
// @Router /projects/{id} [put]
func UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.ProjectPayload
	if err := c.BodyParser(&project); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(project)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DBConn
	projectModel := models.Project{
		Name:        project.Name,
		Description: project.Description,
	}
	db.Model(&projectModel).Where("id = ?", id).Updates(projectModel)
	return c.JSON(fiber.Map{
		"message": "Project updated",
	})
}

// Delete Project
// @Summary Delete Project
// @Description Delete project
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.ProjectResponse
// @Error 404 {object} models.ErrorResponse
// @Router /projects/{id} [delete]
func DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var project models.Project
	db.First(&project, id)
	if project.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}
	db.Delete(&project)
	return c.JSON(fiber.Map{
		"message": "Project deleted",
	})
}
