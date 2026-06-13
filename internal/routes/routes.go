package routes

import (
	"user-api/internal/handler"
	"user-api/internal/middleware"

	_ "user-api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/swaggo/swag"
)

func Setup(app *fiber.App, userHandler *handler.UserHandler) {
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())

	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		doc, err := swag.ReadDoc()
		if err != nil {
			return c.Status(500).SendString("swagger doc not found")
		}
		c.Set("Content-Type", "application/json")
		return c.SendString(doc)
	})

	app.Get("/swagger", func(c *fiber.Ctx) error {
		html := `<!DOCTYPE html>
<html>
  <head>
    <title>User API — Swagger UI</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
      window.onload = function() {
        SwaggerUIBundle({
          url: "/swagger/doc.json",
          dom_id: '#swagger-ui',
          presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
          layout: "StandaloneLayout"
        })
      }
    </script>
  </body>
</html>`
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	users := app.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.ListUsers)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
