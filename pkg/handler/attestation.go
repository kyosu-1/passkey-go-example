package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"

	"github.com/kyosu-1/passkey-go-example/pkg/domain"
)

type AttestationOptionsParam struct {
	Username string
}

func (h *hander)AttestationOptions(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	var p AttestationOptionsParam

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	username := p.Username
	if username == "" {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: "Missing username",
		}, http.StatusBadRequest)
		return
	}

	user, err := h.userDB.GetUser(username)
	if err != nil {
		displayName := strings.Split(username, "@")[0]
		user = domain.NewUser(username, displayName)
		h.userDB.AddUser(user)
	}
	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
		credCreationOpts.AuthenticatorSelection.RequireResidentKey = protocol.ResidentKeyRequired()
		credCreationOpts.AuthenticatorSelection.ResidentKey = protocol.ResidentKeyRequirementRequired
	}
	options, sessionData, err := h.webAuthn.BeginRegistration(user, registerOptions)
	if err != nil {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "registration",
		Value: h.sessionDB.StartSession(sessionData),
		Path:  "/",
	})
	jsonResponse(w, options, http.StatusOK)
}

func (h *hander)AttestationResult(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	cookie, err := r.Cookie("registration")
	if err != nil {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	sessionData, err := h.sessionDB.GetSession(cookie.Value)
	if err != nil {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	user, err := h.userDB.GetUser(string(sessionData.UserID))
	if err != nil {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	credential, err := h.webAuthn.FinishRegistration(user, *sessionData, r)
	if err != nil {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: "user was not verified",
		}, http.StatusBadRequest)
		return
	}

	user.AddCredential(*credential)

	h.sessionDB.DeleteSession(cookie.Value)

	jsonResponse(w, FIDO2Response{
		Status:       "ok",
		ErrorMessage: "",
	}, http.StatusOK)
}
