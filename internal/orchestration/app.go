package orchestration

import (
	"log"

	"github.com/axelhj/docs/internal/config"
	// "github.com/axelhj/docs/internal/repository"
	"github.com/axelhj/docs/internal/server"

	"github.com/gofiber/fiber/v3"
)

func NewApp(cfg *config.Config) {
	// TODO: Enable mocking of the couchster.
	// db := repository.NewCouchDB(cfg.CouchHost, cfg.CouchUser, cfg.CouchPassword)

	// couch := db.EnsureDB(cfg.DbName)
	// docService := service.NewDocumentService(couch)

	contentHandler := server.NewContentHandler()
	schemaHandler := server.NewSchemaHandler()

	// Todo: Use the server/service
	_, err := server.NewDocServer(cfg.CouchHost, cfg.CouchPort, cfg.CouchUser, cfg.CouchPassword)
	if err != nil {
		log.Fatalf("NewApp fatal %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "JustDocs",
	})

	// Perfect spot for adding the middleware etc.

	server.Setup(app, contentHandler, schemaHandler)

	log.Print(cfg.AppBindExpr)
	app.Listen(cfg.AppBindExpr)
}
