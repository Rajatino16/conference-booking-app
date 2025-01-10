package user

type Service interface {
	AddUser(req AddUserRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AddUser(req AddUserRequest) error {
	user := &User{ID: req.ID}
	return s.repo.Create(user)
}
