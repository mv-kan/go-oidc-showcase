package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/config"
)

func NotImplemented(w http.ResponseWriter, r *http.Request) {

}
func New(conf config.HTTP) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(conf.AuthEndpoint, NotImplemented).Methods(http.MethodGet)
	// two endpoints for login, first is for redering login form, second is for validation
	r.HandleFunc(conf.LoginEndpoint, NotImplemented).Methods(http.MethodGet)
	r.HandleFunc(conf.LoginEndpoint, NotImplemented).Methods(http.MethodPost)
	r.HandleFunc(conf.TokenEndpoint, NotImplemented).Methods(http.MethodPost)
	r.HandleFunc(conf.CheckTokenEndpoint, NotImplemented).Methods(http.MethodGet)
	return r
}
