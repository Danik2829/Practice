package service

import (
	"log"
	"restAPI/core"
)

type UserService struct {
	store core.UserStore
}

func NewUserService(store core.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) CreateUser(user core.User) error {
	log.Printf("CreateUser called for ID: %s, email: %s", user.ID, user.Email)
	if user.ID == "" || user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Age == 0 {
		log.Printf("CreateUser validation failed: missing required fields for ID=%s", user.ID)
		return core.InvalidData
	}
	_, exist := s.store.Get(user.ID)
	if exist == nil {
		log.Printf("CreateUser: user with ID %s already exists", user.ID)
		return core.UserExist
	}
	err := s.store.Create(user)
	if err != nil {
		log.Printf("CreateUser: store.Create failed for ID %s: %v", user.ID, err)
		return err
	}
	log.Printf("CreateUser: user %s created successfully", user.ID)
	return nil
}

func (s *UserService) GetUser(id string) (core.User, error) {
	log.Printf("GetUser called for ID: %s", id)
	user, exist := s.store.Get(id)
	if exist != nil {
		log.Printf("GetUser: user %s not found", id)
		return core.User{}, core.NotFound
	}
	log.Printf("GetUser: user %s retrieved successfully", id)
	return user, nil
}

func (s *UserService) UpdateUser(user core.User) error {
	log.Printf("UpdateUser called for ID: %s", user.ID)
	if user.ID == "" || user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Age == 0 {
		log.Printf("UpdateUser validation failed: missing required fields for ID=%s", user.ID)
		return core.InvalidData
	}
	err := s.store.Update(user)
	if err != nil {
		log.Printf("UpdateUser: store.Update failed for ID %s: %v", user.ID, err)
		return err
	}
	log.Printf("UpdateUser: user %s updated successfully", user.ID)
	return nil
}

func (s *UserService) DeleteUser(id string) error {
	log.Printf("DeleteUser called for ID: %s", id)
	err := s.store.Delete(id)
	if err != nil {
		log.Printf("DeleteUser: store.Delete failed for ID %s: %v", id, err)
		return err
	}
	log.Printf("DeleteUser: user %s deleted successfully", id)
	return nil
}
