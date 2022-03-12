package main

import (
	"fmt"
	"os"

	"github.com/fabrv/watchman-server/routes"

	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/jinzhu/gorm"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("postgres", getEnv("DATABASE_URL", "host=localhost port=5432 user=postgres dbname=watchman password=password sslmode=disable"))

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")

	database.DBConn.AutoMigrate(&models.User{}, &models.Project{}, &models.Team{}, &models.Role{})
	fmt.Println("Database migrated")
}

func main() {
	app := fiber.New()
	initDatabase()
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

	err := app.Listen(":" + getEnv("PORT", "3000"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
