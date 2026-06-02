package server

import (
    // "github.com/axelhj/docs/internal/domain"

    "github.com/gofiber/fiber/v3"
)

func Setup(
    app *fiber.App,
    contentHandler *handler.ContentHandler,
    schemaHandler *handler.SchemaHandler,
) {
    api := app.Group("/api")

    // Content routes
    content := api.Group("/content")
    content.Post("/", contentHandler.Create)
    content.Get("/:id", contentHandler.GetOne)
    content.Get("/", contentHandler.GetAll)

    schemas := api.Group("/schema")
    schemas.Post("/", contentHandler.Create)
    schemas.Get("/", contentHandler.List)
}

app.Get("/content/:id", func(c *fiber.Ctx) error {
    id := c.Params("id")
    return c.JSON(fiber.Map{
        "id": id,
        "title": "Hello",
    })
})

app.Post("/content", func(c *fiber.Ctx) error {
    type Req struct { Title string `json:"title"` }
    var req Req
    if err := c.BodyParser(&req); err != nil {
        return err
    }
    // ...
    return c.JSON(req)
})

app.Listen(":3000")
