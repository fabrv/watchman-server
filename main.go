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
)

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
	routes.SetupRoutes(app)

	err := app.Listen(":" + utils.GetEnv("PORT", "3000"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
