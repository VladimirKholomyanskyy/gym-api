package account

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}
func (s *UserService) CreateUser(externalUserID string) (*User, error) {
	user := &User{
		ExternalID: externalUserID,
	}
	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindUserByID(id uint) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) FindUserByExternalID(id string) (*User, error) {
	return s.repo.FindByExternalID(id)
}

func (s *UserService) UpdateUser(user *User) error {
	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}
