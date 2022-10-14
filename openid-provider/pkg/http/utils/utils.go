package utils

// this file package is utter garbage, I need to redo everything, but
// I am not going to do it :)
import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/storage"
)

// this are error codes in callback response
const (
	InvalidRequest          string = "invalid_request"
	UnauthorizedClient      string = "unauthorized_client"
	AccessDenied            string = "access_denied"
	UnsupportedResponseType string = "unsupported_response_type"
	InvalidScope            string = "invalid_scope"
	ServerError             string = "server_error"
	TemporaryUnavailable    string = "temporary_unavailable"
)

// returns json object with errorCode, also it has redirection uri with error params
// errorCode - is OAuth Error code (https://www.rfc-editor.org/rfc/rfc6749#section-4.1.2), you can find then in pkg
func AuthErrorResponse(clientStorage storage.Client, w http.ResponseWriter, r *http.Request, authReq oidc.AuthRequest, errorCode string) {
	callback, err := url.Parse(authReq.RedirectURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// auth server MUST NOT redirect if redirectURI is invalid
	ok, err := ValidateClientRedirectURI(clientStorage, authReq.ClientID, callback.String())
	if errors.Is(err, storage.ErrNotFound) {
		ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "client does not exist"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "redirectURI is not valid. Maybe it does not exist in client registered redirectURIs list"})
		return
	}
	// add to callback query parameters
	values := callback.Query()

	// add state parameter if exists
	if authReq.State != "" {
		values.Add("state", authReq.State)
	}

	// add error parameter according to https://www.rfc-editor.org/rfc/rfc6749#section-4.1.2
	values.Add("error", errorCode)

	callback.RawQuery = values.Encode()

	// oauth specification requires add form urlencoded
	w.Header().Add("Content-Type", "application/x-www-form-urlencoded")
	http.Redirect(w, r, callback.String(), http.StatusFound)
}

// returns true if redirectURI is valid for clientID, otherwise false
func ValidateClientRedirectURI(clientStorage storage.Client, clientID string, redirectURI string) (bool, error) {
	// get client
	client, err := clientStorage.Get(clientID)
	if err != nil {
		return false, err
	}
	// check with ParseRequestURI all urls for validity and then compare them
	_, err = url.ParseRequestURI(redirectURI)
	if err != nil {
		return false, err
	}
	// check if redirect uri is in client registered redirect uris list
	return Contains(client.RedirectURIs, redirectURI), nil
}

func Contains[T comparable](sl []T, elem T) bool {
	for _, value := range sl {
		if value == elem {
			return true
		}
	}
	return false
}
func ResponseJSON(w http.ResponseWriter, code int, payload any) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}
