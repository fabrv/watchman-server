package main

import (
	"fmt"
	"log"

	"github.com/fabrv/watchman-server/routes"
	"github.com/fabrv/watchman-server/utils"
	"github.com/joho/godotenv"

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
	gerr := godotenv.Load()
	if gerr != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	// Init Database
	database.InitDatabase()
	defer func(DBConn *gorm.DB) {
		err := DBConn.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(database.DBConn)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Use(logger.New())

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
	}))

	routes.SetupRoutes(app)

	err := app.Listen(":" + utils.GetEnv("PORT", "3000"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
