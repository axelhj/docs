package just_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
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

	app := Setup()

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
