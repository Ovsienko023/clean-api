package usecase

import (
	"api/internal/domain"
	"context"
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: repo,
	}
}

func (u *UserUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return u.userRepo.Get(ctx, id)
}

func (u *UserUseCase) ListUser(ctx context.Context, query domain.QueryOptions) ([]*domain.User, error) {
	return u.userRepo.Search(ctx, query)
}

func (u *UserUseCase) CreateUser(ctx context.Context, user domain.User) (*string, error) {
	return u.userRepo.Create(ctx, user)
}

func (u *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	return u.userRepo.Delete(ctx, id)
}
