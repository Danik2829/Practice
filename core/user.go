package core

import "errors"

type UserService interface {
	Create(user User) error
	Get(id string) *User
	Update(user User) (User, error)
	Delete(id string) string
}

type UserStore interface {
	Create(user User) error
	Get(id string) *User
	Update(user User) (User, error)
	Delete(id string) string
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
