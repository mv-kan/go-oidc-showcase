package main

import (
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/http/controller"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/http/router"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/storage"
)

// env keys names
const (
	AuthEndpoint       = "AUTH_ENDPOINT"
	LoginEndpoint      = "LOGIN_ENDPOINT"
	TokenEndpoint      = "TOKEN_ENDPOINT"
	CheckTokenEndpoint = "CHECK_TOKEN_ENDPOINT"
	OPHost             = "OP_HOST"
	OPURL              = "OP_URL"

	LogFilePath = "LOG_FILEPATH"

	// allowed redirect uri
	AllowedRedirectURL = "ALLOWED_REDIRECT_URL"
)

func run() {
	// load env file
	err := godotenv.Load("op.env")
	if err != nil {
		log.Fatal(err.Error())
	}
	// set variables config
	httpconf := config.HTTP{
		AuthEndpoint:       os.Getenv(AuthEndpoint),
		LoginEndpoint:      os.Getenv(LoginEndpoint),
		TokenEndpoint:      os.Getenv(TokenEndpoint),
		CheckTokenEndpoint: os.Getenv(CheckTokenEndpoint),
		OPHost:             os.Getenv(OPHost),
		OPURL:              os.Getenv(OPURL),
	}
	logconf := config.Logger{
		FilePath: os.Getenv(LogFilePath),
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
	// get storage
	tokenStorage := storage.NewAccessToken()
	authCodeStorage := storage.NewAuthCode()
	requestStorage := storage.NewAuthRequest()
	clientStorage := storage.NewClient()
	userStorage := storage.NewUser()

	// add test user
	user := oidc.User{ID: "username", Password: "password"}
	userStorage.Add(user)

	// add test client
	client := oidc.Client{ID: "web", Secret: "secret", RedirectURIs: []string{os.Getenv(AllowedRedirectURL)}}
	clientStorage.Add(client)

	// get controllers
	authCtrl := controller.NewAuth(logger, httpconf, requestStorage, clientStorage)
	loginCtrl := controller.NewLogin(logger, userStorage, requestStorage, authCodeStorage)
	tokenCtrl := controller.NewToken(logger, authCodeStorage, clientStorage, tokenStorage)

	// get router
	r := router.New(httpconf, authCtrl, loginCtrl, tokenCtrl)

	logger.Info("run server on " + httpconf.OPHost)
	http.ListenAndServe(httpconf.OPHost, r)
}

func main() {
	run()
}
