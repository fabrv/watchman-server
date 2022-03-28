package controllers

import (
	"fmt"
	"strings"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

// Get Users
// @Summary Get all users
// @Description Get all users
// @Tags Users
// @Produce json
// @Param limit query number false "Limit"
// @Param offset query number false "Offset"
// @Param name query string false "Name"
// @Param ids query string false "IDs"
// @Success 200 {array} models.UserPayload
// @Router /users [get]
func GetUsers(c *fiber.Ctx) error {
	db := database.DBConn
	var users []models.User

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

	query := db.Limit(limit).Offset(offset)

	if ids[0] != "" {
		query = query.Where("id IN (?)", ids)
	}

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// Preload roles and select only id, name, role_id and name
	query.Preload("Role").Select("id, name, email, role_id").Find(&users)

	var usersResponse []models.UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, models.UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			RoleID: user.RoleID,
			Role: models.RoleResponse{
				ID:          user.Role.ID,
				Name:        user.Role.Name,
				Description: user.Role.Description,
			},
		})
	}

	return c.JSON(usersResponse)
}

// Get User
// @Summary Get user
// @Description Get user
// @Tags Users
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.UserResponse
// @Router /users/{id} [get]
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var user models.User
	db.Preload("Role").Select("id, name, email, role_id").First(&user, id)

	if user.ID == 0 {
		return c.Status(404).JSON(models.ErrorResponse{
			Error: "User not found",
		})
	}

	return c.JSON(models.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		RoleID: user.RoleID,
		Role: models.RoleResponse{
			ID:          user.Role.ID,
			Name:        user.Role.Name,
			Description: user.Role.Description,
		},
	})
}

// Create User
// @Summary Create user
// @Description Create user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserPayload true "User"
// @Success 200 {object} models.UserResponse
// @Router /users [post]
func AddUser(c *fiber.Ctx) error {
	var userPayload models.UserPayload
	if err := c.BodyParser(&userPayload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(userPayload)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DBConn
	var role models.Role
	db.First(&role, userPayload.RoleID)

	user := models.User{
		Name:     userPayload.Name,
		Email:    userPayload.Email,
		Password: userPayload.Password,
		Role:     role,
	}

	fmt.Println(role)

	status := db.Create(&user)

	if status.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": status.Error.Error(),
		})
	}
	return c.JSON(models.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		RoleID: user.RoleID,
	})
}

// Update User
// @Summary Update user
// @Description Update user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param user body models.UserPayload true "User"
// @Success 200 {object} models.MessageResponse
// @Router /users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.UserPayload
	if err := c.BodyParser(&user); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := utils.ValidateStruct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DBConn
	userModel := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	db.Model(&userModel).Where("id = ?", id).Updates(userModel)
	return c.JSON(fiber.Map{
		"message": "User updated",
	})
}

// Delete User
// @Summary Delete user
// @Description Delete user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.MessageResponse
// @Router /users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var user models.User
	db.First(&user, id)
	if user.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	db.Delete(&user)
	return c.JSON(fiber.Map{
		"message": "User deleted",
	})
}
