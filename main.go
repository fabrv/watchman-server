package main

import (
	"fmt"

	"github.com/fabrv/watchman-server/routes"
	"github.com/fabrv/watchman-server/utils"

	"github.com/fabrv/watchman-server/database"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/jinzhu/gorm"

	swagger "github.com/arsmn/fiber-swagger/v2"
	_ "github.com/fabrv/watchman-server/docs"
)

// @BasePath /api/v1
func main() {
	app := fiber.New()
	database.InitDatabase()
	defer func(DBConn *gorm.DB) {
		err := DBConn.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(database.DBConn)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(logger.New())

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

	routes.SetupRoutes(app)

	err := app.Listen(":" + utils.GetEnv("PORT", "3000"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
