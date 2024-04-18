package handler

import (
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
)

func (h *hander) AssertionOptions(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// Some options are complemented in the frontend
	// Ref: https://github.com/MasterKale/SimpleWebAuthn/blob/5229cebbcc2d087b7eaaaeb9886f53c9e1d93522/packages/browser/src/methods/startAuthentication.ts#L72-L76
	options, sessionData, err := h.webAuthn.BeginDiscoverableLogin()
	if err != nil {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "authentication",
		Value: h.sessionDB.StartSession(sessionData),
		Path:  "/",
	})
	jsonResponse(w, options, http.StatusOK)
}

func (h *hander)AssertionResult(w http.ResponseWriter, r *http.Request) {
	session, _ := h.sessStore.Get(r, "session-name")

	if err := r.ParseForm(); err != nil {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("authentication")
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

	credential, err := h.webAuthn.FinishDiscoverableLogin(func(rawId []byte, userhandle []byte) (user webauthn.User, err error) {
		return h.userDB.GetUser(string(userhandle))
	}, *sessionData, r)
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

	if credential.Authenticator.CloneWarning {
		jsonResponse(w, FIDO2Response{
			Status:       "failed",
			ErrorMessage: "authenticator is cloned",
		}, http.StatusBadRequest)
		return
	}

	// セッションの設定
	session.Values["authenticated"] = true
	session.Save(r, w)

	h.sessionDB.DeleteSession(cookie.Value)

	jsonResponse(w, FIDO2Response{
		Status:       "ok",
		ErrorMessage: "",
	}, http.StatusOK)
}