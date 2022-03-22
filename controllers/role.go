package controllers

import (
	"strings"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

// Get Roles
// @Summary Get all roles
// @Description Get all roles
// @Tags Roles
// @Produce json
// @Param limit query number false "Limit"
// @Param offset query number false "Offset"
// @Param name query string false "Name"
// @Param ids query string false "IDs"
// @Success 200 {array} models.RoleResponse
// @Router /roles [get]
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
	var response []models.RoleResponse
	for _, role := range roles {
		response = append(response, models.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}
	return c.JSON(response)
}

// Get Role
// @Summary Get role
// @Description Get role
// @Tags Roles
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.RoleResponse
// @Router /roles/{id} [get]
func GetRole(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var role models.Role
	db.First(&role, id)
	return c.JSON(models.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	})
}

// Add Role
// @Summary Add role
// @Description Add role
// @Tags Roles
// @Accept json
// @Produce json
// @Param role body models.RolePayload true "Role"
// @Success 200 {object} models.RoleResponse
// @Router /roles [post]
func AddRole(c *fiber.Ctx) error {
	var role models.RolePayload
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
	roleModel := models.Role{
		Name:        role.Name,
		Description: role.Description,
	}
	status := db.Create(&roleModel)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})
	}
	return c.JSON(models.RoleResponse{
		ID:          roleModel.ID,
		Name:        roleModel.Name,
		Description: roleModel.Description,
		CreatedAt:   roleModel.CreatedAt,
		UpdatedAt:   roleModel.UpdatedAt,
	})
}

// Update Role
// @Summary Update role
// @Description Update role
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param role body models.RolePayload true "Role"
// @Success 200 {object} models.Message
// @Router /roles/{id} [put]
func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var role models.RolePayload
	if err := c.BodyParser(&role); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate role
	errors := utils.ValidateStruct(role)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	roleModel := models.Role{
		Name:        role.Name,
		Description: role.Description,
	}

	db := database.DBConn
	db.Model(&roleModel).Where("id = ?", id).Updates(role)
	return c.JSON(fiber.Map{
		"message": "Role updated",
	})
}

// Delete Role
// @Summary Delete role
// @Description Delete role
// @Tags Roles
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.Message
// @Failure 404 {object} models.Error
// @Router /roles/{id} [delete]
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
