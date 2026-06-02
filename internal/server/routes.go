package server

import (
	// "github.com/axelhj/docs/internal/domain"

	"github.com/gofiber/fiber/v3"
)

func Setup(
	app *fiber.App,
	contentHandler *ContentHandler,
	schemaHandler *SchemaHandler,
) {
	api := app.Group("/api")

	// Content routes
	content := api.Group("/content")
	content.Post("/", contentHandler.Create)
	content.Get("/:id", contentHandler.GetOne)
	content.Get("/", contentHandler.GetAll)

	schemas := api.Group("/schema")
	schemas.Post("/", schemaHandler.Create)
	schemas.Get("/", schemaHandler.GetAll)
}

// app.Get("/content/:id", func(c *fiber.Ctx) error {
// 	id := c.Params("id")
//	return c.JSON(fiber.Map{
//		"id": id,
//		"title": "Hello",
//	})
// })
//
// app.Post("/content", func(c *fiber.Ctx) error {
//	type Req struct { Title string `json:"title"` }
//	var req Req
//	if err := c.BodyParser(&req); err != nil {
//		return err
//	}
//	// ...
//	return c.JSON(req)
// })
//
