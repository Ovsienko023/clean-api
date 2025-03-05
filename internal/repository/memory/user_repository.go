package memory

import (
	"api/internal/domain"
	"api/internal/domain/errdomain"
	"context"
	"github.com/google/uuid"
	"sync"
)

type UserRepositoryMemory struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

var _ domain.UserRepository = &UserRepositoryMemory{}

func NewUserRepositoryMemory() *UserRepositoryMemory {
	return &UserRepositoryMemory{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepositoryMemory) Get(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errdomain.ErrObjectNotFound
	}

	return user, nil
}

func (r *UserRepositoryMemory) Search(ctx context.Context, query domain.QueryOptions) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var users []*domain.User

	for _, u := range r.users {
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepositoryMemory) Create(ctx context.Context, user domain.User) (*string, error) {
	id := uuid.New().String()
	user.ID = id

	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = &user

	return &id, nil
}

func (r *UserRepositoryMemory) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errdomain.ErrObjectNotFound
	}

	delete(r.users, id)

	return nil
}
