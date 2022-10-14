package controller

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/http/utils"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/storage"
)

type Auth struct {
	logger         log.Logger
	httpconf       config.HTTP
	requestStorage storage.AuthRequest
	clientStorage  storage.Client
}

func (a Auth) Auth(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()

	// parse parameters
	randID := uuid.New().String()
	authRequest := oidc.AuthRequest{
		ID: randID,
	}

	// clientID check
	clientID, ok := params["client_id"]
	if !ok {
		a.logger.Debug("auth req params missing client_id")
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.InvalidRequest)
		return
	}
	authRequest.ClientID = clientID[0]

	// redirectURI check
	redirectURI, ok := params["redirect_uri"]
	if !ok {
		a.logger.Debug("auth req params missing redirect uri")
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.InvalidRequest)
		return
	}
	authRequest.RedirectURI = redirectURI[0]

	// state check
	state, ok := params["state"]
	if !ok {
		a.logger.Debug("auth req params missing state")
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.InvalidRequest)
		return
	}
	authRequest.State = state[0]

	// scope check
	scope, ok := params["scope"]
	if !ok {
		a.logger.Debug("auth req params missing scope")
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.InvalidRequest)
		return
	}
	authRequest.Scope = scope

	// response type check
	responseType, ok := params["response_type"]
	if !ok {
		a.logger.Debug("auth req params missing response_type")
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.InvalidRequest)
		return
	}
	authRequest.ResponseType = responseType

	// verify that openid is present in scope
	if !utils.Contains(authRequest.Scope, "openid") {
		a.logger.Debug("auth req params: in scope missing openid")
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.InvalidScope)
		return
	}

	// the only available flow is Authorization code flow
	if len(authRequest.ResponseType) != 1 || authRequest.ResponseType[0] != "code" {
		a.logger.Debug("only supported response type is code")
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.UnsupportedResponseType)
		return
	}

	err := a.requestStorage.Add(authRequest)
	if err != nil {
		a.logger.Error(err.Error())
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.ServerError)
		return
	}
	// check if redirect URI is valid
	ok, err = utils.ValidateClientRedirectURI(a.clientStorage, authRequest.ClientID, authRequest.RedirectURI)
	if err != nil {
		a.logger.Error(err.Error())
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.ServerError)
		return
	}
	if !ok {
		a.logger.Debug("not registered regirect_uri")
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.ServerError)
		return
	}
	// check if client is registered (in storage)
	clients, err := a.clientStorage.GetAll()
	if err != nil {
		a.logger.Error(err.Error())
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.ServerError)
		return
	}
	ids := func() []string {
		tmp := make([]string, 0)

		for _, client := range clients {
			tmp = append(tmp, client.ID)
		}
		return tmp
	}()
	if !utils.Contains(ids, authRequest.ClientID) {
		a.logger.Debug("not registered client")
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.UnauthorizedClient)
		return
	}

	// authenticate user redirecting him to login page
	loginRedirect, err := url.JoinPath(a.httpconf.OPURL, a.httpconf.LoginEndpoint)
	loginRedirectParams := fmt.Sprintf("?authRequestID=%s", authRequest.GetID())
	if err != nil {
		a.logger.Error(err.Error())
		utils.AuthErrorResponse(a.clientStorage, w, r, authRequest, utils.ServerError)
		return
	}
	a.logger.Info("Successfully processed this request and redirect to " + loginRedirect + loginRedirectParams)
	http.Redirect(w, r, loginRedirect+loginRedirectParams, http.StatusFound)
}
