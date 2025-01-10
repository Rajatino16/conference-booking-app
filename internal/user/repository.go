package user

import (
	"sync"

	"conference-booking/pkg/errors"
)

type Repository interface {
	Create(user *User) error
	FindByID(id string) (*User, error)
}

type inMemoryRepository struct {
	users map[string]*User
	mutex sync.Mutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		users: make(map[string]*User),
	}
}

func (r *inMemoryRepository) Create(user *User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.ErrConflict
	}

	r.users[user.ID] = user
	return nil
}

func (r *inMemoryRepository) FindByID(id string) (*User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.ErrNotFound
	}

	return user, nil
}
