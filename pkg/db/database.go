package db

import (
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/kyosu-1/passkey-go-example/pkg/domain"
)

type UserDB interface {
	GetUser(id string) (*domain.User, error)
	AddUser(user *domain.User) error
}

type SessionDB interface {
	GetSession(sessionID string) (*webauthn.SessionData, error)
	DeleteSession(sessionID string)
	StartSession(data *webauthn.SessionData) string
}
