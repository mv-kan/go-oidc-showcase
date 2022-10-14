package controller

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"

	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/http/utils"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/storage"
)

type Token struct {
	logger          log.Logger
	authCodeStorage storage.AuthCode
	clientStorage   storage.Client
	tokenStorage    storage.AccessToken
}

func (t Token) SwitchCodeToToken(w http.ResponseWriter, r *http.Request) {
	// authorization code check
	err := r.ParseForm()
	if err != nil {
		t.logger.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")
	clientID := r.FormValue("client_id")
	grantType := r.FormValue("grant_type")
	// validate grantType
	if grantType != "authorization_code" {
		t.logger.Debug("only grant type is authorization_code, grant_type = " + grantType)
		http.Error(w, "the only grant type is authorization_code, grant_type="+grantType, http.StatusBadRequest)
		return
	}
	// get auth code obj from storage
	authCode, err := t.authCodeStorage.Get(code)
	if err != nil {
		t.logger.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot get auth code:%s", err), http.StatusInternalServerError)
		return
	}
	// compare redirect uri
	if authCode.RedirectURI != redirectURI {
		t.logger.Debug(fmt.Sprintf("redirect uris are not the same: %s(in auth req) %s (in token req)", authCode.RedirectURI, redirectURI))
		http.Error(w, "redirect uris are not the same", http.StatusBadRequest)
		return
	}
	// validate clientID
	if authCode.ClientID != clientID {
		t.logger.Debug(fmt.Sprintf("clientIDs are not the same clientid auth req = %s, clientid token req = %s", authCode.ClientID, clientID))
		http.Error(w, "client ids are not the same", http.StatusBadRequest)
		return
	}
	// Check Authorization basic Client
	client, err := t.clientStorage.Get(clientID)
	if errors.Is(err, storage.ErrNotFound) {
		t.logger.Debug(fmt.Sprintf("clientID does not exist clientID=%s", clientID))
		http.Error(w, "unauthorized client does not exist", http.StatusUnauthorized)
		return
	} else if err != nil {
		t.logger.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot get auth request:%s", err), http.StatusInternalServerError)
		return
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		t.logger.Debug("no basic auth is presented")
		http.Error(w, "unauthorized no basic auth is presented", http.StatusUnauthorized)
		return
	}
	clientIDHash := sha256.Sum256([]byte(clientID))
	secretHash := sha256.Sum256([]byte(clientSecret))
	expectedClientIDHash := sha256.Sum256([]byte(client.ID))
	expectedSecretHash := sha256.Sum256([]byte(client.Secret))

	// Use the subtle.ConstantTimeCompare() function to check if
	// the provided username and password hashes equal the
	// expected username and password hashes. ConstantTimeCompare
	// will return 1 if the values are equal, or 0 otherwise.
	// Importantly, we should to do the work to evaluate both the
	// username and password before checking the return values to
	// avoid leaking information.
	clientIDMatch := (subtle.ConstantTimeCompare(clientIDHash[:], expectedClientIDHash[:]) == 1)
	secretMatch := (subtle.ConstantTimeCompare(secretHash[:], expectedSecretHash[:]) == 1)
	if !(clientIDMatch && secretMatch) {
		t.logger.Debug(fmt.Sprintf("clientid or password is invalid, clientID=%s", clientID))
		http.Error(w, "clientid or password is invalid, unauthorized", http.StatusUnauthorized)
		return
	}

	// after successful validation we generate token and add it to token storage
	token, err := utils.SwitchCodeToToken(t.authCodeStorage, t.tokenStorage, authCode, redirectURI, grantType)
	if errors.Is(err, storage.ErrNotFound) {
		t.logger.Debug(fmt.Sprintf("clientID does not exist clientID=%s", clientID))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	} else if err != nil {
		t.logger.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot generate access token:%s", err), http.StatusInternalServerError)
		return
	}
	// then write successful response
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Cache-Control", "no-store")

	msg := map[string]any{
		"access_token": token.ID,
		"token_type":   "Bearer",
		"expires_in":   99999999,
		"id_token":     token.GetID(),
	}
	err = utils.ResponseJSON(w, http.StatusOK, msg)
	if err != nil {
		t.logger.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot generate access token:%s", err), http.StatusInternalServerError)
	}

}

func (t Token) CheckToken(w http.ResponseWriter, r *http.Request) {
	// authorization code check
	err := r.ParseForm()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	accessToken := r.FormValue("token")
	// check token by getting it from "db"
	token, err := t.tokenStorage.Get(accessToken)
	if errors.Is(storage.ErrNotFound, err) {
		log.Debug("token is not present in token storage")
		http.Error(w, "no such token", http.StatusUnauthorized)
		return
	} else if err != nil {
		t.logger.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot get token:%s", err), http.StatusInternalServerError)
		return
	}
	// yes I do understand that we don't use scopes at all
	// but this project is for example and learning reasons
	response := map[string]string{
		"access_token": token.ID,
		"username":     token.UserID,
	}
	utils.ResponseJSON(w, http.StatusOK, response)
}
