package controllers

import (
	"strings"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

// Get Teams
// @Summary Get all teams
// @Description Get all teams
// @Tags Teams
// @Produce json
// @Param limit query number false "Limit"
// @Param offset query number false "Offset"
// @Param name query string false "Name"
// @Param ids query string false "IDs"
// @Success 200 {array} models.TeamResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /teams [get]
func GetTeams(c *fiber.Ctx) error {
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

	var teams []models.Team
	query := db.Limit(limit).Offset(offset)

	if ids[0] != "" {
		query = query.Where("id IN (?)", ids)
	}

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	query.Find(&teams)

	var response []models.TeamResponse
	for _, team := range teams {
		response = append(response, models.TeamResponse{
			ID:          team.ID,
			Name:        team.Name,
			Description: team.Description,
			CreatedAt:   team.CreatedAt,
			UpdatedAt:   team.UpdatedAt,
		})
	}

	return c.JSON(response)
}

// Get Team
// @Summary Get team
// @Description Get team
// @Tags Teams
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} models.TeamResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /teams/{id} [get]
func GetTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var team models.Team
	db.First(&team, id)

	if team.ID == 0 {
		return c.Status(fiber.ErrNotFound.Code).JSON(models.ErrorResponse{
			Error: "Team not found",
		})
	}

	return c.JSON(models.TeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	})
}

// Create Team
// @Summary Create team
// @Description Create team
// @Tags Teams
// @Accept json
// @Produce json
// @Param team body models.TeamPayload true "Team"
// @Success 200 {object} models.TeamResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /teams [post]
func AddTeam(c *fiber.Ctx) error {
	var team models.TeamPayload
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
	teamModel := models.Team{
		Name:        team.Name,
		Description: team.Description,
	}
	status := db.Create(&teamModel)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})
	}
	return c.JSON(models.TeamResponse{
		ID:          teamModel.ID,
		Name:        teamModel.Name,
		Description: teamModel.Description,
		CreatedAt:   teamModel.CreatedAt,
		UpdatedAt:   teamModel.UpdatedAt,
	})
}

// Update Team
// @Summary Update team
// @Description Update team
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Param team body models.TeamPayload true "Team"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /teams/{id} [put]
func UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	var team models.TeamPayload
	if err := c.BodyParser(&team); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(team)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DBConn
	teamModel := models.Team{
		Name:        team.Name,
		Description: team.Description,
	}

	db.Model(&teamModel).Where("id = ?", id).Updates(teamModel)
	return c.JSON(fiber.Map{
		"message": "Team updated",
	})
}

// Delete Team
// @Summary Delete team
// @Description Delete team
// @Tags Teams
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /teams/{id} [delete]
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
