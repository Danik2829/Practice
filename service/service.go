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

func (s *UserService) Create(user core.User) error {
	log.Printf("CreateUser called for ID: %s, email: %s", user.ID, user.Email)
	if user.ID == "" || user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Age == 0 {
		log.Printf("CreateUser validation failed: missing required fields for ID=%s", user.ID)
		return core.InvalidData
	}
	u := s.store.Get(user.ID)
	if u == nil {
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

func (s *UserService) Get(id string) *core.User {
	log.Printf("GetUser called for ID: %s", id)
	user := s.store.Get(id)
	if user == nil {
		log.Printf("GetUser: user %s not found", id)
		return nil
	}
	log.Printf("GetUser: user %s retrieved successfully", id)
	return user
}

func (s *UserService) Update(user core.User) (core.User, error) {
	log.Printf("UpdateUser called for ID: %s", user.ID)

	if user.ID == "" || user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Age == 0 {
		log.Printf("UpdateUser validation failed: missing required fields for ID=%s", user.ID)
		return core.User{}, core.InvalidData
	}

	updatedUser, err := s.store.Update(user)
	if err != nil {
		log.Printf("UpdateUser: store.Update failed for ID %s: %v", user.ID, err)
		return core.User{}, err
	}

	log.Printf("UpdateUser: user %s updated successfully", user.ID)
	return updatedUser, nil
}

func (s *UserService) Delete(id string) string {
	log.Printf("DeleteUser called for ID: %s", id)
	deletedId := s.store.Delete(id)
	if deletedId == "" {
		log.Printf("DeleteUser: store.Delete failed for ID %s", id)
		return ""
	}
	log.Printf("DeleteUser: user %s deleted successfully", id)
	return deletedId
}
