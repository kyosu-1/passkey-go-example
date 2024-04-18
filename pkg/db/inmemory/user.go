package inmemory

import (
	"fmt"
	"sync"

	"github.com/kyosu-1/passkey-go-example/pkg/domain"
)

type userdb struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

var usersDB *userdb = &userdb{
	users: make(map[string]*domain.User),
}

func NewUserDB() *userdb {
	return usersDB
}

func (db *userdb) GetUser(id string) (*domain.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	user, ok := db.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (db *userdb) AddUser(user *domain.User) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	_, ok := db.users[string(user.WebAuthnID())]
	if ok {
		return fmt.Errorf("user already exists")
	}
	_, ok = db.users[string(user.UserName)]
	if ok {
		return fmt.Errorf("user already exists")
	}
	db.users[string(user.WebAuthnID())] = user
	db.users[user.UserName] = user
	return nil
}
