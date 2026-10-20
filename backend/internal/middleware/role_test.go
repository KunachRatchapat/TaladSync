package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"taladsync/internal/auth"
	"taladsync/internal/domain"
	"taladsync/internal/middleware"
)

func TestRequireRoles(t *testing.T) {
	tokens := auth.NewTokenService([]byte(testSecret))

	app := fiber.New()
	app.Get("/walkin-only", middleware.Auth(tokens), middleware.RequireRoles(domain.RoleWalkin),
		func(c *fiber.Ctx) error {
			return c.SendString("ok")
		},
	)

	walkinToken, err := tokens.Sign(1, domain.RoleWalkin)
	require.NoError(t, err)
	staffToken, err := tokens.Sign(2, domain.RoleStaff)
	require.NoError(t, err)

	test := []struct {
		description  string
		token        string
		expectstatus int
	}{
		{
			description:  "ร้านค้าที่ไม่มีที่ประจำเข้าได้",
			token:        walkinToken,
			expectstatus: fiber.StatusOK,
		},
		{
			description:  "เจ้าหน้าที่ไม่สามารถเข้าได้",
			token:        staffToken,
			expectstatus: fiber.StatusForbidden,
		},
		{
			description:  "ไม่มี token เลย ไม่รู้ว่าเป็นใคร",
			token:        "",
			expectstatus: fiber.StatusUnauthorized,
		},
	}

	for _, tt := range test {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodGet, "/walkin-only", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)

			require.Equal(t, tt.expectstatus, resp.StatusCode)
		})
	}
}

func TestRequireRoles_MultipleRoles(t *testing.T) {
	tokens := auth.NewTokenService([]byte(testSecret))

	app := fiber.New()
	app.Get("/walkin-or-staff",
		middleware.Auth(tokens),
		middleware.RequireRoles(domain.RoleWalkin, domain.RoleStaff),
		func(c *fiber.Ctx) error {
			return c.SendString("ok")
		},
	)

	test := []struct {
		description  string
		role         domain.Role
		expectstatus int
	}{
		{"WALKIN เข้าได้", domain.RoleWalkin, fiber.StatusOK},
		{"STAFF เข้าได้", domain.RoleStaff, fiber.StatusOK},
		{"REGULAR เข้าไม่ได้", domain.RoleRegular, fiber.StatusForbidden},
	}

	for _, tt := range test {
		t.Run(tt.description, func(t *testing.T) {
			token, err := tokens.Sign(1, tt.role)
			require.NoError(t, err)

			req := httptest.NewRequest(fiber.MethodGet, "/walkin-or-staff", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			require.NoError(t, err)

			require.Equal(t, tt.expectstatus, resp.StatusCode)
		})
	}
}
