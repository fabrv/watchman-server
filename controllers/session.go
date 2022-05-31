package controllers

import (
	"context"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateSession(c *fiber.Ctx) error {
	store := database.SessionInstance()
	firebase := database.FirebaseInstance()
	ctx := context.Background()

	payload := models.SessionPayload{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if utils.ValidateStruct(payload) != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Token is required",
		})
	}

	IDToken := payload.Token

	session, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	if IDToken == utils.GetEnv("DEV_TOKEN", "devtoken") {
		session.Set("user", "dev")
		session.Save()
		return c.SendString("Session created")
	}

	client, cerr := firebase.Auth(ctx)
	if cerr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": cerr.Error(),
		})
	}

	tokenInfo, terr := client.VerifyIDToken(ctx, IDToken)
	if terr != nil {
		return c.Status(401).JSON(models.ErrorResponse{
			Error: "Not authorized",
		})
	}

	session.Set("user", tokenInfo.Claims)
	session.Set("token", IDToken)
	saveErr := session.Save()
	if saveErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": saveErr.Error(),
		})
	}
	return c.JSON(models.MessageResponse{
		Message: "Session created",
	})
}

func DeleteSession(c *fiber.Ctx) error {
	store := database.SessionInstance()
	session, err := store.Get(c)
	if err != nil {
		panic(err)
	}
	session.Destroy()
	session.Save()
	return c.JSON(models.MessageResponse{
		Message: "Session deleted",
	})
}
