package main

import (
	"fmt"
	"os"

	"github.com/fabrv/watchman-server/controllers"
	"github.com/fabrv/watchman-server/database"
	"github.com/fabrv/watchman-server/models"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/jinzhu/gorm"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open("postgres", getenv("DATABASE_URL", "host=localhost port=5432 user=postgres dbname=watchman password=password sslmode=disable"))

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")

	database.DBConn.AutoMigrate(&models.Book{})
	fmt.Println("Database migrated")
}

func setupRoutes(app *fiber.App) {
	var baseRoute = "/api/v1"
	app.Get(baseRoute+"/books", controllers.GetBooks)
	app.Get(baseRoute+"/books/:id", controllers.GetBook)
	app.Post(baseRoute+"/books", controllers.AddBook)
	app.Put(baseRoute+"/books/:id", controllers.UpdateBook)
	app.Delete(baseRoute+"/books/:id", controllers.DeleteBook)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    404,
			"message": "404: Not Found",
		})
	})
}

func main() {
	app := fiber.New()
	initDatabase()
	defer database.DBConn.Close()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // comma string format e.g. "localhost, nikschaefer.tech"
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(logger.New())
	setupRoutes(app)

	app.Listen(":" + getenv("PORT", "3000"))
}
