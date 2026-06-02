package server

import (
	"github.com/axelhj/docs/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type ContentHandler struct {
	// Stuff that the handler needs, like the handler
	// private members etc.
	// service *service.DocumentService
}

type SchemaHandler struct {
}

// func NewContentHandler(serviceParam *service.DocumentService) *ContentHandler {
// 	return &ContentHandler{service: serviceParam}
// }

func NewContentHandler() *ContentHandler {
	return &ContentHandler{}
}

func NewSchemaHandler() *SchemaHandler {
	return &SchemaHandler{}
}

func (h *ContentHandler) Create(c fiber.Ctx) error {
	var doc domain.Document
	if err := c.JSON(&doc); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

//TODO: Implement a service or go directly to the repository.
	// saved, err := h.service.Save(c.Context(), &doc)
	// if err != nil {
		// return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	// }

	// return c.Status(291).JSON(saved)
	return nil
}

func (h *ContentHandler) GetOne(c fiber.Ctx) error {
	// id := c.Params("id")
	// doc, err := h.service.GetById(c.Context(), id)
	// if err != nil {
	// 	return c.Status(404).JSON(fiber.Map{"error": "Not found. Did you submit the wrong ID?"})
	// }
	// return c.JSON(doc)
	return nil
}

func (h *ContentHandler) GetAll(c fiber.Ctx) error {
	// id := c.Params("id")
	// doc, err := h.service.GetById(c.Context(), id)
	// if err != nil {
	// 	return c.Status(404).JSON(fiber.Map{"error": "Not found. Did you submit the wrong ID?"})
	// }
	// return c.JSON(doc)
	return nil
}

func (h *SchemaHandler) Create(c fiber.Ctx) error {
	var doc domain.Document
	if err := c.JSON(&doc); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

//TODO: Implement a service or go directly to the repository.
	// saved, err := h.service.Save(c.Context(), &doc)
	// if err != nil {
		// return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	// }

	// return c.Status(291).JSON(saved)
	return nil
}

func (h *SchemaHandler) GetOne(c fiber.Ctx) error {
	// id := c.Params("id")
	// doc, err := h.service.GetById(c.Context(), id)
	// if err != nil {
	// 	return c.Status(404).JSON(fiber.Map{"error": "Not found. Did you submit the wrong ID?"})
	// }
	// return c.JSON(doc)
	return nil
}

func (h *SchemaHandler) GetAll(c fiber.Ctx) error {
	// id := c.Params("id")
	// doc, err := h.service.GetById(c.Context(), id)
	// if err != nil {
	// 	return c.Status(404).JSON(fiber.Map{"error": "Not found. Did you submit the wrong ID?"})
	// }
	// return c.JSON(doc)
	return nil
}

