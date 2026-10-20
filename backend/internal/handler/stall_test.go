package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"taladsync/internal/usecase"
)

type fakeStallRepo struct {
	zones []string
	err   error
}

func (f fakeStallRepo) ListZones(ctx context.Context) ([]string, error) {
	return f.zones, f.err
}

func TestListZones(t *testing.T) {
	tests := []struct {
		description  string
		repo         fakeStallRepo
		expectstatus int
		expectzones  []string
	}{
		{
			description:  "returns zones",
			repo:         fakeStallRepo{zones: []string{"โซน A", "โซน B"}},
			expectstatus: fiber.StatusOK,
			expectzones:  []string{"โซน A", "โซน B"},
		},
		{
			description:  "db error",
			repo:         fakeStallRepo{err: errors.New("db down")},
			expectstatus: fiber.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			app := fiber.New()
			app.Get("/zones", NewStallHandler(usecase.NewStallUsecase(test.repo)).ListZones)

			resp, err := app.Test(httptest.NewRequest("GET", "/zones", nil))
			require.NoError(t, err)
			require.Equal(t, test.expectstatus, resp.StatusCode)

			if test.expectzones != nil {
				var body struct {
					Zones []string `json:"zones"`
				}
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
				require.Equal(t, test.expectzones, body.Zones)
			}
		})
	}
}
