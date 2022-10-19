package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mv-kan/go-oidc-showcase/relying-party/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/relying-party/pkg/http/controller"
)

func New(conf config.HTTP, callbackCtrl controller.Callback, indexCtrl controller.Index) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc(conf.IndexEndpoint, indexCtrl.GetMessage).Methods(http.MethodGet)

	r.HandleFunc(conf.CallbackEndpoint, callbackCtrl.GetToken).Methods(http.MethodGet)
	return r
}
