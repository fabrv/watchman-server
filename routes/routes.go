package routes

import (
	"github.com/fabrv/watchman-server/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	var baseRoute = "/api/v1"

	// Roles routes
	app.Get(baseRoute+"/roles", controllers.GetRoles)
	app.Get(baseRoute+"/roles/:id", controllers.GetRole)
	app.Post(baseRoute+"/roles", controllers.AddRole)
	app.Put(baseRoute+"/roles/:id", controllers.UpdateRole)
	app.Delete(baseRoute+"/roles/:id", controllers.DeleteRole)

	// Project routes
	app.Get(baseRoute+"/projects", controllers.GetProjects)
	app.Get(baseRoute+"/projects/:id", controllers.GetProject)
	app.Post(baseRoute+"/projects", controllers.AddProject)
	app.Put(baseRoute+"/projects/:id", controllers.UpdateProject)
	app.Delete(baseRoute+"/projects/:id", controllers.DeleteProject)

	// Team routes
	app.Get(baseRoute+"/teams", controllers.GetTeams)
	app.Get(baseRoute+"/teams/:id", controllers.GetTeam)
	app.Post(baseRoute+"/teams", controllers.AddTeam)
	app.Put(baseRoute+"/teams/:id", controllers.UpdateTeam)
	app.Delete(baseRoute+"/teams/:id", controllers.DeleteTeam)

	// Log Type routes
	app.Get(baseRoute+"/log-types", controllers.GetLogTypes)
	app.Get(baseRoute+"/log-types/:id", controllers.GetLogType)
	app.Post(baseRoute+"/log-types", controllers.AddLogType)
	app.Put(baseRoute+"/log-types/:id", controllers.UpdateLogType)
	app.Delete(baseRoute+"/log-types/:id", controllers.DeleteLogType)

	// User routes
	app.Get(baseRoute+"/users", controllers.GetUsers)
	app.Get(baseRoute+"/users/:id", controllers.GetUser)
	app.Post(baseRoute+"/users", controllers.AddUser)
	app.Put(baseRoute+"/users/:id", controllers.UpdateUser)
	app.Delete(baseRoute+"/users/:id", controllers.DeleteUser)

	// Time Log routes
	app.Get(baseRoute+"/time-logs", controllers.GetTimeLogs)
	app.Get(baseRoute+"/time-logs/:id", controllers.GetTimeLog)
	app.Post(baseRoute+"/time-logs", controllers.AddTimeLog)
	app.Put(baseRoute+"/time-logs/:id", controllers.UpdateTimeLog)
	app.Delete(baseRoute+"/time-logs/:id", controllers.DeleteTimeLog)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    404,
			"message": "404: Not Found",
		})
	})
}
