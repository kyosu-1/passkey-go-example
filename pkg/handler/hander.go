package handler

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gorilla/sessions"

	"github.com/kyosu-1/passkey-go-example/pkg/db"
)

type hander struct {
	sessStore *sessions.CookieStore
	webAuthn  *webauthn.WebAuthn
	sessionDB db.SessionDB
	userDB   db.UserDB
}

func NewHandler(sessStore *sessions.CookieStore, webAuthn *webauthn.WebAuthn, sessionDB db.SessionDB, userDB db.UserDB) *hander {
	return &hander{
		sessStore: sessStore,
		webAuthn:  webAuthn,
		sessionDB: sessionDB,
		userDB: userDB,
	}
}
