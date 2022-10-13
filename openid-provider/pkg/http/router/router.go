package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/http/handler"
)

func NotImplemented(w http.ResponseWriter, r *http.Request) {

}
func New(conf config.HTTP) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(conf.AuthEndpoint, handler.Auth).Methods(http.MethodGet)
	// two endpoints for login, first is for redering login form, second is for validation
	r.HandleFunc(conf.LoginEndpoint, handler.Login).Methods(http.MethodGet)
	r.HandleFunc(conf.LoginEndpoint, handler.CheckLogin).Methods(http.MethodPost)
	r.HandleFunc(conf.TokenEndpoint, handler.Token).Methods(http.MethodPost)
	r.HandleFunc(conf.CheckTokenEndpoint, handler.CheckToken).Methods(http.MethodGet)
	return r
}
