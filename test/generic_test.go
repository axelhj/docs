package test

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"

	"github.com/axelhj/docs/internal/config"
	// "github.com/axelhj/docs/internal/repository"
	"github.com/axelhj/docs/internal/server"
)

func TestIndexRoute(t *testing.T) {
	tests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "index route",
			route:         "/",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "OK",
		},
		{
			description:   "non existing route",
			route:         "/i-dont-exist",
			expectedError: false,
			expectedCode:  404,
			expectedBody:  "Cannot GET /i-dont-exist",
		},
	}
	cfg := config.Config{}

	// TODO: Enable mocking of the couchster.
	// db := repository.NewCouchDB(cfg.CouchHost, cfg.CouchUser, cfg.CouchPassword)

	// couch := db.EnsureDB(cfg.DbName)
	// docService := service.NewDocumentService(couch)

	contentHandler := server.NewContentHandler()
	schemaHandler := server.NewSchemaHandler() // etc.

	if docSrv, err := server.NewDocServer(cfg.CouchHost, cfg.CouchUser, cfg.CouchPassword); err != nil {
		log.Fatal(err)
	}

	if err := docSrv.Start(cfg.AppPort); err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		AppName: "JustDocs",
	})

	server.Setup(app, contentHandler, schemaHandler)

	app.Listen(":" + cfg.AppPort)
	for _, test := range tests {
		req, err := http.NewRequest("GET", test.route, nil)
		assert.Nilf(t, err, test.description)
		res, err := app.Test(req, fiber.TestConfig{Timeout: 0, FailOnTimeout: false})
		assert.Equalf(t, test.expectedError, err != nil, test.description)
		if test.expectedError {
			continue
		}
		defer res.Body.Close()
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
		body, err := io.ReadAll(res.Body)
		assert.Nilf(t, err, test.description)
		assert.Equalf(t, test.expectedBody, string(body), test.description)
	}
}

func TestContentHandler_Create(t *testing.T) {
	// Mock service
	// mockService := &service.MockDocumentService{}
	// h := handler.NewContentHandler(mockService)

	// app := server.NewApp(h, nil) // reuse the same setup

	// Test request
	// resp, err := app.Test(fiber.Post("/api/content", `{"type":"note","title":"Test"}`))
	// assert.NoError(t, err)
	// assert.Equal(t, 201, resp.StatusCode)
}
