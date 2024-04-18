package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gorilla/sessions"

	"github.com/kyosu-1/passkey-go-example/pkg/db/inmemory"
	"github.com/kyosu-1/passkey-go-example/pkg/handler"
)

var env struct {
	SessionSecret []byte
}

func loadEnv() {
	var err error
	env.SessionSecret, err = base64.StdEncoding.DecodeString(os.Getenv("SESSION_SECRET"))
	if err != nil {
		env.SessionSecret = []byte("your-secret-key")
	}
}

func main() {
	loadEnv()
	wconfig := &webauthn.Config{
		RPDisplayName: "Go WebAuthn",
		RPID:          "localhost",
		RPOrigin:      "http://localhost:8080",
	}
	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		log.Fatal(err)
	}

	sessStore := sessions.NewCookieStore(env.SessionSecret)

	r := chi.NewRouter()

	userDB := inmemory.NewUserDB()
	sessDB := inmemory.NewSessionDB()

	h := handler.NewHandler(sessStore, webAuthn, sessDB, userDB)

	r.Post("/attestation/options", h.AttestationOptions)
	r.Post("/attestation/result", h.AttestationResult)

	r.Post("/assertion/options", h.AssertionOptions)
	r.Post("/assertion/result", h.AssertionResult)

	r.Handle("/", http.FileServer(http.Dir("./templates")))
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/login.html")
	})

	r.Get("/success", authMiddleware(successPageHandler, sessStore))

	log.Println("server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func authMiddleware(next http.HandlerFunc, sessStore sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessStore.Get(r, "session-name")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

func successPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "authenticated/success.html")
}
