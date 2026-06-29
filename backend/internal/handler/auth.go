package handler

import (
	"github.com/gofiber/fiber/v2"

	"taladsync/internal/middleware"
	"taladsync/internal/usecase"
)

type AuthHandler struct {
	uc *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

type mockLoginRequest struct {
	Role string `json:"role"`
}

type mockLoginResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

type meResponse struct {
	User any `json:"user"`
}

func (h *AuthHandler) MockLogin(c *fiber.Ctx) error {
	var req mockLoginRequest
	if err := c.BodyParser(&req); err != nil || req.Role == "" {
		return badRequest(c)
	}

	token, user, err := h.uc.MockLogin(c.Context(), req.Role)
	if err != nil {
		return RespondError(c, err)
	}

	return c.JSON(mockLoginResponse{Token: token, User: user})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims, err := middleware.GetClaims(c)
	if err != nil {
		return RespondError(c, errUnauthorized)
	}

	user, err := h.uc.GetMe(c.Context(), claims.UserID)
	if err != nil {
		return RespondError(c, err)
	}

	return c.JSON(meResponse{User: user})
}
