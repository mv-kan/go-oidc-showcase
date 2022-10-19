package controller

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
	"github.com/mv-kan/go-oidc-showcase/relying-party/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/relying-party/pkg/http/utils"
)

type Index struct {
	opconf     config.OP
	clientconf config.Client
	rsconf     config.RS
	logger     log.Logger
}

func NewIndex(logger log.Logger, opconf config.OP, clientconf config.Client, rsconf config.RS) Index {
	return Index{logger: logger, opconf: opconf, clientconf: clientconf, rsconf: rsconf}

}
func (i Index) GetMessage(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetCookie(r, "token")
	// if token does not exist, then redirect to auth server's login page
	if errors.Is(err, http.ErrNoCookie) {
		// save state in cookies
		state := uuid.New().String()
		utils.SetCookie(w, "state", state)
		// get openid provider authenticate url
		opEndpoint, err := url.JoinPath(i.opconf.URL, i.opconf.AuthEndpoint)
		if err != nil {
			i.logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		redirectURL, err := utils.GetAuthCodeURL(opEndpoint, i.clientconf.ID, i.clientconf.RedirectURI, state, "openid", "code")
		if err != nil {
			i.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else if err != nil {
		i.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// if token exists send token to protected resource and send message from this resource
	url, err := url.JoinPath(i.rsconf.URL, i.rsconf.ProtectedEndpoint)
	if err != nil {
		i.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		i.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		i.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		i.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseBody := []byte("Message from protected server: ")
	responseBody = append(responseBody, body...)
	w.Write(responseBody)
	//i.writeProtectedInfo(token, w, r)
}
func (i Index) writeProtectedInfo(token string, w http.ResponseWriter, r *http.Request) {

	// if token exists send token to protected resource and send message from this resource
	url, err := url.JoinPath(i.rsconf.URL, i.rsconf.ProtectedEndpoint)
	if err != nil {
		i.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		i.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		i.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		i.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseBody := []byte("Message from protected server: ")
	responseBody = append(responseBody, body...)
	w.Write(responseBody)
}
