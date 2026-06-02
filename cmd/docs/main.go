package main

import (
	"log"

	"github.com/axelhj/docs/internal/config"
	"github.com/axelhj/docs/internal/repository"
	"github.com/axelhj/docs/internal/server"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.Load()

	db := repository.NewCouchDB(cfg.CouchHost, cfg.CouchUser, cfg.CouchPassword)

	couch := db.EnsureDB(cfg.DbName)           // your couch client
	// docService := service.NewDocumentService(couch)

	contentHandler := server.NewContentHandler(couch)
	schemaHandler := handler.NewSchemaHandler(/*...*/) // etc.

	if srv, err := server.NewDocServer(cfg.CouchHost, cfg.CouchUser, cfg.CouchPassword); err != nil {
		log.Fatal(err)
	}

	if err := srv.Start(cfg.Server.Port); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	router.Setup(app, contentHandler, schemaHandler)

	app.Listen(":3000")
}
