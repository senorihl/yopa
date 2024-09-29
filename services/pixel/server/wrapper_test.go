package server

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestStatus(t *testing.T) {
	app := Setup(func(query string, remoteAddr string) {})
	req, _ := http.NewRequest(
		"GET",
		"/status",
		nil,
	)
	res, err := app.Test(req, -1)
	assert.Equalf(t, 200, res.StatusCode, "/status should always respond as a 200 as soon as the server is started")
	body, err := io.ReadAll(res.Body)
	assert.Nilf(t, err, "/status should not throw an error")
	assert.Equalf(t, "OK", string(body), "/status should respond with `OK`")
}

func TestPixel(t *testing.T) {
	tests := []struct {
		description string

		// Test input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "pixel route",
			route:         "/pixel.gif",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  transPixel,
		},
	}

	// Iterate through test single test cases
	for _, test := range tests {
		app := Setup(func(query string, remoteAddr string) {

		})

		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := io.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.expectedBody, string(body), test.description)

		_ = app.Shutdown()
	}
}
