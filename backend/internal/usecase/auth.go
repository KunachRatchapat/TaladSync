package usecase

import (
	"context"

	"taladsync/internal/auth"
	"taladsync/internal/domain"
)

type UserReader interface {
	FindByRole(ctx context.Context, role domain.Role) (*domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
}

type AuthUsecase struct {
	users  UserReader
	tokens auth.TokenService
}

func NewAuthUsecase(users UserReader, tokens auth.TokenService) *AuthUsecase {
	return &AuthUsecase{users: users, tokens: tokens}
}

func (uc *AuthUsecase) MockLogin(ctx context.Context, roleStr string) (string, *domain.User, error) {
	role := domain.Role(roleStr)
	if role != domain.RoleRegular && role != domain.RoleWalkin && role != domain.RoleStaff {
		return "", nil, domain.ErrInvalidRole
	}

	user, err := uc.users.FindByRole(ctx, role)
	if err != nil {
		return "", nil, err
	}

	token, err := uc.tokens.Sign(user.ID, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (uc *AuthUsecase) GetMe(ctx context.Context, userID uint) (*domain.User, error) {
	return uc.users.FindByID(ctx, userID)
}
