package service

import "restAPI/core"

type UserService struct {
	stor core.UserStorage
}

func NewUserService(stor core.UserStorage) *UserService {
	return &UserService{
		stor: stor,
	}
}

func (s *UserService) CreateUser(user core.User) error {
	if user.ID == "" || user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Age == 0 {
		return core.InvalidData
	}
	if _, exist := s.stor.Get(user.ID); exist == nil {
		return core.UserExist
	}
	if err := s.stor.Create(user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUser(id string) (core.User, error) {
	if user, exist := s.stor.Get(id); exist != nil {
		return core.User{}, core.NotFound
	} else {
		return user, nil
	}
}

func (s *UserService) UpdateUser(user core.User) error {
	if user.ID == "" || user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Age == 0 {
		return core.InvalidData
	}
	return s.stor.Update(user)
}

func (s *UserService) DeleteUser(id string) error {
	return s.stor.Delete(id)
}
