package routes

import (
	"github.com/fabrv/watchman-server/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	var baseRoute = "/api/v1"
	app.Get(baseRoute+"/roles", controllers.GetRoles)
	app.Get(baseRoute+"/roles/:id", controllers.GetRole)
	app.Post(baseRoute+"/roles", controllers.AddRole)
	app.Put(baseRoute+"/roles/:id", controllers.UpdateRole)
	app.Delete(baseRoute+"/roles/:id", controllers.DeleteRole)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    404,
			"message": "404: Not Found",
		})
	})
}
