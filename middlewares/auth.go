package middlewares

import (
	"context"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
	"github.com/gofiber/fiber/v2"
)

func Auth() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		session, err := database.SessionInstance().Get(c)

		if err != nil {
			panic(err)
		}

		firebase := database.FirebaseInstance()
		ctx := context.Background()
		IDToken := session.Get("token")

		tokenString := ""

		if IDToken != nil {
			tokenString = IDToken.(string)
		}

		if tokenString == "" {
			return c.Status(401).JSON(models.ErrorResponse{
				Error: "Not authorized",
			})
		}

		if tokenString == utils.GetEnv("DEV_TOKEN", "devtoken") {
			return c.Next()
		}

		client, cerr := firebase.Auth(ctx)
		if cerr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": cerr.Error(),
			})
		}

		_, terr := client.VerifyIDToken(ctx, tokenString)
		if terr != nil {
			return c.Status(401).JSON(models.ErrorResponse{
				Error: "Not authorized",
			})
		}

		return c.Next()
	}
}
