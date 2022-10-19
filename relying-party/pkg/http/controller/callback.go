package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/http/response"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
	"github.com/mv-kan/go-oidc-showcase/relying-party/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/relying-party/pkg/http/utils"
)

type Callback struct {
	opconf     config.OP
	clientconf config.Client
	httpconf   config.HTTP
	logger     log.Logger
}

func NewCallback(logger log.Logger, opconf config.OP, clientconf config.Client, httpconf config.HTTP) Callback {
	return Callback{logger: logger, opconf: opconf, clientconf: clientconf, httpconf: httpconf}

}

func (c Callback) GetToken(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	codes, ok := params["code"]
	if !ok {
		c.logger.Debug("parsing of code params failed")
		http.Error(w, "parsing of code params failed", http.StatusBadRequest)
		return
	}
	code := codes[0]
	states, ok := params["state"]
	if !ok {
		c.logger.Debug("parsing of state params failed")
		http.Error(w, "parsing of state params failed", http.StatusBadRequest)
		return
	}
	state := states[0]
	cookieState, err := utils.GetCookie(r, "state")
	if err != nil {
		c.logger.Debug(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if cookieState != state {
		c.logger.Debug("states are not the same")
		http.Error(w, "states are not the same", http.StatusBadRequest)
		return
	}
	opEndpoint, err := url.JoinPath(c.opconf.URL, c.opconf.TokenEndpoint)
	if err != nil {
		c.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// exchange code to token
	tokenRequest, err := utils.GetExchangeCodeToTokenRequest(opEndpoint, c.clientconf.ID, c.clientconf.Secret, c.clientconf.RedirectURI, code, "authorization_code")
	if err != nil {
		c.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	res, err := http.DefaultClient.Do(tokenRequest)
	if err != nil {
		c.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// read all json body
	tokenResponse := response.GrantToken{}
	jsonBytes, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	json.Unmarshal(jsonBytes, &tokenResponse)
	utils.SetCookie(w, "token", tokenResponse.AccessToken)
	http.Redirect(w, r, c.httpconf.RPHost, http.StatusFound)
}
