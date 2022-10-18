package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mv-kan/go-oidc-showcase/resource-server/pkg/http/controller"
)

func New(messageCtrl controller.ProtectedMessage) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", messageCtrl.GetProtectedSuperSecret).Methods(http.MethodGet)
	return r
}
