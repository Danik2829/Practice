package core

import "errors"

type UserStorage interface {
	Create(user User) error
	Get(id string) (User, error)
	Update(user User) error
	Delete(id string) error
}

type User struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Age       uint   `json:"age"`
}

var (
	InvalidData = errors.New("invalid data")
	UserExist   = errors.New("user already exists")
	NotFound    = errors.New("User with this ID not exist")
)
