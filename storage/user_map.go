package storage

import "restAPI/core"

type MapStorage struct {
	users map[string]core.User
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		users: make(map[string]core.User),
	}
}

func (m *MapStorage) Create(u core.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *MapStorage) Get(id string) (core.User, error) {
	if user, exist := m.users[id]; !exist {
		return core.User{}, core.NotFound
	} else {
		return user, nil
	}
}

func (m *MapStorage) Update(user core.User) error {
	if _, exist := m.users[user.ID]; !exist {
		return core.NotFound
	} else {
		m.users[user.ID] = user
		return nil
	}
}

func (m *MapStorage) Delete(id string) error {
	if _, exist := m.users[id]; !exist {
		return core.NotFound
	} else {
		delete(m.users, id)
		return nil
	}
}
