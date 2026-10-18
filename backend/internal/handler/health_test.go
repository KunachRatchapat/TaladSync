package handler

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type fakePinger struct {
	err error
}

func (f fakePinger) PingContext(ctx context.Context) error {
	return f.err
}

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		description  string
		pingErr      error
		expectstatus int
	}{
		{
			description:  "db up",
			pingErr:      nil,
			expectstatus: fiber.StatusOK,
		},
		{
			description:  "db down",
			pingErr:      errors.New("connection refused"),
			expectstatus: fiber.StatusServiceUnavailable,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			app := fiber.New()
			app.Get("/health", NewHealthHandler(fakePinger{err: test.pingErr}).Check)

			resp, err := app.Test(httptest.NewRequest("GET", "/health", nil))
			require.NoError(t, err)
			require.Equal(t, test.expectstatus, resp.StatusCode)
		})
	}
}
