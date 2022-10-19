package main

import (
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
	"github.com/mv-kan/go-oidc-showcase/relying-party/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/relying-party/pkg/http/controller"
	"github.com/mv-kan/go-oidc-showcase/relying-party/pkg/http/router"
)

// env keys names
const (
	// client
	ClientID          = "CLIENT_ID"
	ClientSecret      = "CLIENT_SECRET"
	ClientRedirectURI = "CLIENT_REDIRECT_URI"

	// relying party
	RPURL            = "RP_URL"
	RPHost           = "RP_HOST"
	IndexEndpoint    = "RP_INDEX_ENDPOINT"
	CallbackEndpoint = "RP_CALLBACK_ENDPOINT"

	// openid provider
	AuthEndpoint       = "OP_AUTH_ENDPOINT"
	LoginEndpoint      = "OP_LOGIN_ENDPOINT"
	TokenEndpoint      = "OP_TOKEN_ENDPOINT"
	CheckTokenEndpoint = "OP_CHECK_TOKEN_ENDPOINT"
	OPHost             = "OP_HOST"
	OPURL              = "OP_URL"

	RSURL               = "RS_URL"
	RSProtectedEndpoint = "RS_PROTECTED_ENDPOINT"

	LogFilePath = "LOG_FILEPATH"
)

func run() {
	// init all configs
	// load env file
	err := godotenv.Load("rp.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	clientconf := config.Client{
		ID:          os.Getenv(ClientID),
		Secret:      os.Getenv(ClientSecret),
		RedirectURI: os.Getenv(ClientRedirectURI),
	}
	httpconf := config.HTTP{
		RPHost:           os.Getenv(RPHost),
		RPURL:            os.Getenv(RPURL),
		IndexEndpoint:    os.Getenv(IndexEndpoint),
		CallbackEndpoint: os.Getenv(CallbackEndpoint),
	}
	logconf := config.Logger{
		FilePath: os.Getenv(LogFilePath),
	}
	opconf := config.OP{
		URL:           os.Getenv(OPURL),
		TokenEndpoint: os.Getenv(TokenEndpoint),
		AuthEndpoint:  os.Getenv(AuthEndpoint),
	}
	rsconf := config.RS{
		URL:               os.Getenv(RSURL),
		ProtectedEndpoint: os.Getenv(RSProtectedEndpoint),
	}
	// create logger
	f, err := os.OpenFile(logconf.FilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	logger := log.New(map[log.LogOutput]io.Writer{
		log.FileOutput:    f,
		log.ConsoleOutput: os.Stdout,
	})
	// get all controllers
	indexCtrl := controller.NewIndex(logger, opconf, clientconf, rsconf)
	callbackCtrl := controller.NewCallback(logger, opconf, clientconf, httpconf)
	// get router
	r := router.New(httpconf, callbackCtrl, indexCtrl)
	// run
	logger.Info("run server on " + httpconf.RPHost)
	http.ListenAndServe(httpconf.RPHost, r)
}
func main() {
	run()
}
