package conference

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	conferences map[string]*Conference
}

func NewMockRepository() Repository {
	return &mockRepository{conferences: make(map[string]*Conference)}
}

func (m *mockRepository) Create(conference *Conference) error {
	m.conferences[conference.Name] = conference
	return nil
}

func (m *mockRepository) FindByName(name string) (*Conference, error) {
	if conf, ok := m.conferences[name]; ok {
		return conf, nil
	}
	return nil, nil
}

func TestAddConference(t *testing.T) {
	repo := NewMockRepository()
	service := NewService(repo)

	req := AddConferenceRequest{
		Name:           "Test Conference",
		StartTime:      time.Now().UTC(),
		EndTime:        time.Now().UTC().Add(2 * time.Hour),
		AvailableSlots: 10,
	}

	err := service.AddConference(req)
	assert.NoError(t, err)
	// assert.NotNil(t, repo.FindByName(req.Name))
}
