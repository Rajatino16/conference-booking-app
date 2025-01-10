package conference

import (
	"conference-booking/pkg/errors"
	"sync"
)

type Repository interface {
	Create(conference *Conference) error
	FindByName(name string) (*Conference, error)
}

type inMemoryRepository struct {
	conferences map[string]*Conference
	mutex       sync.Mutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		conferences: make(map[string]*Conference),
	}
}

func (r *inMemoryRepository) Create(conference *Conference) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.conferences[conference.Name]; exists {
		return errors.ErrConflict
	}

	r.conferences[conference.Name] = conference
	return nil
}

func (r *inMemoryRepository) FindByName(name string) (*Conference, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	conference, exists := r.conferences[name]
	if !exists {
		return nil, errors.ErrNotFound
	}

	return conference, nil
}
