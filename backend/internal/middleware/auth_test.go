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

const testSecret = "test-secret"

func newProtectedApp(tokens auth.TokenService) *fiber.App {
	app := fiber.New()
	app.Get("/protected", middleware.Auth(tokens), func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	return app
}

func TestAuth_Unauthorized(t *testing.T) {
	tokens := auth.NewTokenService([]byte(testSecret))
	app := newProtectedApp(tokens)

	validToken, err := tokens.Sign(1, domain.RoleWalkin)
	require.NoError(t, err)

	otherTokens := auth.NewTokenService([]byte("another-secret"))
	foreignToken, err := otherTokens.Sign(1, domain.RoleWalkin)
	require.NoError(t, err)

	test := []struct {
		description  string
		authHeader   string
		expectstatus int
	}{
		{
			description:  "ไม่มี Authorization header",
			authHeader:   "",
			expectstatus: fiber.StatusUnauthorized,
		},
		{
			description:  "ไม่มีคำว่า Bearer",
			authHeader:   "abc123",
			expectstatus: fiber.StatusUnauthorized,
		},
		{
			description:  "มี Bearer แต่ token ว่าง",
			authHeader:   "Bearer ",
			expectstatus: fiber.StatusUnauthorized,
		},
		{
			description:  "token มั่ว",
			authHeader:   "Bearer not-a-jwt",
			expectstatus: fiber.StatusUnauthorized,
		},
		{
			description:  "token เซ็นด้วย secret อื่น",
			authHeader:   "Bearer " + foreignToken,
			expectstatus: fiber.StatusUnauthorized,
		},
		{
			description:  "token ถูกต้อง",
			authHeader:   "Bearer " + validToken,
			expectstatus: fiber.StatusOK,
		},
	}

	for _, tt := range test {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)

			require.Equal(t, tt.expectstatus, resp.StatusCode)
		})
	}
}

func TestAuth_PassesClaimsToHandler(t *testing.T) {
	tokens := auth.NewTokenService([]byte(testSecret))

	app := fiber.New()
	app.Get("/me", middleware.Auth(tokens), func(c *fiber.Ctx) error {
		claims, err := middleware.GetClaims(c)
		require.NoError(t, err)
		require.Equal(t, uint(42), claims.UserID)
		require.Equal(t, domain.RoleStaff, claims.Role)
		return c.SendString("ok")
	})

	token, err := tokens.Sign(42, domain.RoleStaff)
	require.NoError(t, err)

	req := httptest.NewRequest(fiber.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
}
