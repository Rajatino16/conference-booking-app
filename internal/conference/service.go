package conference

import (
	"conference-booking/pkg/errors"
)

type Service interface {
	AddConference(req AddConferenceRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AddConference(req AddConferenceRequest) error {
	if req.EndTime.Before(req.StartTime) || req.EndTime.Sub(req.StartTime).Hours() > 12 {
		return errors.ErrInvalidInput
	}

	existing, _ := s.repo.FindByName(req.Name)
	if existing != nil {
		return errors.ErrConflict
	}

	conference := &Conference{
		Name:           req.Name,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		AvailableSlots: req.AvailableSlots,
	}

	return s.repo.Create(conference)
}
