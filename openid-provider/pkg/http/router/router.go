package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/http/controller"
)

func New(conf config.HTTP, authCtrl controller.Auth, loginCtrl controller.Login, tokenCtrl controller.Token) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc(conf.AuthEndpoint, authCtrl.Auth).Methods(http.MethodGet)
	// two endpoints for login, first is for rendering login form, second is for validation
	r.HandleFunc(conf.LoginEndpoint, loginCtrl.Login).Methods(http.MethodGet)
	r.HandleFunc(conf.LoginEndpoint, loginCtrl.CheckLogin).Methods(http.MethodPost)
	r.HandleFunc(conf.TokenEndpoint, tokenCtrl.SwitchCodeToToken).Methods(http.MethodPost)
	r.HandleFunc(conf.CheckTokenEndpoint, tokenCtrl.CheckToken).Methods(http.MethodGet)
	return r
}
